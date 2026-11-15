//nolint:exhaustruct_v5
package magicnumber_test

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

// go test -v -run "TestText"

func TestTextASCII(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, asciiFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(pathUncompress(t, txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(pathUncompress(t, gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))

	r, err = os.Open(pathUncompress(t, badFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ASCII(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(pathFile(t, manualFile))
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
	r, err := os.Open(pathUncompress(t, ansiFile))
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

	r, err = os.Open(pathUncompress(t, txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(pathUncompress(t, gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Ansi(r))

	r, err = os.Open(pathUncompress(t, badFile))
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

	r, err := os.Open(pathUncompress(t, ansiFile))
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

	r, err = os.Open(pathUncompress(t, txtFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PlainText)

	r, err = os.Open(pathUncompress(t, gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.CSI(r))

	r, err = os.Open(pathUncompress(t, badFile))
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

	r, err := os.Open(pathUncompress(t, rtfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Rtf(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.RichTextFormat)
}

func TestTextPDF(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, pdfFile))
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

	r, err := os.Open(pathUncompress(t, utf16File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Utf16(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.UTF16Text)
}

func TestTextISO7(t *testing.T) {
	t.Parallel()
	r, err := os.Open(pathUncompress(t, iso7File))
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

	r, err := os.Open(pathFile(t, "TRIAD.TXT"))
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

	r, err := os.Open(pathFile(t, "binarytxt.bin"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.XBin(r))

	r, err = os.Open(pathFile(t, "binarytxt.xb"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.XBin(r))
}

func TestTextAnsi_PositionBug(t *testing.T) {
	t.Parallel()

	const count = 1030
	const bold = "\x1b[1;"
	padding := bytes.Repeat([]byte{'A'}, count)
	b := []byte(bold)
	b = append(padding, b...)

	r := bytes.NewReader(b)
	var w bytes.Buffer
	sl := slog.New(slog.NewJSONHandler(&w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	got := magicnumber.AnsiWithLogger(sl, r)
	be.True(t, got)

	s := w.String()
	be.True(t, strings.Contains(s, `"n":1030`))
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
		got := magicnumber.TxtWithLogger(nil, r)
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
		got := magicnumber.TxtWithLogger(nil, r)
		be.True(t, got)
	})

	t.Run("boundaries", func(t *testing.T) {
		t.Parallel()

		b := make([]byte, 1500)
		for i := range b {
			b[i] = 'B'
		}
		r := bytes.NewReader(b)
		got := magicnumber.TxtWithLogger(nil, r)
		be.True(t, got)
	})
}

func TestDumpCGA(t *testing.T) {
	t.Parallel()

	t.Log("TestDumpCGA")
	r, err := os.Open(pathFile(t, "SPYHUNT.TPG"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.DumpCGA(r))

	got := magicnumber.Find(r)
	be.Equal(t, got, magicnumber.CGAVideoDump)
	s := got.Title()
	be.Equal(t, s, "IBM CGA video RAM dump")

	// 4096-byte plain ASCII text should be PlainText, not CGAVideoDump.
	ascii4k := strings.NewReader(strings.Repeat("Hello world!\n", 316)[:4096])
	be.Equal(t, magicnumber.Find(ascii4k), magicnumber.PlainText)

	// All-zero 4096-byte file must not match DumpCGA.
	zeros := bytes.NewReader(make([]byte, 4096))
	be.True(t, !magicnumber.DumpCGA(zeros))

	// BSAVE header (0xFD, B800:0000, length 4000) + 80x25 CGA screen.
	bsave := make([]byte, 7+4000)
	bsave[0] = 0xFD
	bsave[1], bsave[2] = 0x00, 0xB8
	bsave[5], bsave[6] = 0xA0, 0x0F // 4000 little-endian
	for i := 7; i < len(bsave); i += 2 {
		bsave[i] = ' '
		bsave[i+1] = 0x07
	}
	br := bytes.NewReader(bsave)
	be.True(t, magicnumber.DumpCGA(br))
	be.Equal(t, magicnumber.Find(br), magicnumber.CGAVideoDump)

	// 80x43 EGA text dump (6880 bytes) with blink (0x8E) on black background.
	ega := make([]byte, 6880)
	for i := 0; i < len(ega); i += 2 {
		ega[i] = '#'
		if i < 160 {
			ega[i+1] = 0x8E
		} else {
			ega[i+1] = 0x0F
		}
	}
	be.True(t, magicnumber.DumpCGA(bytes.NewReader(ega)))
}

func TestDumpCGAGraphics(t *testing.T) {
	t.Parallel()

	f, err := os.Open(pathFile(t, "cga/OPEN"))
	be.Err(t, err, nil)
	defer f.Close()
	be.True(t, magicnumber.DumpCGAGraphics(f))
	got := magicnumber.Find(f)
	be.Equal(t, got, magicnumber.CGAGraphicsDump)
	be.Equal(t, got.Title(), "IBM CGA graphics RAM dump")

	raw := make([]byte, 16384)
	raw[0] = 0x55
	raw[8192] = 0xAA
	r := bytes.NewReader(raw)
	be.True(t, magicnumber.DumpCGAGraphics(r))
	got = magicnumber.Find(r)
	be.Equal(t, got, magicnumber.CGAGraphicsDump)
	be.Equal(t, got.Title(), "IBM CGA graphics RAM dump")

	// BSAVE-wrapped 16KB CGA graphics dump
	bsave := make([]byte, 7+16384)
	bsave[0] = 0xFD
	bsave[1], bsave[2] = 0x00, 0xB8
	bsave[5], bsave[6] = 0x00, 0x40 // 16384 little-endian
	copy(bsave[7:], raw)
	br := bytes.NewReader(bsave)
	be.True(t, magicnumber.DumpCGAGraphics(br))
	be.Equal(t, magicnumber.Find(br), magicnumber.CGAGraphicsDump)

	// All-zero 16KB file must not match
	be.True(t, !magicnumber.DumpCGAGraphics(bytes.NewReader(make([]byte, 16384))))

	// Non-zero byte in inter-bank gap must not match
	badGap := append([]byte{}, raw...)
	badGap[8000] = 0x01
	be.True(t, !magicnumber.DumpCGAGraphics(bytes.NewReader(badGap)))
}
