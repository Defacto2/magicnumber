//nolint:cyclop
package magicnumber

// Package file text.go contains the functions that parse bytes as common text and document formats.

import (
	"bufio"
	"bytes"
	"io"
	"log/slog"
)

const esc = 0x1b // ANSI escape

// NotASCII returns true if the byte is not a printable ASCII character.
// Most control characters are not printable ASCII characters, but an exception
// is made for the ESC (escape) character which is used in ANSI escape codes and
// the EOF (end of file) character which is used in DOS.
func NotASCII(b byte) bool {
	if b >= 0x20 && b <= 0x7F {
		return false
	}

	const (
		nul = 0x00 // Null character / CP437 blank
		bel = '\a' // Alert / Bell
		bak = '\b' // Backspace
		tab = '\t' // Horizontal Tab
		nl  = '\n' // Line Feed / Newline
		vt  = '\v' // Vertical Tab
		ff  = '\f' // Form Feed
		cr  = '\r' // Carriage Return
		eof = 0x1A // DOS End-Of-File
	)

	switch b {
	case nul, bel, bak, tab, nl, vt, ff, cr, eof, esc:
		return false // Allowed control / formatting character
	default:
		return true // Non-printable / non-ASCII
	}
}

// NotPlainText returns true if the byte is not a printable plain text character.
// This includes any printable ASCII character as well as any "extended ASCII".
func NotPlainText(b byte) bool {
	if !NotASCII(b) {
		return false
	}
	const extendedBegin = 0x80
	const extendedEnd = 0xff
	extASCII := b >= extendedBegin && b <= extendedEnd
	return !extASCII
}

// NonISO889591 returns true if the byte is not a printable ISO/IEC-8859-1 character.
func NonISO889591(b byte) bool {
	if !NotASCII(b) {
		return false
	}
	const extendedBegin = 0xa0
	const extendedEnd = 0xff
	extASCII := b >= extendedBegin && b <= extendedEnd
	return !extASCII
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
	extTypography := b != unused81 && b != unused8d && b != unused8f && b != unused90 && b != unused9d
	return b < extendedBegin || b > extendedEnd || !extTypography
}

// ASCII returns true if the reader exclusively contains printable ASCII characters.
// Today, ASCII characters are the first characters of the Unicode character set
// but historically it was a 7 and 8-bit character encoding standard found on
// most microcomputers, personal computers, and the early Internet.
func ASCII(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	const chunkSize = 1024
	p := make([]byte, chunkSize)
	var off int64

	for {
		n, err := r.ReadAt(p, off)

		for i := range n {
			if NotASCII(p[i]) {
				return false
			}
		}

		off += int64(n)
		if err != nil {
			if err == io.EOF {
				break
			}
			return false
		}

		if n == 0 {
			break
		}
	}

	return true
}

// CodePage returns true if the reader is a potential IBM code page
// text file that was often in use on DOS and ancient Windows systems.
func CodePage(r io.ReaderAt) bool {
	return CodePageW(io.Discard, r)
}

// Deprecated: use [CodePageWithLogger] instead.
// The io.Writer is unused.
func CodePageW(_ io.Writer, r io.ReaderAt) bool {
	return CodePageWithLogger(nil, r)
}

