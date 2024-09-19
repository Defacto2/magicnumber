package magicnumber_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestASCII(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(asciiFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.ASCII(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(txtFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.ASCII(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(gifFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.ASCII(r))

	r, err = os.Open(uncompress(badFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.ASCII(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(tdfile(manualFile))
	require.NoError(t, err)
	defer r.Close()

	assert.False(t, magicnumber.ASCII(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PlainText, sign)

	sign, err = magicnumber.Document(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PlainText, sign)

	sign, err = magicnumber.Document(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PlainText, sign)
}

func TestANSI(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ansiFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Ansi(r))
	assert.Equal(t, magicnumber.ANSIEscapeText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.ANSIEscapeText, sign)
	sign, err = magicnumber.Document(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.ANSIEscapeText, sign)

	r, err = os.Open(uncompress(txtFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.Ansi(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(gifFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.Ansi(r))

	r, err = os.Open(uncompress(badFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.Ansi(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	s := "ANSI \x1b[2Jtext"
	nr := strings.NewReader(s)
	assert.True(t, magicnumber.Ansi(nr))
	s = "ANSI \x1b[0;text"
	nr = strings.NewReader(s)
	assert.True(t, magicnumber.Ansi(nr))
	s = "ANSI \x1b[1;text"
	nr = strings.NewReader(s)
	assert.True(t, magicnumber.Ansi(nr))

	s = "ANSI \x1btext"
	nr = strings.NewReader(s)
	assert.False(t, magicnumber.Ansi(nr))
}

func TestCSI(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ansiFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.CSI(r))
	assert.Equal(t, magicnumber.ANSIEscapeText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.ANSIEscapeText, sign)
	sign, err = magicnumber.Document(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.ANSIEscapeText, sign)

	r, err = os.Open(uncompress(txtFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.CSI(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	r, err = os.Open(uncompress(gifFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.CSI(r))

	r, err = os.Open(uncompress(badFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.CSI(r))
	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))

	s := "ANSI \x1b[2Jtext"
	nr := strings.NewReader(s)
	assert.False(t, magicnumber.CSI(nr))
	s = "ANSI \x1b[0;text"
	nr = strings.NewReader(s)
	assert.False(t, magicnumber.CSI(nr))
	s = "ANSI \x1b[1;t\x1b[2Je\x1b[0;xt"
	nr = strings.NewReader(s)
	assert.True(t, magicnumber.CSI(nr))

	s = "ANSI \x1btext"
	nr = strings.NewReader(s)
	assert.False(t, magicnumber.CSI(nr))
}

func TestRTF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(rtfFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Rtf(r))
	assert.Equal(t, magicnumber.RichTextFormat, magicnumber.Find(r))
}

func TestPDF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pdfFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Pdf(r))
	assert.Equal(t, magicnumber.PortableDocumentFormat, magicnumber.Find(r))
	sign, err := magicnumber.Document(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PortableDocumentFormat, sign)
}

func TestUTF16(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(utf16File))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Utf16(r))
	assert.Equal(t, magicnumber.UTF16Text, magicnumber.Find(r))
}

func TestISO7(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(iso7File))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.ASCII(r))
	assert.False(t, magicnumber.Ansi(r))
	assert.True(t, magicnumber.Txt(r))
	assert.True(t, magicnumber.TxtLatin1(r))
	assert.True(t, magicnumber.TxtWindows(r))
	assert.False(t, magicnumber.Utf8(r))
	assert.False(t, magicnumber.Utf16(r))
	assert.False(t, magicnumber.Utf32(r))
}

func TestByte(t *testing.T) {
	t.Parallel()
	b := byte(0x90)
	assert.True(t, magicnumber.NonWindows1252(b))
	b = byte('a')
	assert.False(t, magicnumber.NonWindows1252(b))
}

func TestCodePage(t *testing.T) {
	t.Parallel()
	r, err := os.Open(tdfile("TRIAD.TXT"))
	require.NoError(t, err)
	defer r.Close()

	assert.False(t, magicnumber.Txt(r))
	assert.True(t, magicnumber.CodePage(r))

	assert.Equal(t, magicnumber.PlainText, magicnumber.Find(r))
	sign, err := magicnumber.Text(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PlainText, sign)
	sign, err = magicnumber.Document(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PlainText, sign)
}
