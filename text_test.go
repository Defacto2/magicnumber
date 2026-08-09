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

func TestASCII(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(asciiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ASCII(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ASCII(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))

	r, err = os.Open(uncompress(badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(tdfile(manualFile))
	be.Err(t, err, nil)
	defer r.Close()

	be.True(t, !magicnumber.ASCII(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.PlainText, sign)

	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.PlainText, sign)

	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.PlainText, sign)
}

func TestANSI(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ansiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ansi(r))
	be.Equal(t, magicnumber.ANSIEscapeText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.ANSIEscapeText, sign)
	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.ANSIEscapeText, sign)

	r, err = os.Open(uncompress(txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))

	r, err = os.Open(uncompress(badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

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

func TestCSI(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ansiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.CSI(r))
	be.Equal(t, magicnumber.ANSIEscapeText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.ANSIEscapeText, sign)
	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.ANSIEscapeText, sign)

	r, err = os.Open(uncompress(txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))

	r, err = os.Open(uncompress(badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))
	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

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

func TestRTF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(rtfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Rtf(r))
	be.Equal(t, magicnumber.RichTextFormat, magicnumber.Find(r))
}

func TestPDF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pdfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pdf(r))
	be.Equal(t, magicnumber.PortableDocumentFormat, magicnumber.Find(r))
	sign, err := magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.PortableDocumentFormat, sign)
}

func TestUTF16(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(utf16File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Utf16(r))
	be.Equal(t, magicnumber.UTF16Text, magicnumber.Find(r))
}

func TestISO7(t *testing.T) {
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

func TestByte(t *testing.T) {
	t.Parallel()
	b := byte(0x90)
	be.True(t, magicnumber.NonWindows1252(b))
	b = byte('a')
	be.True(t, !magicnumber.NonWindows1252(b))
}

func TestCodePage(t *testing.T) {
	t.Parallel()
	r, err := os.Open(tdfile("TRIAD.TXT"))
	be.Err(t, err, nil)
	defer r.Close()

	be.True(t, magicnumber.Txt(r))
	be.True(t, magicnumber.CodePage(r))

	be.Equal(t, magicnumber.PlainText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.PlainText, sign)
	sign, err = magicnumber.Document(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.PlainText, sign)
}

func TestBinaryTexts(t *testing.T) {
	t.Parallel()
	r, err := os.Open(tdfile("binarytxt.bin"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.XBin(r))
	r, err = os.Open(tdfile("binarytxt.xb"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.XBin(r))
}

func TestAnsiW_PositionBug(t *testing.T) {
	// Create a payload larger than chunkSize (1024) so it spans multiple chunks.
	// Fill first 1030 bytes with dummy letters 'A', then insert a bold sequence "\x1b[1;"
	// Target sequence starts at exact byte index 1030.

	const count = 1030
	const bold = "\x1b[1;"
	padding := bytes.Repeat([]byte{'A'}, count)
	b := []byte(bold)
	b = append(padding, b...)

	reader := bytes.NewReader(b)
	var buf bytes.Buffer

	got := magicnumber.AnsiW(&buf, reader)
	be.True(t, got)

	s := buf.String()
	expectedPosition := "position 1030"
	got = strings.Contains(s, expectedPosition)
	be.True(t, got)
	if !got {
		fmt.Fprintf(os.Stderr, "%q: expected a total size of %d", s, count)
	}
}

func TestTxtW(t *testing.T) {
	t.Run("exceeds 2%", func(t *testing.T) {
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
		b := make([]byte, 1500)
		for i := range b {
			b[i] = 'B'
		}
		r := bytes.NewReader(b)
		got := magicnumber.TxtW(io.Discard, r)
		be.True(t, got)
	})
}