// CodePageWithLogger returns true if the reader is a potential IBM code page
// text file that was often in use on DOS and ancient Windows systems.
//
// This function is heuristic and checks for unique patterns, otherwise
// it will return false, and should be used in combination with [Txt].
//   - no multiple nulls before the EOF marker
//   - require IBM PC/Microsoft newlines
//   - number of newlines should be at least (80 columns / length of file) / halved
//   - geometric triangle pairs, ▲▼ ◄► ►◄
func CodePageWithLogger(sl *slog.Logger, r io.ReaderAt) bool { //nolint:funlen
	const msg = "heuristic codepage"
	if sl == nil {
		sl = slog.New(slog.DiscardHandler)
	}
	if r == nil {
		return false
	}

	const (
		chunkSize = 1024
		columns   = 80
		split     = 2
		minWords  = 10
	)

	size := Length(r)
	if size == 0 {
		return true // an empty file is treated as valid text
	}
	sep := []byte{0x0D, 0x0A} // ms-dos era new line
	newlines := 0
	p := make([]byte, chunkSize)
	sl.Debug(msg+" reading chunks", slog.Int("chunk size B", chunkSize), slog.Int64("size B", size))

	// track the last byte of the previous chunk to detect CRLF split across boundaries
	var lastByte byte
	hasLastByte := false

	for off := int64(0); off < size; {
		n, err := r.ReadAt(p, off)
		if err != nil && err != io.EOF {
			sl.Debug(msg+" read at error", slog.Int64("offset", off), slog.Any("error", err))
			return false
		}

		chunk := p[:n]
		// check boundary CRLF split (CR at end of prev chunk, LF at start of current)
		if hasLastByte && lastByte == 0x0D && n > 0 && chunk[0] == 0x0A {
			newlines++
		}
		// heuristic check for unique CP437 character pairs
		if ok, match := charPairs(sl, n, chunk); ok {
			return match
		}

		newlines += bytes.Count(chunk, sep)
		if n > 0 {
			lastByte = chunk[n-1] //nolint:gosec
			hasLastByte = true
		}
		off += int64(n)
		if err == io.EOF || n == 0 {
			break
		}
	}

	if size > columns {
		threshold := (size / columns) / split
		hasEnoughNewlines := int64(newlines) >= threshold
		sl.Debug(msg+"newline count >=",
			slog.Int("value", newlines),
			slog.Int64("threshold", threshold),
			slog.Int64("size", size),
			slog.Int("columns", columns),
			slog.Int("split", split),
			slog.Bool("has enough newlines", hasEnoughNewlines),
		)
		if hasEnoughNewlines {
			return true
		}
		return hasMinWords(r, size, minWords)
	}
	sl.Debug(msg + " matched a textfile")
	return true
}

func charPairs(sl *slog.Logger, n int, buf []byte) (bool, bool) {
	const msg = "character pairs"
	const binary, textfile = false, true
	const match = true

	// always use fixed arrays for performance
	nulpair := [2]byte{0x00, 0x00}
	updown := [2]byte{0x1E, 0x1F}    // ▲▼
	leftright := [2]byte{0x11, 0x10} // ◄►
	rightleft := [2]byte{0x10, 0x11} // ►◄

	s := buf[:n]

	if pos := bytes.Index(s, nulpair[:]); pos != -1 {
		sl.Debug(msg + " read to the eof without a marker")
		return match, binary
	}

	if pos := bytes.Index(s, updown[:]); pos != -1 {
		sl.Debug(msg + " returning textfile up-down ▲▼ match")
		return match, textfile
	}

	if pos := bytes.Index(s, leftright[:]); pos != -1 {
		sl.Debug(msg + " returning textfile left-right ◄► match")
		return match, textfile
	}

	if pos := bytes.Index(s, rightleft[:]); pos != -1 {
		sl.Debug(msg + " returning textfile right-left ►◄ match")
		return match, textfile
	}

	return !match, false
}

// hasMinWords checks if the reader contains at least minWords without reading the whole file.
func hasMinWords(r io.ReaderAt, size int64, minWords int) bool {
	if r == nil || size <= 0 {
		return false
	}

	sr := io.NewSectionReader(r, 0, size)
	scanner := bufio.NewScanner(sr)
	scanner.Split(bufio.ScanWords)

	count := 0
	for scanner.Scan() {
		count++
		if count >= minWords {
			return true // Early exit as soon as threshold is met
		}
	}
	_ = scanner.Err()

	return false
}

