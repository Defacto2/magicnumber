package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
)

// BenchmarkISODetection measures the performance of ISO 9660 detection.
func BenchmarkISODetection(b *testing.B) {
	f, err := os.Open(imgfile("uncompress.iso"))
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	b.ResetTimer()
	for b.Loop() {
		// Reset file position for each iteration
		_, err := f.Seek(0, 0)
		if err != nil {
			b.Fatal(err)
		}
		// Benchmark the ISO detection function
		_ = magicnumber.ISO(f)
	}
}

// BenchmarkISOFind measures the performance of ISO detection through the Find function.
func BenchmarkISOFind(b *testing.B) {
	f, err := os.Open(imgfile("uncompress.iso"))
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err := f.Seek(0, 0)
		if err != nil {
			b.Fatal(err)
		}
		// Benchmark detection through the general Find function
		_ = magicnumber.Find(f)
	}
}

// BenchmarkISODiscImage measures the performance of ISO detection through the DiscImage function.
func BenchmarkISODiscImage(b *testing.B) {
	f, err := os.Open(imgfile("uncompress.iso"))
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	b.ResetTimer()
	for b.Loop() {
		_, err := f.Seek(0, 0)
		if err != nil {
			b.Fatal(err)
		}
		// Benchmark detection through the DiscImage function
		_, _ = magicnumber.DiscImage(f)
	}
}
