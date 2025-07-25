package magicnumber

// Package file text.go contains the functions that parse bytes as common text and document formats.

import (
	"bytes"
	"io"
)

// NotASCII returns true if the byte is not an printable ASCII character.
// Most control characters are not printable ASCII characters, but an exception
// is made for the ESC (escape) character which is used in ANSI escape codes and
// the EOF (end of file) character which is used in DOS.
func NotASCII(b byte) bool {
	// a list of rune literals for the control characters
	// https://go.dev/ref/spec#Rune_literals
	const (
		nul = 0x0
		tab = byte('\t')
		nl  = byte('\n')
		vt  = byte('\v')
		ff  = byte('\f')
		cr  = byte('\r')
		bel = byte('\a')
		bak = byte('\b')
		eof = 0x1a // end of file character commonly used in DOS
		esc = 0x1b // escape character used in ANSI escape codes
	)
	return (b < 0x20 || b > 0x7f) &&
		b != nul && b != tab && b != nl && b != vt && b != ff && b != cr && b != bel && b != bak &&
		b != esc && b != eof
}

// NotPlainText returns true if the byte is not a printable plain text character.
// This includes any printable ASCII character as well as any "extended ASCII".
func NotPlainText(b byte) bool {
	if !NotASCII(b) {
		return false
	}
	const extendedBegin = 0x80
	const extendedEnd = 0xff
	ExtendedASCII := b >= extendedBegin && b <= extendedEnd
	return !ExtendedASCII
}

// NonISO889591 returns true if the byte is not a printable ISO/IEC-8895-1 character.
func NonISO889591(b byte) bool {
	if !NotASCII(b) {
		return false
	}
	const extendedBegin = 0xa0
	const extendedEnd = 0xff
	ExtendedASCII := b >= extendedBegin && b <= extendedEnd
	return !ExtendedASCII
}

// NonWindows1252 returns true if the byte is not a printable Windows-1252 character.
func NonWindows1252(b byte) bool {
	if !NonISO889591(b) {
		return false
	}
	const (
		extendedBegin = 0x80
		extendedEnd   = 0xff
		unused81      = 0x81
		unused8d      = 0x8d
		unused8f      = 0x8f
		unused90      = 0x90
		unused9d      = 0x9d
	)
	ExtraTypography := b != unused81 && b != unused8d && b != unused8f && b != unused90 && b != unused9d
	return b < extendedBegin || b > extendedEnd || !ExtraTypography
}

// ASCII returns true if the reader exclusively contains printable ASCII characters.
// Today, ASCII characters are the first characters of the Unicode character set
// but historically it was a 7 and 8-bit character encoding standard found on
// most microcomputers, personal computers, and the early Internet.
func ASCII(r io.ReaderAt) bool {
	size := Length(r)
	const chunkSize = 1024
	buf := make([]byte, chunkSize)
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return false
		}
		for i := range n {
			if NotASCII(buf[i]) {
				return false
			}
		}
		if err == io.EOF {
			break
		}
	}
	return true
}

// CodePage returns true if the reader contains is a possible IBM code page
// text file that was often found on DOS and 16-bit Windows computers.
//
// This function is heuristic and checks for the following:
//   - no multiple nulls before the EOF marker
//   - require IBM PC/Microsoft newlines
//   - number of newlines should be at least (80 columns / length of file) / halved
func CodePage(r io.ReaderAt) bool {
	nulpair := []byte{0x0, 0x0}
	msdosNL := []byte{0x0d, 0x0a}
	size := Length(r)
	const chunkSize = 1024
	const binary, textfile = false, true
	newlineCount := 0
	buf := make([]byte, chunkSize)
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return binary
		}
		if pos := bytes.Index(buf[:n], nulpair); pos != -1 {
			return binary
		}
		newlineCount += bytes.Count(buf[:n], msdosNL)
		if err == io.EOF {
			break
		}
	}
	const columns = int64(80)
	if size > columns {
		return int64(newlineCount) >= (size/columns)/2
	}
	return textfile
}

