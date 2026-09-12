//nolint:exhaustruct_v5,paralleltest
package magicnumber_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
)

var errTestSeek = errors.New("seek failed")

// errSeeker allows injecting failures into specific Seek calls.
type errSeeker struct {
	seekCount  int
	failOnCall int // 1 = SeekCurrent, 2 = SeekEnd, 3 = Restore
}

// mockSeeker implements io.ReaderAt and io.Seeker without implementing sizer.
type mockSeeker struct {
	*bytes.Reader
}

// mockReaderAt implements io.ReaderAt only (no sizer, no seeker).
type mockReaderAt struct{}

func (m mockReaderAt) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, io.EOF
}

func (e *errSeeker) ReadAt(_ []byte, _ int64) (int, error) {
	return 0, io.EOF
}

func (e *errSeeker) Seek(_ int64, whence int) (int64, error) {
	e.seekCount++
	if e.seekCount == e.failOnCall {
		return 0, errTestSeek
	}
	if whence == io.SeekEnd {
		return 100, nil
	}
	return 0, nil
}

func TestLength(t *testing.T) {
	content := []byte("hello world")

	tests := []struct {
		name string
		r    io.ReaderAt
		want int64
	}{
		{
			name: "nil reader",
			r:    nil,
			want: 0,
		},
		{
			name: "bytes.Reader (implements sizer)",
			r:    bytes.NewReader(content),
			want: int64(len(content)),
		},
		{
			name: "io.SectionReader (implements sizer)",
			r:    io.NewSectionReader(bytes.NewReader(content), 0, 5),
			want: 5,
		},
		{
			name: "unsupported reader type",
			r:    mockReaderAt{},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := magicnumber.Length(tt.r); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLength_SeekerFallbackAndRestore(t *testing.T) {
	content := []byte("hello world")
	base := bytes.NewReader(content)

	// Advance pointer to offset 3 before calling Length
	const initialOffset int64 = 3
	_, err := base.Seek(initialOffset, io.SeekStart)
	if err != nil {
		t.Fatalf("failed to setup seeker: %v", err)
	}

	seeker := mockSeeker{Reader: base}

	got := magicnumber.Length(seeker)
	want := int64(len(content))
	if got != want {
		t.Errorf("Length() = %v, want %v", got, want)
	}

	// Verify offset was restored back to initial position (3)
	currentOffset, err := base.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatalf("failed to query offset: %v", err)
	}
	if currentOffset != initialOffset {
		t.Errorf("offset after Length() = %v, want %v", currentOffset, initialOffset)
	}
}

func TestLength_File(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "length_test_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	data := []byte("temp file content")
	if _, err := tmp.Write(data); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if got := magicnumber.Length(tmp); got != int64(len(data)) {
		t.Errorf("Length(file) = %v, want %v", got, len(data))
	}
}

func TestLength_SeekErrors(t *testing.T) {
	tests := []struct {
		name       string
		failOnCall int
		want       int64
	}{
		{
			name:       "fail on first Seek (SeekCurrent)",
			failOnCall: 1,
			want:       0,
		},
		{
			name:       "fail on second Seek (SeekEnd)",
			failOnCall: 2,
			want:       0,
		},
		{
			name:       "fail on third Seek (restore offset)",
			failOnCall: 3,
			want:       100, // Length returns the length found at SeekEnd even if restoring offset fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &errSeeker{failOnCall: tt.failOnCall}
			if got := magicnumber.Length(s); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}
