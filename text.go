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

	const chunkSize = 1024
	p := make([]byte, chunkSize)

	finds := 0
	var off int64

	// track up to 2 trailing bytes from the previous chunk to detect boundary-straddling CSI sequences
	var tail [2]byte
	tails := 0

	const minRequired = 3
	for {
		n, err := r.ReadAt(p, off)
		if err != nil && err != io.EOF {
			return false
		}
		if n == 0 {
			break
		}

		// prepare a search slice that includes overlap from the previous chunk
		var searchBuf []byte
		if tails > 0 {
			searchBuf = append(tail[:tails], p[:n]...)
		} else {
			searchBuf = p[:n]
		}

		// count occurrences of CSI codes (\x1b[ + char) in this chunk
		// common is a list of ANSI CSI termination/parameter characters
		const common = `0123456789JK=su#`
		for _, c := range []byte(common) {
			pattern := [3]byte{esc, '[', c}
			finds += bytes.Count(searchBuf, pattern[:])
			if finds >= minRequired {
				return true
			}
		}

		// save the last 2 bytes for the next boundary check
		const last2 = 2
		if n >= last2 {
			tail[0] = p[n-2]
			tail[1] = p[n-1]
			tails = 2
		} else if n == 1 {
			tail[0] = p[n-1]
			tails = 1
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
