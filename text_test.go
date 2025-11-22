package magicnumber_test

import (
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
