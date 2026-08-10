package magicnumber

// Package file id3.go contains the functions that parse bytes as common ID3 tag formats usually found in MP3 files.

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// ID3v1Size is the minimum buffer size of an ID3 v1 tag.
const ID3v1Size = 128

// MusicID3v1 reads the [ID3 v1] tag in the byte slice and returns "Song by Artist (Year)".
// The ID3 v1 tag is a 128 byte tag at the end of an MP3 audio file.
//
// [ID3 v1]: http://id3.org/ID3v1
func MusicID3v1(r io.ReaderAt) string {
	if r == nil {
		return ""
	}

	length := Length(r)
	if length < ID3v1Size {
		return ""
	}

	off := length - ID3v1Size
	var p [128]byte
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < ID3v1Size {
		return ""
	}

	if p[0] != 'T' || p[1] != 'A' || p[2] != 'G' {
		return ""
	}

	trimReader := func(b []byte) string {
		// cut at first null byte if present
		if idx := bytes.IndexByte(b, 0); idx != -1 {
			b = b[:idx]
		}
		return strings.TrimSpace(string(b))
	}
	song := trimReader(p[3:33])
	artist := trimReader(p[33:63])
	year := trimReader(p[93:97])
	return formatID3(song, artist, year)
}

// MusicID3v2 reads the ID3 v2 tag at the start of an MP3 audio file and returns "Song by Artist (Year)".
//
// [ID3 v2]: https://id3.org/id3v2-00
func MusicID3v2(r io.ReaderAt) string {
	if r == nil {
		return ""
	}

	var p [10]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 10 {
		return ""
	}

	if p[0] != 'I' || p[1] != 'D' || p[2] != '3' {
		return ""
	}

	version := p[3]
	if version < 2 || version > 4 {
		return ""
	}

	size := Syncsafe(p[6:10])
	if size <= 0 {
		return ""
	}

	// read the tag body (bytes following the 10-byte header)
	const off = 10
	b := make([]byte, size)
	if n, err := r.ReadAt(b, off); (err != nil && err != io.EOF) || int64(n) < size {
		return ""
	}

	const v22 = 2
	if version == v22 {
		return parseID3v22(b)
	}
	return parseID3(b, version)
}

// Syncsafe converts the p value into a 7-bit, synchsafe integer.
func Syncsafe(p []byte) int64 {
	const bitsPerByte = 7
	var size int64
	for _, b := range p {
		if b&0x80 != 0 {
			return 0
		}
		size = (size << bitsPerByte) | int64(b)
	}
	return size
}

// Sequential frame parser.
func frames(p []byte, version byte) map[string]string {
	frames := make(map[string]string)
	offset := 0

	for offset+10 <= len(p) {
		id := string(p[offset : offset+4])
		const padding = 0
		if id[0] == padding {
			break
		}

		var size int
		const v24 = 4
		if version == v24 {
			size = int(Syncsafe(p[offset+4 : offset+8]))
		} else {
			size = int(p[offset+4])<<24 | int(p[offset+5])<<16 | int(p[offset+6])<<8 | int(p[offset+7])
		}

		// skip 10-byte frame header
		offset += 10
		if size <= 0 || offset+size > len(p) {
			break
		}

		frames[id] = frame(p[offset : offset+size])
		offset += size
	}

	return frames
}

// Sequential 3-byte, frame parser for the older ID3v2.2 standard.
func frames3B(p []byte) map[string]string {
	frames := make(map[string]string)
	offset := 0

	for offset+6 <= len(p) {
		id := string(p[offset : offset+3])
		const padding = 0
		if id[0] == padding {
			break
		}

		size := int(p[offset+3])<<16 | int(p[offset+4])<<8 | int(p[offset+5])
		offset += 6

		if size <= 0 || offset+size > len(p) {
			break
		}

		frames[id] = frame(p[offset : offset+size])
		offset += size
	}

	return frames
}

func frame(p []byte) string {
	const minimum = 2
	if len(p) < minimum {
		return ""
	}

	const (
		iso8859_1 = 0
		cutset    = "\x00"
	)

	body := bytes.Trim(p[1:], cutset)

	if encoding := p[0]; encoding == iso8859_1 {
		v, err := Latin1(body)
		if err != nil {
			return string(body)
		}
		return v
	}
	return string(body)
}

func parseID3(p []byte, version byte) string {
	frames := frames(p, version)

	const (
		tit1 = "TIT1"
		tit2 = "TIT2"
		talb = "TALB"
		tpe1 = "TPE1"
		tyer = "TYER"
		tdrc = "TDRC"
	)

	song := frames[tit2]
	if song == "" {
		song = frames[talb] // fallback to album title if song title missing
	}
	if song == "" {
		return ""
	}

	artist := frames[tpe1]
	if artist == "" {
		artist = frames[tit1]
	}

	year := frames[tyer]
	if year == "" {
		// ID3v2.4 replaces TYER with TDRC (Recording Time)
		year = frames[tdrc]
	}
	const digits = 4
	if len(year) >= digits {
		year = year[:4]
	}
	if _, err := strconv.Atoi(year); err != nil {
		year = ""
	}

	return formatID3(song, artist, year)
}

func parseID3v22(payload []byte) string {
	frames := frames3B(payload)

	song := frames["TT2"]
	if song == "" {
		song = frames["TAL"] // fallback to album title if song title missing
	}
	if song == "" {
		return ""
	}

	artist := frames["TP1"]
	if artist == "" {
		artist = frames["TP2"]
	}

	year := frames["TYE"]
	if year != "" {
		if _, err := strconv.Atoi(year); err != nil {
			year = ""
		}
	}

	return formatID3(song, artist, year)
}

func formatID3(song, artist, year string) string {
	var sb strings.Builder
	sb.WriteString(song)

	if artist != "" {
		sb.WriteString(" by ")
		sb.WriteString(artist)
	}
	if year != "" {
		sb.WriteString(" (")
		sb.WriteString(year)
		sb.WriteString(")")
	}

	return sb.String()
}

// Latin1 converts a byte slice to a Latin-1 (ISO-8859-1) string.
func Latin1(b []byte) (string, error) {
	decoder := charmap.ISO8859_1.NewDecoder()
	s, err := decoder.Bytes(b)
	if err != nil {
		const format = "magicnumber iso 8859-1 decoder: %w"
		return "", fmt.Errorf(format, err)
	}
	return string(s), nil
}

// Deprecated: use [Latin1] instead.
func ConvLatin1(b []byte) (string, error) {
	return Latin1(b)
}

// Deprecated: use [Syncsafe] instead.
func ConvSize(b []byte) int64 {
	return Syncsafe(b)
}