// CSI returns true if the reader contains three or more common Control Sequence Introducer (CSI) escape codes
// that are used in ANSI encoded texts. This is a heuristic function and does not guarantee that the reader
// contains ANSI encoded text.
func CSI(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	const (
		chunkSize   = 1024
		maxTail     = 2
		minRequired = 3
		common      = `0123456789mHJK?=suABCD#`
	)

	// track up to 2 trailing bytes from the previous chunk to detect boundary-straddling CSI sequences
	var buf [maxTail + chunkSize]byte
	var tail [maxTail]byte
	p := buf[maxTail:]

	finds := 0
	tails := 0
	var off int64

	for {
		n, err := r.ReadAt(p, off)
		if err != nil && err != io.EOF {
			return false
		}
		if n == 0 {
			break
		}

		// prepare a contiguous search slice that includes overlap from the previous chunk without allocating
		start := maxTail - tails
		if tails > 0 {
			copy(buf[start:maxTail], tail[:tails])
		}
		searchBuf := buf[start : maxTail+n]

		// count occurrences of CSI codes (\x1b[ + char) in this chunk
		if bytes.IndexByte(searchBuf, esc) != -1 {
			for _, c := range []byte(common) {
				pattern := [3]byte{esc, '[', c}
				finds += bytes.Count(searchBuf, pattern[:])
				if finds >= minRequired {
					return true
				}
			}
		}

		// save up to the last 2 bytes for the next boundary check
		if len(searchBuf) >= maxTail {
			copy(tail[:], searchBuf[len(searchBuf)-maxTail:])
			tails = maxTail
		} else {
			copy(tail[:], searchBuf)
			tails = len(searchBuf)
		}

		off += int64(n)
		if err == io.EOF {
			break
		}
	}

	return finds >= minRequired
}

// Ansi returns true if the reader contains some common ANSI escape codes.
// It for speed and to avoid false positives it only matches the ANSI escape codes
// for bold, normal and reset text.
func Ansi(r io.ReaderAt) bool {
	return AnsiWithLogger(nil, r)
}

// Deprecated: use [AnsiWithLogger] instead.
// The io.Writer is unused.
func AnsiW(_ io.Writer, r io.ReaderAt) bool {
	return AnsiWithLogger(nil, r)
}

// AnsiWithLogger returns true if the reader contains some common ANSI escape codes.
// It for speed and to avoid false positives it only matches the ANSI escape codes
// for bold, normal and reset text.
func AnsiWithLogger(sl *slog.Logger, r io.ReaderAt) bool { //nolint:funlen
	const msg = "common ansi escape codes"
	if sl == nil {
		sl = slog.New(slog.DiscardHandler)
	}
	if r == nil {
		return false
	}

	size := Length(r)
	const chunkSize = 1024
	sl.Debug(msg, slog.Int("chunk size", chunkSize), slog.Int64("total size B", size))

	reset := [4]byte{esc, '[', '0', 'm'}
	restart := [4]byte{esc, '[', '2', 'J'}
	bold := [4]byte{esc, '[', '1', ';'}
	normal := [4]byte{esc, '[', '0', ';'}

	const maxTail = 3
	var buf [maxTail + chunkSize]byte
	p := buf[maxTail:]

	tails := 0
	var off int64
	var tail [maxTail]byte

	for {
		n, err := r.ReadAt(p, off)
		if err != nil && err != io.EOF {
			sl.Debug(msg+" read at error", slog.Any("error", err))
			return false
		}
		if n == 0 {
			break
		}

		var s []byte
		if tails > 0 {
			// Copy tail into front of buf, slicing s without dynamic heap allocation
			copy(buf[:tails], tail[:tails])
			s = buf[:tails+n]
		} else {
			s = p[:n]
		}

		if pos := bytes.Index(s, reset[:]); pos != -1 {
			n := off - int64(tails) + int64(pos)
			sl.Debug(msg+" reset", slog.Int("position", pos), slog.Int64("n", n))
			return true
		}
		if pos := bytes.Index(s, restart[:]); pos != -1 {
			n := off - int64(tails) + int64(pos)
			sl.Debug(msg+" restart", slog.Int("position", pos), slog.Int64("n", n))
			return true
		}
		if pos := bytes.Index(s, bold[:]); pos != -1 {
			n := off - int64(tails) + int64(pos)
			sl.Debug(msg+" bold", slog.Int("position", pos), slog.Int64("n", n))
			return true
		}
		if pos := bytes.Index(s, normal[:]); pos != -1 {
			n := off - int64(tails) + int64(pos)
			sl.Debug(msg+" normal", slog.Int("position", pos), slog.Int64("n", n))
			return true
		}

		const size = 3
		if n >= size {
			copy(tail[:], p[n-size:n])
			tails = size
		} else {
			copy(tail[:n], p[:n])
			tails = n
		}

		off += int64(n)

		if err == io.EOF {
			sl.Debug(msg+" end of file", slog.Int64("total bytes read", off))
			break
		}
	}
	sl.Debug(msg + " scan found nothing")
	return false
}

