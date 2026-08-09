package magicnumber

// File file cdimage.go contains the file type signature for physical media disk image formats.

import (
	"encoding/binary"
	"io"
)

// Daa returns true if the reader contains the PowerISO DAA CD image signature.
func Daa(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 12 {
		return false
	}
	if p[0] != 'D' || p[1] != 'A' || p[2] != 'A' || p[3] != 0x00 {
		return false
	}
	return daa([4]byte(p[8:12]))
}

func daa(p [4]byte) bool {
	ver := binary.LittleEndian.Uint32(p[:])
	const (
		sanity0 = 10
		sanity1 = 99
		shift   = 8
		andVal  = 0xff
	)
	if major := ver >> shift; major > sanity0 {
		return false
	}
	minor := ver & andVal
	return minor <= sanity1
}

// ISO returns true if the reader contains the ISO 9660 CD-ROM file-system signature.
// To be accurate, it requires at least 36KB of data to be read.
func ISO(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	offsets := [4]int64{32769, 34817, 36865, 0}
	var p [5]byte

	for _, off := range offsets {
		n, err := r.ReadAt(p[:], off)
		if (err != nil && err != io.EOF) || n < 5 {
			continue // continue to the next descriptor
		}
		if p == [5]byte{'C', 'D', '0', '0', '1'} {
			return true
		}
	}

	return false
}

// Mdf returns true if the reader contains the Alcohol 120% MDF CD image signature.
func Mdf(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [16]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 16 {
		return false
	}
	if p == [16]byte{'M', 'E', 'D', 'I', 'A', ' ', 'D', 'E', 'S', 'C', 'R', 'I', 'P', 'T', 'O', 'R'} {
		return true
	}

	cdSyncHeader := [12]byte{
		0x00, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x00,
	}
	var p12 [12]byte
	copy(p12[:], p[:12]) // comparing a copy more performant
	return p12 == cdSyncHeader
}

// Nero returns true if the reader contains either a Nero ISO compilation or a Nero Disc Image.
func Nero(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	if n, err := r.ReadAt(p[:], 0); (err == nil || err == io.EOF) && n >= 12 {
		if p == [12]byte{0x0e, 'N', 'e', 'r', 'o', 'I', 'S', 'O', 'C', 'o', 'm', 'p'} {
			return true
		}
	}

	type sizer interface {
		Size() int64
	}
	const minimum = 12
	if s, ok := r.(sizer); ok {
		size := s.Size()
		if size >= minimum {
			var p [4]byte
			if n, err := r.ReadAt(p[:], size-minimum); (err == nil || err == io.EOF) && n >= 4 {
				return p == [4]byte{'N', 'E', 'R', '5'} || p == [4]byte{'N', 'E', 'R', 'O'}
			}
		}
	}

	return false
}

// Deprecated: use [Nero] instead.
func Nri(r io.ReaderAt) bool {
	return Nero(r)
}