// CSI returns true if the reader contains three or more common Control Sequence Introducer (CSI) escape codes
// that are used in ANSI encoded texts. This is a heuristic function and does not guarantee that the reader
// contains ANSI encoded text.
func CSI(r io.ReaderAt) bool {
	const esc, leftBracket = 0x1b, 0x5b
	const minRequired = 3
	csi := []byte{esc, leftBracket, 0x0}
	codes := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'J', 'K', '=', 's', 'u', '#'}
	finds := 0
	size := Length(r)
	const chunkSize = 1024
	buf := make([]byte, chunkSize)
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return false
		}
		for _, c := range codes {
			if finds >= minRequired {
				return true
			}
			csi[2] = byte(c)
			if pos := bytes.Index(buf[:n], csi); pos > -1 {
				finds++
				continue
			}
		}
		if err == io.EOF {
			break
		}
	}
	return false
}

// Ansi returns true if the reader contains some common ANSI escape codes.
// It for speed and to avoid false positives it only matches the ANSI escape codes
// for bold, normal and reset text.
func Ansi(r io.ReaderAt) bool {
	const esc = 0x1b
	var (
		reset   = []byte{esc, '[', '0', 'm'}
		restart = []byte{esc, '[', '2', 'J'}
		bold    = []byte{esc, '[', '1', ';'}
		normal  = []byte{esc, '[', '0', ';'}
	)
	// check for the common ANSI escape codes
	size := Length(r)
	const chunkSize = 1024
	buf := make([]byte, chunkSize)
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return false
		}

		if pos := bytes.Index(buf[:n], reset); pos != -1 {
			return true
		}
		if pos := bytes.Index(buf[:n], restart); pos != -1 {
			return true
		}
		if pos := bytes.Index(buf[:n], bold); pos != -1 {
			return true
		}
		if pos := bytes.Index(buf[:n], normal); pos != -1 {
			return true
		}
		if err == io.EOF {
			break
		}
	}
	return false
}

// Hlp returns true if the reader contains the Windows Help File signature.
// This is a generic signature for Windows help files and does not differentiate between
// the various versions of the help file format.
func Hlp(r io.ReaderAt) bool {
	const size = 4
	p := make([]byte, size)
	sr := io.NewSectionReader(r, 0, size)
	if n, err := sr.Read(p); err != nil || n < size {
		return false
	}
	compiledHTML := []byte{'I', 'T', 'S', 'F'}
	windowsHelpLN := []byte{'L', 'N', 0x2, 0x0}
	windowsHelp := []byte{'?', 0x5f, 0x3, 0x0}
	help := bytes.Equal(p, compiledHTML) ||
		bytes.Equal(p, windowsHelp) ||
		bytes.Equal(p, windowsHelpLN)
	if help {
		return true
	}
	const offset, size6b = 6, 6
	p = make([]byte, size6b)
	sr = io.NewSectionReader(r, offset, size6b)
	if n, err := sr.Read(p); err != nil || n < size6b {
		return false
	}
	windowsHelp6byte := []byte{0x0, 0x0, 0xff, 0xff, 0xff, 0xff}
	return bytes.Equal(p, windowsHelp6byte)
}

// Pdf returns true if the reader contains the Portable Document Format signature.
func Pdf(r io.ReaderAt) bool {
	const size = 4
	p := make([]byte, size)
	sr := io.NewSectionReader(r, 0, size)
	if n, err := sr.Read(p); err != nil || n < size {
		return false
	}
	if !bytes.Equal(p, []byte{'%', 'P', 'D', 'F'}) {
		return false
	}
	length := Length(r)
	endoffileMarks := [][]byte{
		{0x0a, '%', '%', 'E', 'O', 'F'},
		{0x0a, '%', '%', 'E', 'O', 'F', 0x0a},
		{0x0d, 0x0a, '%', '%', 'E', 'O', 'F', 0x0d, 0x0a},
		{0x0d, '%', '%', 'E', 'O', 'F', 0x0d},
	}
	for _, eof := range endoffileMarks {
		eofSize := int64(len(eof))
		offset := length - eofSize
		p := make([]byte, eofSize)
		sr := io.NewSectionReader(r, offset, eofSize)
		if n, err := sr.Read(p); err != nil || int64(n) < eofSize {
			continue
		}
		if bytes.HasSuffix(p, eof) {
			return true
		}
	}
	return false
}

