package magicnumber

// Package file synthesismusic.go contains the functions that parse bytes as common synthesis and tracker music formats.

import (
	"bytes"
	"fmt"
	"io"
)

// Midi matches the Musical Instrument Digital Interface (MIDI) format.
func Midi(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	return p == [4]byte{'M', 'T', 'h', 'd'}
}

// MTM matches the MultiTracker music format.
func MTM(r io.ReaderAt) bool {
	return MusicMTM(r) != ""
}

// XM matches the eXtended Module tracked music format.
func XM(r io.ReaderAt) bool {
	return MusicXM(r) != ""
}

// IT matches the Impulse Tracker music format.
func IT(r io.ReaderAt) bool {
	return MusicIT(r) != ""
}

// MK matches the ProTracker MOD music format.
func MK(r io.ReaderAt) bool {
	return MusicMK(r) != ""
}

// MusicTracker returns the tracked music format in the byte slice and
// the name or title of the song if available.
// The tracked music formats include MultiTracker, Impulse Tracker,
// Extended Module, and 4 channel MODule music.
//
// [Modland] has a large collection of tracked music format documentation.
//
// [Modland]: https://ftp.modland.com/pub/documents/format_documentation/
func MusicTracker(r io.ReaderAt) string {
	if s := MusicMTM(r); s != "" {
		return s
	}
	if s := MusicIT(r); s != "" {
		return s
	}
	if s := MusicXM(r); s != "" {
		return s
	}
	if s := MusicMK(r); s != "" {
		return s
	}
	return ""
}

// MusicMTM returns the [MultiTracker] song or title in the byte slice if available.
// The MultiTracker format is a tracked music format created by the scene group Renaissance.
//
// [MultiTracker]: https://ftp.modland.com/pub/documents/format_documentation/MultiTracker%20(.mtm).txt
func MusicMTM(r io.ReaderAt) string {
	if r == nil {
		return ""
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return ""
	}
	if p != [4]byte{'M', 'T', 'M', 0x10} {
		return ""
	}

	const off = 4
	var s [20]byte
	if n, err := r.ReadAt(s[:], off); (err != nil && err != io.EOF) || n < 20 {
		return ""
	}

	title := bytes.TrimSpace(bytes.TrimRight(s[:], "\x00"))
	const match = "MultiTrack song"
	if len(title) > 0 {
		return fmt.Sprintf(`%s, "%s"`, match, title)
	}
	return match
}

// MusicIT returns the [Impulse Tracker] song or title in the byte slice if available.
// The Impulse Tracker format is a tracked music format created by Jeffrey Lim.
//
// [Impulse Tracker]: https://ftp.modland.com/pub/documents/format_documentation/Impulse%20Tracker%20v2.04%20(.it).html
func MusicIT(r io.ReaderAt) string {
	if r == nil {
		return ""
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return ""
	}
	if p != [4]byte{'I', 'M', 'P', 'M'} {
		return ""
	}

	const off = 4
	var s [26]byte
	if n, err := r.ReadAt(s[:], off); (err != nil && err != io.EOF) || n < 26 {
		return ""
	}

	title := bytes.TrimSpace(bytes.TrimRight(s[:], "\x00"))
	const match = "Impulse Tracker song"
	if len(title) > 0 {
		return fmt.Sprintf(`%s, "%s"`, match, title)
	}
	return match
}

// MusicXM returns the [eXtended Module] song or title in the byte slice if available.
// The XM format was originally used by FastTracker II (FT2) and later modified by other trackers.
//
// [eXtended Module]: https://ftp.modland.com/pub/documents/format_documentation/FastTracker%202%20v2.04%20(.xm).html
func MusicXM(r io.ReaderAt) string {
	if r == nil {
		return ""
	}

	var p [17]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 17 {
		return ""
	}
	const idText = "Extended Module: "
	if string(p[:]) != idText {
		return ""
	}

	const off = 17
	var s [20]byte
	if n, err := r.ReadAt(s[:], off); (err != nil && err != io.EOF) || n < 20 {
		return ""
	}

	title := bytes.TrimSpace(bytes.TrimRight(s[:], "\x00"))
	const match = "extended module tracked music"
	if len(title) > 0 {
		return fmt.Sprintf(`%s, "%s"`, match, title)
	}
	return match
}

// MusicMK returns the MOD song or title in the byte slice if available.
// The Soundtracker MOD format is a tracked music format created by Karsten Obarski on the Commodore Amiga.
// The original MOD format had no signatures. However, Mahoney & Kaktus later used a 'M.K.' signature
// with their sound samples and it got adopted as the de-facto signature for the MOD format.
//
// Common MOD formats include the original The Ultimate Soundtracker, Protracker, FastTracker II...
//
// [ProTracker]: https://ftp.modland.com/pub/documents/format_documentation/ProTracker%20v1.0%20(.mod).html
func MusicMK(r io.ReaderAt) string {
	if r == nil {
		return ""
	}

	var p [4]byte
	const off = 1080
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 4 {
		return ""
	}

	switch p {
	case [4]byte{'2', 'C', 'H', 'N'}:
		return modSong("ProTracker 2-channel song", r)
	case
		[4]byte{'M', '.', 'K', '.'},
		[4]byte{'M', '!', 'K', '!'},
		[4]byte{'4', 'C', 'H', 'N'},
		[4]byte{'F', 'L', 'T', '4'}:
		return modSong("ProTracker 4-channel song", r)
	case [4]byte{'6', 'C', 'H', 'N'}:
		return modSong("ProTracker 6-channel song", r)
	case
		[4]byte{'F', 'L', 'T', '8'},
		[4]byte{'O', 'C', 'T', 'A'},
		[4]byte{'8', 'C', 'H', 'N'}:
		return modSong("ProTracker 8-channel song", r)
	default:
		return ""
	}
}

func modSong(match string, r io.ReaderAt) string {
	if r == nil {
		return match
	}
	const off = 0
	var s [20]byte
	if n, err := r.ReadAt(s[:], off); (err != nil && err != io.EOF) || n < 20 {
		return match
	}

	title := bytes.TrimSpace(bytes.TrimRight(s[:], "\x00"))
	if len(title) > 0 {
		return fmt.Sprintf(`%s, "%s"`, match, title)
	}
	return match
}
