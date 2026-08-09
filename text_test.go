package magicnumber_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

// go test -v -run "TestText"

func TestTextASCII(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(asciiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(uncompress(txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))

	r, err = os.Open(uncompress(badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(tdfile(t, manualFile))
	be.Err(t, err, nil)
	defer r.Close()

	be.True(t, !magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.PlainText)

	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.PlainText)

	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.PlainText)
}

func TestTextANSI(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ansiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ansi(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.ANSIEscapeText)
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.ANSIEscapeText)
	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.ANSIEscapeText)

	r, err = os.Open(uncompress(txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))

	r, err = os.Open(uncompress(badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	s := "ANSI \x1b[2Jtext"
	nr := strings.NewReader(s)
	be.True(t, magicnumber.Ansi(nr))
	s = "ANSI \x1b[0;text"
	nr = strings.NewReader(s)
	be.True(t, magicnumber.Ansi(nr))
	s = "ANSI \x1b[1;text"
	nr = strings.NewReader(s)
	be.True(t, magicnumber.Ansi(nr))

	s = "ANSI \x1btext"
	nr = strings.NewReader(s)
	be.True(t, !magicnumber.Ansi(nr))
}

func TestTextCSI(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ansiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.CSI(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.ANSIEscapeText)
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.ANSIEscapeText)
	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.ANSIEscapeText)

	r, err = os.Open(uncompress(txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))

	r, err = os.Open(uncompress(badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	s := "ANSI \x1b[2Jtext"
	nr := strings.NewReader(s)
	be.True(t, !magicnumber.CSI(nr))
	s = "ANSI \x1b[0;text"
	nr = strings.NewReader(s)
	be.True(t, !magicnumber.CSI(nr))
	s = "ANSI \x1b[1;t\x1b[2Je\x1b[0;xt"
	nr = strings.NewReader(s)
	be.True(t, magicnumber.CSI(nr))

	s = "ANSI \x1btext"
	nr = strings.NewReader(s)
	be.True(t, !magicnumber.CSI(nr))
}

func TestTextRTF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(rtfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Rtf(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.RichTextFormat)
}

func TestTextPDF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pdfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pdf(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PortableDocumentFormat)
	sign, err := magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.PortableDocumentFormat)
}

func TestTextUTF16(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(utf16File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Utf16(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.UTF16Text)
}

func TestTextISO7(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(iso7File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))
	be.True(t, !magicnumber.Ansi(r))
	be.True(t, magicnumber.Txt(r))
	be.True(t, magicnumber.TxtLatin1(r))
	be.True(t, magicnumber.TxtWindows(r))
	be.True(t, !magicnumber.Utf8(r))
	be.True(t, !magicnumber.Utf16(r))
	be.True(t, !magicnumber.Utf32(r))
}

func TestTextByte(t *testing.T) {
	t.Parallel()
	b := byte(0x90)
	be.True(t, magicnumber.NonWindows1252(b))
	b = byte('a')
	be.True(t, !magicnumber.NonWindows1252(b))
}

func TestTextCodePage(t *testing.T) {
	t.Parallel()
	r, err := os.Open(tdfile(t, "TRIAD.TXT"))
	be.Err(t, err, nil)
	defer r.Close()

	be.True(t, magicnumber.Txt(r))
	be.True(t, magicnumber.CodePage(r))

	const want = magicnumber.PlainText
	be.Equal(t, magicnumber.Find(r), want)
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, want)
	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, want)
}

func TestTextBinaries(t *testing.T) {
	t.Parallel()
	r, err := os.Open(tdfile(t, "binarytxt.bin"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.XBin(r))
	r, err = os.Open(tdfile(t, "binarytxt.xb"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.XBin(r))
}

func TestTextAnsiW_PositionBug(t *testing.T) {
	t.Parallel()
	const count = 1030
	const bold = "\x1b[1;"
	padding := bytes.Repeat([]byte{'A'}, count)
	b := []byte(bold)
	b = append(padding, b...)

	var w bytes.Buffer
	r := bytes.NewReader(b)
	got := magicnumber.AnsiW(&w, r)
	be.True(t, got)

	s := w.String()
	substr := "position 1030"
	got = strings.Contains(s, substr)
	be.True(t, got)
	if !got {
		const format = "%q: expected a total size of %d"
		fmt.Fprintf(os.Stderr, format, s, count)
	}
}

func TestTextTxtW(t *testing.T) {
	t.Parallel()
	t.Run("exceeds 2%", func(t *testing.T) {
		t.Parallel()
		b := make([]byte, 100)
		for i := range b {
			b[i] = 'A'
		}
		b[0], b[1], b[2] = 0x01, 0x02, 0x03 // 3% bad bytes
		r := bytes.NewReader(b)
		got := magicnumber.TxtW(io.Discard, r)
		be.True(t, !got)
	})

	t.Run("1% bad byte", func(t *testing.T) {
		t.Parallel()
		b := make([]byte, 100)
		for i := range b {
			b[i] = 'A'
		}
		b[0] = 0x01
		r := bytes.NewReader(b)
		got := magicnumber.TxtW(io.Discard, r)
		be.True(t, got)
	})

	t.Run("boundaries", func(t *testing.T) {
		t.Parallel()
		b := make([]byte, 1500)
		for i := range b {
			b[i] = 'B'
		}
		r := bytes.NewReader(b)
		got := magicnumber.TxtW(io.Discard, r)
		be.True(t, got)
	})
}