// Rtf returns true if the reader contains the Rich Text Format signature.
func Rtf(r io.ReaderAt) bool {
	const size = 5
	p := make([]byte, size)
	sr := io.NewSectionReader(r, 0, size)
	if n, err := sr.Read(p); err != nil || n < size {
		return false
	}
	if !bytes.Equal(p, []byte{'{', 0x5c, 'r', 't', 'f'}) {
		return false
	}
	length := Length(r)
	p = make([]byte, 1)
	sr = io.NewSectionReader(r, length-1, 1)
	if n, err := sr.Read(p); err != nil || n < 1 {
		return false
	}
	return p[0] == '}'
}

// Txt returns true if the reader exclusively contains plain text ASCII characters,
// control characters or "extended ASCII characters".
//
// There is a 2% threshold for non-plain text characters such as ASCII control characters
// which are not printable but often found in plain text files for 8-bit microcomputers.
func Txt(r io.ReaderAt) bool {
	const chunkSize = 1024
	size := Length(r)
	buf := make([]byte, chunkSize)
	nonPlainText := 0
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return false
		}
		for i := range n {
			if NotPlainText(buf[i]) {
				nonPlainText++
				if !threshold(nonPlainText, size) {
					return false
				}
			}
		}
		if err == io.EOF {
			break
		}
	}
	return threshold(nonPlainText, size)
}

// If count is greater than 2% of the file size, then it is not plain text.
func threshold(count int, size int64) bool {
	const percentage = 0.02
	return float64(count)/float64(size) < percentage
}

// TxtLatin1 returns true if the reader exclusively contains plain text ISO/IEC-8895-1 characters,
// commonly known as the Latin-1 character set.
func TxtLatin1(r io.ReaderAt) bool {
	size := Length(r)
	const chunkSize = 1024
	buf := make([]byte, chunkSize)
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return false
		}
		for i := range n {
			if NonISO889591(buf[i]) {
				return false
			}
		}
		if err == io.EOF {
			break
		}
	}
	return true
}

// TxtWindows returns true if the reader exclusively contains plain text Windows-1252 characters.
// This is an extension of the Latin-1 character set with additional typography characters and was
// the default character set for English in Microsoft Windows up to Windows 7?
func TxtWindows(r io.ReaderAt) bool {
	size := Length(r)
	const chunkSize = 1024
	buf := make([]byte, chunkSize)
	for offset := int64(0); offset < size; offset += chunkSize {
		bytesToRead := chunkSize
		if offset+int64(chunkSize) > size {
			bytesToRead = int(size - offset)
		}
		n, err := r.ReadAt(buf[:bytesToRead], offset)
		if err != nil && err != io.EOF {
			return false
		}
		for i := range n {
			if NonWindows1252(buf[i]) {
				return false
			}
		}
		if err == io.EOF {
			break
		}
	}
	return true
}

// Utf8 returns true if the reader begins with the UTF-8 Byte Order Mark signature.
func Utf8(r io.ReaderAt) bool {
	const size = 3
	p := make([]byte, size)
	sr := io.NewSectionReader(r, 0, size)
	if n, err := sr.Read(p); err != nil || n < size {
		return false
	}
	return bytes.Equal(p, []byte{0xef, 0xbb, 0xbf})
}

// Utf16 returns true if the reader beings with the UTF-16 Byte Order Mark signature.
func Utf16(r io.ReaderAt) bool {
	const size = 2
	p := make([]byte, size)
	sr := io.NewSectionReader(r, 0, size)
	if n, err := sr.Read(p); err != nil || n < size {
		return false
	}
	return bytes.Equal(p, []byte{0xff, 0xfe}) || bytes.Equal(p, []byte{0xfe, 0xff})
}

// Utf32 returns true if the reader beings with the UTF-32 Byte Order Mark signature.
func Utf32(r io.ReaderAt) bool {
	const size = 4
	p := make([]byte, size)
	sr := io.NewSectionReader(r, 0, size)
	if n, err := sr.Read(p); err != nil || n < size {
		return false
	}
	return bytes.Equal(p, []byte{0xff, 0xfe, 0x0, 0x0}) || bytes.Equal(p, []byte{0x0, 0x0, 0xfe, 0xff})
}