// Hlp returns true if the reader contains the Windows Help File signature.
// This is a generic signature for Windows help files and does not differentiate between
// the various versions of the help file format.
func Hlp(r io.ReaderAt) bool {
	if r == nil {
		return false
	}
	// read first 4 bytes directly from offset 0
	var header4 [4]byte
	off := int64(0)
	if n, err := r.ReadAt(header4[:], off); (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	itsf := [4]byte{'I', 'T', 'S', 'F'}
	winHelpLN := [4]byte{'L', 'N', 0x02, 0x00}
	winHelp := [4]byte{'?', 0x5F, 0x03, 0x00}

	if header4 == itsf || header4 == winHelpLN || header4 == winHelp {
		return true
	}

	// read the first 6 bytes directly from offset 6
	var header6 [6]byte
	off = 6
	if n, err := r.ReadAt(header6[:], off); (err != nil && err != io.EOF) || n < 6 {
		return false
	}
	winHelp6B := [6]byte{0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF}
	return header6 == winHelp6B
}

// Pdf returns true if the reader contains the Portable Document Format signature.
func Pdf(r io.ReaderAt) bool {
	if r == nil {
		return false
	}
	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	if p != [4]byte{'%', 'P', 'D', 'F'} {
		return false
	}

	length := Length(r)
	var tail [9]byte

	eofMarkers := [4]string{
		"\x0a%%EOF",
		"\x0a%%EOF\x0a",
		"\x0d\x0a%%EOF\x0d\x0a",
		"\x0d%%EOF\x0d",
	}

	for _, eof := range eofMarkers {
		eofSize := int64(len(eof))
		if length < eofSize {
			continue
		}

		off := length - eofSize
		buf := tail[:eofSize]
		if n, err := r.ReadAt(buf, off); (err != nil && err != io.EOF) || int64(n) < eofSize {
			continue
		}
		if string(buf) == eof {
			return true
		}
	}
	return false
}

// Rtf returns true if the reader contains the Rich Text Format signature.
func Rtf(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var header [5]byte
	if n, err := r.ReadAt(header[:], 0); (err != nil && err != io.EOF) || n < 5 {
		return false
	}
	expected := [5]byte{'{', '\\', 'r', 't', 'f'}
	if header != expected {
		return false
	}

	length := Length(r)
	const sanity = 6
	if length < sanity {
		return false
	}

	var tail [1]byte
	if n, err := r.ReadAt(tail[:], length-1); (err != nil && err != io.EOF) || n < 1 {
		return false
	}
	return tail[0] == '}'
}

// Txt returns true if the reader exclusively contains plain text ASCII characters,
// control characters or "extended ASCII characters".
func Txt(r io.ReaderAt) bool {
	return TxtWithLogger(nil, r)
}

// Deprecated: use [TxtWithLogger] instead.
// The io.Writer is unused.
func TxtW(_ io.Writer, r io.ReaderAt) bool {
	return TxtWithLogger(nil, r)
}

// TxtWithLogger returns true if the reader exclusively contains plain text ASCII characters,
// control characters or "extended ASCII characters".
//
// There is a 2% threshold for non-plain text characters such as ASCII control characters
// which are not printable but often found in plain text files for 8-bit microcomputers.
func TxtWithLogger(sl *slog.Logger, r io.ReaderAt) bool {
	const msg = "txt "
	if sl == nil {
		sl = slog.New(slog.DiscardHandler)
	}
	if r == nil {
		return false
	}

	const chunkSize = 1024
	size := Length(r)
	if size == 0 {
		return true
	}

	var buf [chunkSize]byte
	count := 0
	sl.Debug(msg+" reading chunks", slog.Int("chunk size B", chunkSize), slog.Int64("size B", size))

	for off := int64(0); off < size; off += chunkSize {
		bytesToRead := chunkSize
		if rem := size - off; rem < int64(chunkSize) {
			bytesToRead = int(rem)
		}

		n, err := r.ReadAt(buf[:bytesToRead], off)
		if err != nil && err != io.EOF {
			sl.Debug(msg+" read at error", slog.Any("error", err))
			return false
		}

		for i := range n {
			if NotPlainText(buf[i]) {
				count++
				if !threshold(count, size) {
					sl.Debug(msg+" count is >= than the two percent threshold",
						slog.Int("count", count), slog.Int64("size B", size))
					return false
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	return threshold(count, size)
}

// If count is greater than 2% of the file size, then it is not plain text.
func threshold(count int, size int64) bool {
	const percentage = 0.02
	return float64(count)/float64(size) < percentage
}

// TxtLatin1 returns true if the reader exclusively contains plain text ISO/IEC 8859-1 characters,
// commonly known as the Latin-1 character set.
func TxtLatin1(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	size := Length(r)
	if size == 0 {
		return true
	}

	const chunkSize = 1024
	var buf [chunkSize]byte

	for off := int64(0); off < size; off += chunkSize {
		bytesToRead := int(min(int64(chunkSize), size-off))

		n, err := r.ReadAt(buf[:bytesToRead], off)
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
	if r == nil {
		return false
	}

	size := Length(r)
	if size == 0 {
		return true
	}

	const chunkSize = 1024
	var buf [chunkSize]byte

	for off := int64(0); off < size; off += chunkSize {
		bytesToRead := int(min(int64(chunkSize), size-off))

		n, err := r.ReadAt(buf[:bytesToRead], off)
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
	if r == nil {
		return false
	}

	var p [3]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 3 {
		return false
	}
	return p == [3]byte{0xef, 0xbb, 0xbf}
}

// Utf16 returns true if the reader beings with the UTF-16 Byte Order Mark signature.
func Utf16(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [2]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 2 {
		return false
	}
	return p == [2]byte{0xff, 0xfe} || p == [2]byte{0xfe, 0xff}
}

// Utf32 returns true if the reader beings with the UTF-32 Byte Order Mark signature.
func Utf32(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	return p == [4]byte{0xff, 0xfe, 0x0, 0x0} || p == [4]byte{0x0, 0x0, 0xfe, 0xff}
}

// XBin matches the eXtender BInary text format.
func XBin(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [5]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 5 {
		return false
	}
	return p == [5]byte{'X', 'B', 'I', 'N', 0x1a}
}

// DumpCGA uses heuristics to discover 4-bit, RGBI Video RAM screen dumps
// of colored text that also included some styling bits.
// While CGA was the most common IBM PC hardware to use the 4-bit palette,
// it was also used by Tandy's PCs, and the IBM PCjr.
//
// The dump uses simple byte pairs and is far more efficient for PCs
// of the 1980s than parsing ANSI encoded text.
//
// The first byte is an 8-bit character, usually representing a
// code point from IBM's Code Page 437.
//
// The second byte is the display attributes:
//   - 0-2 bits are the three bits for one of eight foreground colors
//   - bit 3 is the foreground intensity, often a "bright" toggle
//   - 4-6 bits are the three bits for the one of eight background colors
//   - bit 7 is the background intensity, often a "blink" toggle
//
// This test does two reads while looking for specific patterns:
//   - A CGA screen dump is always 80 columns and 25 rows.
//   - The reader is equal to or multiples of 4000 characters
//     or padded to 4KB (4096 characters).
//   - On the first read, it looks at every attribute byte to
//     confirm the blinking bit is always off,
//     leaving 128 of 256 possible bit combinations.
//   - If the first read doesn't pass, the second read assumes
//     the blink attribute is in use, and looks for an 1980s
//     era screen dump, and instead confirms every attribute byte is using
//     a black background.
//
// Hopefully, this is enough to stop the majority of false positives.
func DumpCGA(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	const (
		cols   = 80
		rows   = 25
		screen = cols * rows * 2

		egaRows = 43
		ega43   = cols * egaRows * 2

		padding = 96
		padded  = screen + padding

		attrib  = 2
		blackBG = 0x70
		blink   = 0x80
	)

	size := Length(r)
	if size == 0 {
		return false
	}

	offset, ok := cgaOffset(r, size, screen, padded, ega43)
	if !ok {
		return false
	}

	buf := make([]byte, screen)
	_, err := r.ReadAt(buf, offset)
	if err != nil && err != io.EOF {
		return false
	}

	blinkFree := true
	blackOnly := true
	allZeroAttr := true

	for i := 1; i < screen; i += attrib {
		attr := buf[i]
		if attr != 0 {
			allZeroAttr = false
		}
		if (attr & blink) != 0 {
			blinkFree = false
		}
		if (attr & blackBG) != 0 {
			blackOnly = false
		}
	}

	if allZeroAttr || asciiLike(buf) {
		return false
	}
	return blinkFree || blackOnly
}

func asciiLike(buf []byte) bool {
	hasNewline := false
	for _, b := range buf {
		switch {
		case b == '\n':
			hasNewline = true
		case b == '\t' || b == '\r' || (b >= 0x20 && b <= 0x7E):
		default:
			return false
		}
	}
	return hasNewline
}

func bsave(r io.ReaderAt, size int64) (int64, int64, bool) {
	const (
		header     = 7
		bsaveMagic = 0xFD
		segCGA     = 0xB800
		segMDA     = 0xB000
		shift      = 8
	)
	if size < header {
		return 0, 0, false
	}

	var hdr [header]byte
	n, err := r.ReadAt(hdr[:], 0)
	if (err != nil && err != io.EOF) || n < header || hdr[0] != bsaveMagic {
		return 0, 0, false
	}

	seg := uint16(hdr[1]) | uint16(hdr[2])<<shift
	payload := int64(uint16(hdr[5]) | uint16(hdr[6])<<shift)
	if seg != segCGA && seg != segMDA {
		return 0, 0, false
	}
	if size == header+payload {
		return header, payload, true
	}
	if size == header+payload+1 {
		var eof [1]byte
		n, err := r.ReadAt(eof[:], size-1)
		if (err == nil || err == io.EOF) && n == 1 && (eof[0] == 0x00 || eof[0] == 0x1A) {
			return header, payload, true
		}
	}
	return 0, 0, false
}

func cgaOffset(r io.ReaderAt, size, screen, padded, ega43 int64) (int64, bool) {
	offset, payload, ok := bsave(r, size)
	if ok && (payload%screen == 0 || payload%padded == 0 || payload == ega43) {
		return offset, true
	}

	if size%screen == 0 || size%padded == 0 || size == ega43 {
		return 0, true
	}
	return 0, false
}

// DumpCGAGraphics uses heuristics to discover 16KB, IBM CGA VRAM dumps.
//   - 320x200 4-color.
//   - 640x200 2-color.
func DumpCGAGraphics(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	const (
		size        = 16384
		bankSize    = 8192
		bankActive  = 8000
		minGapZeros = 160
	)

	offset, read, bsave, ok := cgaSpan(r)
	if !ok {
		return false
	}

	buf := make([]byte, read)
	n, err := r.ReadAt(buf, offset)
	if (err != nil && err != io.EOF) || int64(n) < read {
		return false
	}

	bank1End := min(read, int64(bankSize+bankActive))
	if allZeroBytes(buf[:bankActive]) || allZeroBytes(buf[bankSize:bank1End]) {
		return false
	}
	if bsave {
		return true
	}

	// Require at least one 192-byte gap.
	// gap0 or gap1 when full 16KB is present to be mostly 0x00 or 0x1A padding.
	gap0OK := gapCount(buf[bankActive:bankSize]) >= minGapZeros
	gap1OK := read == size && gapCount(buf[bankSize+bankActive:size]) >= minGapZeros
	return gap0OK || gap1OK
}

// cgaSpan returns the CGA graphics payload offset, byte length to read,
// whether the source is a BSAVE file, and whether a valid span was found.
func cgaSpan(r io.ReaderAt) (int64, int64, bool, bool) {
	const (
		size       = 16384
		bankSize   = 8192
		bankActive = 8000
		minBank1   = 7936
		comMaxSize = 65536
		jmpNear    = 0xE9
		jmpShort   = 0xEB
	)

	l := Length(r)
	if l == size {
		return 0, size, false, true
	}
	if l < size || l > comMaxSize {
		return 0, 0, false, false
	}

	off, payload, matched := bsave(r, l)
	if matched && payload == size {
		return off, size, true, true
	}

	off, matched = cgaB800(r, l, bankSize+minBank1)
	if matched {
		return off, min(int64(size), l-off), false, true
	}

	off, matched = cgaView(r, l, bankSize+bankActive)
	if matched {
		return off, min(int64(size), l-off), false, true
	}

	var first [1]byte
	n, err := r.ReadAt(first[:], 0)
	if (err != nil && err != io.EOF) || n < 1 ||
		(first[0] != jmpNear && first[0] != jmpShort) {
		return 0, 0, false, false
	}
	return l - size, size, false, true
}

func gapCount(b []byte) int {
	n := 0
	for _, v := range b {
		if v == 0 || v == 0x1A {
			n++
		}
	}
	return n
}

// cgaB800 finds `mov ax, 0xB800; mov es, ax` (B8 00 B8 8E C0) with a nearby
// `mov si, imm16` (BE lo hi) in the first 4KB of a DOS .COM file and returns the embedded
// payload offset (si - 0x0100).
func cgaB800(r io.ReaderAt, l, payload int64) (int64, bool) {
	const (
		comOrg       = 0x0100
		scanMax      = 4096
		windowBack   = 12
		windowAhead  = 10
		minMovSiLen  = 3
		opMovSiImm16 = 0xBE
		opMovAxImm16 = 0xB8
		opMovRegReg  = 0x8B
		modRmSiAx    = 0xF0
	)

	size := min(l, int64(scanMax))
	buf := make([]byte, size)
	n, err := r.ReadAt(buf, 0)
	if (err != nil && err != io.EOF) || n < 8 {
		return 0, false
	}

	// matches `mov ax, 0xB800; mov es, ax` (point ES at CGA video memory)
	pointsToCGA := func(i int) bool {
		return buf[i] == 0xB8 && buf[i+1] == 0x00 && buf[i+2] == 0xB8 &&
			buf[i+3] == 0x8E && buf[i+4] == 0xC0
	}

	// reads the 16-bit source address loaded into SI (`mov si, imm16` or `mov ax, imm16; mov si, ax`)
	readSourceAddr := func(j int) (int64, bool) {
		isDirect := buf[j] == opMovSiImm16
		isViaAX := j+4 < n && buf[j] == opMovAxImm16 && buf[j+3] == opMovRegReg && buf[j+4] == modRmSiAx
		if !isDirect && !isViaAX {
			return 0, false
		}
		return int64(uint16(buf[j+1]) | uint16(buf[j+2])<<8), true
	}

	for i := 0; i+4 < n; i++ {
		if !pointsToCGA(i) {
			continue
		}
		lo := max(0, i-windowBack)
		hi := min(n-minMovSiLen, i+windowAhead)
		for j := lo; j <= hi; j++ {
			if si, ok := readSourceAddr(j); ok && si >= comOrg && si-comOrg+payload <= l {
				return si - comOrg, true
			}
		}
	}
	return 0, false
}

// cgaView detects DOS .COM CGA mode-4 viewers (mov ax, 0004h; int 10h)
// that copy embedded pixel data from `mov si, imm16` into B800h.
func cgaView(r io.ReaderAt, l, payload int64) (int64, bool) {
	const (
		comOrg        = 0x0100
		headerScan    = 24
		opMovSiImm16  = 0xBE
		movAxInt10Len = 5
		movAlInt10Len = 4
	)
	var hdr [headerScan]byte
	n, err := r.ReadAt(hdr[:], 0)
	if (err != nil && err != io.EOF) || n < 16 {
		return 0, false
	}

	// setsCGAMode4 returns the byte length of a leading `int 10h` CGA 320x200 mode-4 setup.
	setsCGAMode4 := func() int {
		switch {
		case hdr[0] == 0xB8 && hdr[1] == 0x04 && hdr[2] == 0x00 && hdr[3] == 0xCD && hdr[4] == 0x10:
			return movAxInt10Len
		case hdr[0] == 0xB0 && hdr[1] == 0x04 && hdr[2] == 0xCD && hdr[3] == 0x10:
			return movAlInt10Len
		default:
			return 0
		}
	}
	// loadsCGASegment matches `mov ax, 0xB800` (loading the CGA video memory segment).
	loadsCGASegment := func(i int) bool {
		return hdr[i] == 0xB8 && hdr[i+1] == 0x00 && hdr[i+2] == 0xB8
	}
	// loadsSourceAddr matches `mov si, imm16` and returns the .COM payload offset.
	loadsSourceAddr := func(i int) (int64, bool) {
		if hdr[i] != opMovSiImm16 {
			return 0, false
		}
		si := int64(uint16(hdr[i+1]) | uint16(hdr[i+2])<<8)
		return si - comOrg, si >= comOrg
	}

	start := setsCGAMode4()
	if start == 0 {
		return 0, false
	}

	var hasB800 bool
	var siOff int64 = -1
	for i := start; i+2 < n; i++ {
		if loadsCGASegment(i) {
			hasB800 = true
		}
		if siOff < 0 {
			if off, ok := loadsSourceAddr(i); ok {
				siOff = off
			}
		}
	}
	if !hasB800 || siOff < 0 || siOff+payload > l {
		return 0, false
	}
	return siOff, true
}

func allZeroBytes(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

// DumpSPICGA detects Software Pirates, Inc. DOS .COM loaders containing an embedded
// "SPI" RLE-compressed CGA graphics stream.
//
//	Example: "SPI" + uint16 pairCount + pairCount * [val, repeatMinus1]).
func DumpSPICGA(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	const (
		spi     = "SPI"
		size    = 65536
		minSize = 16000
		maxSize = 16384
		header  = 5 // "SPI" + uint16 pairCount
	)

	l := Length(r)
	if l < header || l > size {
		return false
	}

	buf := make([]byte, l)
	if n, err := r.ReadAt(buf, 0); (err != nil && err != io.EOF) || int64(n) < l {
		return false
	}

	// Must start with a DOS .COM jump (0xE9/0xEB) or directly with "SPI".
	if buf[0] != 0xE9 && buf[0] != 0xEB && !bytes.HasPrefix(buf, []byte(spi)) {
		return false
	}

	validSPI := func(tail []byte) bool {
		if len(tail) < header {
			return false
		}
		pairCount := int(uint16(tail[3]) | uint16(tail[4])<<8)
		if header+2*pairCount != len(tail) {
			return false
		}
		decodedLen := 0
		for pos := header; pos < len(tail); pos += 2 {
			decodedLen += int(tail[pos+1]) + 1
		}
		return decodedLen >= minSize && decodedLen <= maxSize
	}

	for searchFrom := 0; ; {
		rel := bytes.Index(buf[searchFrom:], []byte(spi))
		if rel < 0 {
			return false
		}
		idx := searchFrom + rel
		if validSPI(buf[idx:]) {
			return true
		}
		searchFrom = idx + 1
	}
}
