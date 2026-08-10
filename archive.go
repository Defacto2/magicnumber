package magicnumber

// Package file archive.go contains the functions that parse bytes as common file archive,
// compression and disk image formats.

import (
	"io"
)

// Zip64 returns true if the reader contains a valid PKWARE Zip64 archive.
//
// This is an extension to the original ZIP format that allows for larger archives.
func Zip64(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	if p != [4]byte{'P', 'K', 0x03, 0x04} {
		return false
	}

	const minimum = 32
	size := Length(r)
	if size < minimum {
		return false
	}

	var tail [1024]byte
	read := int(min(size, int64(len(tail))))
	off := size - int64(read)
	n, err := r.ReadAt(tail[:read], off)
	if (err != nil && err != io.EOF) || n < 20 {
		return false
	}

	// find the End of Central Directory Record or Central Directory Locator
	buf := tail[:n]
	for i := 0; i <= len(buf)-4; i++ {
		if buf[i] == 'P' && buf[i+1] == 'K' && buf[i+2] == 0x06 {
			if buf[i+3] == 0x06 || buf[i+3] == 0x07 {
				return true
			}
		}
	}

	return false
}

// Pkzip matches the zip archive format.
func Pkzip(r io.ReaderAt) bool {
	return pkzip(r) == pkZip
}

// PkImplode matches the PKWARE Implode method zip archive format.
// This is a legacy method and is generally not supported in modern ZIP tools and libraries.
func PkImplode(r io.ReaderAt) bool {
	return pkzip(r) == pkImplode
}

// PkReduce matches the PKWARE Reduce method zip archive format.
// This is a legacy method and is generally not supported in modern ZIP tools and libraries.
func PkReduce(r io.ReaderAt) bool {
	return pkzip(r) == pkReduce
}

// PkShrink matches the PKWARE Shrink method zip archive format.
// This is a legacy method and is generally not supported in modern ZIP tools and libraries.
func PkShrink(r io.ReaderAt) bool {
	return pkzip(r) == pkSkrink
}

// PkzipMulti matches the PKWARE Multi-Volume Zip archive format.
func PkzipMulti(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	return p == [4]byte{'P', 'K', 0x07, 0x08}
}

type pkComp int

const (
	pkNone pkComp = iota
	pkZip
	pkSkrink
	pkReduce
	pkImplode
)

// pkzip matches the PKWARE Zip archive format.
// This is the most common ZIP format and is widely supported and has been
// tested against many discontinued and legacy ZIP methods and packagers.
//
// Due to the complex history of the ZIP format, 4 possible return values
// maybe returned.
//   - pkNone is returned if the file is not a ZIP archive.
//   - pkOkay is returned if the file is a ZIP archive, except for the compression methods below.
//   - pkSkrink is returned if the ZIP archive uses the PKWARE shrink method, found in PKZIP v0.9.
//   - pkReduce is returned if the ZIP archive uses the PKWARE reduction method, found in PKZIP v0.8.
//   - pkImplode is returned if the ZIP archive uses the PKWARE implode method, found in PKZIP v1.01.
//
// Compression methods Shrink, Reduce and Implode are legacy and are generally
// not supported in modern ZIP tools and libraries.
func pkzip(r io.ReaderAt) pkComp {
	if r == nil {
		return pkNone
	}

	var p [30]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 30 {
		return pkNone
	}
	// local file header signature
	if p[0] != 'P' || p[1] != 'K' || p[2] != 0x03 || p[3] != 0x04 {
		return pkNone
	}

	// the compression method is a 16-bit little-endian integer at offsets 8 and 9
	const shift = 8
	method := uint16(p[8]) | uint16(p[9])<<shift
	return pkMethod(method)
}

func pkMethod(n uint16) pkComp {
	const (
		store       uint16 = 0x0
		shrink      uint16 = 0x1
		reduce1     uint16 = 0x2
		reduce2     uint16 = 0x3
		reduce3     uint16 = 0x4
		reduce4     uint16 = 0x5
		implode     uint16 = 0x6
		deflate     uint16 = 0x8
		deflate64   uint16 = 0x9
		ibmTerse    uint16 = 0xa
		bzip2       uint16 = 0xc
		lzma        uint16 = 0xe
		ibmCMPSC    uint16 = 0x10
		ibmTerseNew uint16 = 0x12
		ibmLZ77z    uint16 = 0x13
		zstd        uint16 = 0x5d
		mp3         uint16 = 0x5e
		xz          uint16 = 0x5f
		jpeg        uint16 = 0x60
		wavPack     uint16 = 0x61
		ppmd        uint16 = 0x62
		ae          uint16 = 0x63
	)

	switch n {
	case store, deflate, deflate64:
		return pkZip
	case shrink:
		return pkSkrink
	case reduce1, reduce2, reduce3, reduce4:
		return pkReduce
	case implode:
		return pkImplode
	case ibmTerse, bzip2, lzma, ibmCMPSC, ibmTerseNew, ibmLZ77z, zstd, mp3, xz, jpeg, wavPack, ppmd, ae:
		return pkZip
	default:
		return pkNone
	}
}

// Tar matches the Tape ARchive format.
func Tar(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	const off = 257
	var p [8]byte
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 8 {
		return false
	}

	if p[0] != 'u' || p[1] != 's' || p[2] != 't' || p[3] != 'a' || p[4] != 'r' {
		return false
	}

	const (
		posixTar = 0x00
		gnuTar   = ' '
	)
	return p[5] == posixTar || p[5] == gnuTar
}

// Rar matches the Roshal ARchive format, RAR v1 to RAR v4.
func Rar(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [7]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 7 {
		return false
	}

	return p == [7]byte{'R', 'a', 'r', 0x21, 0x1a, 0x07, 0x00}
}

// Rarv5 matches the Roshal ARchive format, RAR v5.
func Rarv5(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [8]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 8 {
		return false
	}

	return p == [8]byte{'R', 'a', 'r', 0x21, 0x1a, 0x07, 0x01, 0x00}
}

// Gzip matches the Gzip Compress archive format.
func Gzip(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [3]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 3 {
		return false
	}

	return p[0] == 0x1f && p[1] == 0x8b && p[2] == 0x08
}

// Bzip2 matches the Bzip2 Compress archive format.
func Bzip2(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	// signature 'B', 'Z', 'h' plus the block size digit '1'-'9'
	return p[0] == 'B' && p[1] == 'Z' && p[2] == 'h' && p[3] >= '1' && p[3] <= '9'
}

// X7z matches the 7z Compress archive format.
func X7z(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [6]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 6 {
		return false
	}

	return p == [6]byte{'7', 'z', 0xbc, 0xaf, 0x27, 0x1c}
}

// XZ matches the XZ Compress archive format.
func XZ(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [6]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 6 {
		return false
	}

	return p == [6]byte{0xfd, '7', 'z', 'X', 'Z', 0x00}
}

// ZStd matches the ZStandard archive format.
func ZStd(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	standard := [4]byte{0x28, 0xb5, 0x2f, 0xfd}
	if p == standard {
		return true
	}

	skippable := (p[0] >= 0x50 && p[0] <= 0x5f) && p[1] == 0x2a && p[2] == 0x4d && p[3] == 0x18 //nolint:mnd
	return skippable
}

// ArcFree matches the FreeArc compression format.
func ArcFree(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	return p == [4]byte{'A', 'r', 'C', 0x01}
}

// ArcSEA matches the ARChive SEA compression format.
func ArcSEA(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [2]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 2 {
		return false
	}

	const (
		id        = 0x1a
		maxMethod = 0x11 // max method ID for traditional ARC compression formats
	)
	// must use a valid method byte (1 through 17)
	return p[0] == id && p[1] > 0 && p[1] <= maxMethod
}

// LzhLha matches the LHA and LZH compression formats.
func LzhLha(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [5]byte
	const off = 2
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 5 {
		return false
	}

	// must begin and end with '-'
	if p[0] != '-' || p[4] != '-' {
		return false
	}

	// must be either "-lh#-" or "-lz#-"
	return p[1] == 'l' && (p[2] == 'h' || p[2] == 'z')
}

// Zoo matches the Zoo compression format.
func Zoo(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	const off = 20
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	return p == [4]byte{0xdc, 0xa7, 0xc4, 0xfd}
}

// Arj matches ARJ compression format.
func Arj(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [11]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 11 {
		return false
	}

	if p[0] != 0x60 || p[1] != 0xea {
		return false
	}

	const shift = 8
	headerSize := uint16(p[2]) | (uint16(p[3]) << shift)
	if headerSize < 15 || headerSize > 2600 {
		return false
	}

	const (
		msdos     = 0
		primeComp = 1
		unix      = 2
		amiga     = 3
		macintosh = 4
		ibmOS2    = 5
		apple2    = 6
		atariST   = 7
		next      = 8
		vax       = 9
		win95     = 10
		win32     = 11
	)
	switch p[6] {
	case msdos, primeComp, unix, amiga, macintosh, ibmOS2, apple2, atariST, next, vax, win95, win32:
		return true
	default:
		return false
	}
}

// Cab matches the Microsoft CABinet archive format.
func Cab(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [28]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 28 {
		return false
	}

	return p[0] == 'M' && p[1] == 'S' && p[2] == 'C' && p[3] == 'F' &&
		p[24] == 0x03 && p[25] == 0x01
}

// Pak matches the NoGate Consulting PAK format.
func Pak(r io.ReaderAt) bool {
	size := Length(r)
	const minimum = 4
	if size < minimum {
		return false
	}

	var p [2]byte
	if _, err := r.ReadAt(p[:], 0); err != nil {
		return false
	}

	const arcMarker = 0x1a
	if p[0] != arcMarker {
		return false
	}

	const crushed = 0x0A
	const distilled = 0x0B
	methodNoGate := p[1] == crushed || p[1] == distilled

	const two = 2
	off := size - two
	if _, err := r.ReadAt(p[:], off); err != nil {
		return false
	}

	const (
		nul = 0x00
		arc = 0x1A
		pak = 0xFE
	)
	eofARC := p[0] == arc && p[1] == nul
	eofPAK := p[0] == pak && p[1] == nul
	if methodNoGate && (eofARC || eofPAK) {
		return true
	}
	return p[0] == pak && p[1] == nul
}
