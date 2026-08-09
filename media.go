package magicnumber

// Package file media.go contains the functions that parse bytes as common image, digital audio and video formats.
// A number of these media containers could support multiple modes,
// such as audio only, audio+video, video only, static images, animated images, etc.

import (
	"io"
)

// AAC matches the Advanced Audio Coding audio format.
func AAC(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [3]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 3 {
		return false
	}

	const lead = 0xff
	if p[0] != lead {
		return false
	}
	// byte 1 upper 4 bits must be 0xF (0xF0)
	// layer bits (bits 1-2) must be 0 (0x06 mask == 0)
	const byte1 = 0xf0
	const mask = 0xf6
	return (p[1] & mask) == byte1
}

// Avi matches the Microsoft Audio Video Interleave video format.
func Avi(r io.ReaderAt) bool {
	if r == nil || !RIFF(r) {
		return false
	}

	const off = 8
	var p [4]byte
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	return p == [4]byte{'A', 'V', 'I', 0x20}
}

// Avif returns true if the reader contains an AV1 Image File Format (AVIF) container.
func Avif(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [32]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 12 {
		return false
	}

	const ftyp = "ftyp"
	if string(p[4:8]) != ftyp {
		return false
	}

	const avif = "avif"
	const avis = "avis"
	major := string(p[8:12])
	if major == avif || major == avis {
		return true
	}

	for i := 16; i+4 <= n; i += 4 {
		brand := string(p[i : i+4])
		if brand == avif || brand == avis {
			return true
		}
	}

	return false
}

// Bmp matches the BMP image format.
func Bmp(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [18]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 18 {
		return false
	}
	// check magic
	if p[0] != 'B' || p[1] != 'M' {
		return false
	}
	// che3ck for reserved bits
	if p[6] != 0 || p[7] != 0 || p[8] != 0 || p[9] != 0 {
		return false
	}

	// DIB header size (uint32 Little-Endian at offset 14)
	dibHeader := uint32(p[14]) | uint32(p[15])<<8 | uint32(p[16])<<16 | uint32(p[17])<<24 //nolint:mnd
	const (
		os2v1   = 12
		os2v2   = 16
		winv3   = 40
		bmpv2   = 52
		bmpv3   = 56
		os2full = 64
		bmpv4   = 108
		bmpv5   = 124
	)
	switch dibHeader {
	case os2v1, os2v2, winv3, bmpv2, bmpv3, os2full, bmpv4, bmpv5:
		return true
	default:
		return false
	}
}

// Flac matches the Free Lossless Audio Codec audio format.
func Flac(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [5]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 5 {
		return false
	}

	if p[0] != 'f' || p[1] != 'L' || p[2] != 'a' || p[3] != 'C' {
		return false
	}

	// byte 4: bit 0 is 'last block' flag
	// bits 1-7 are block type (0-6 are standard FLAC block types)
	const valid = 0x7f
	const types = 6
	blockType := p[4] & valid
	return blockType <= types
}

// Flv matches the Shockwave Flash Video format.
func Flv(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [9]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 9 {
		return false
	}
	if p[0] != 'F' || p[1] != 'L' || p[2] != 'V' {
		return false
	}
	const version = 0x01
	if p[3] != version {
		return false
	}
	n := uint32(p[5])<<24 | uint32(p[6])<<16 | uint32(p[7])<<8 | uint32(p[8])
	const header = 9
	return n == header
}

// Gif matches the image Graphics Interchange Format.
// There are two versions of the GIF format, GIF87a and GIF89a.
func Gif(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [6]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 6 {
		return false
	}
	const gif87a = "GIF87a"
	const gif89a = "GIF89a"
	s := string(p[:])
	return s == gif87a || s == gif89a
}

// Ico matches the Microsoft Icon image format.
func Ico(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [6]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 6 {
		return false
	}
	if p[0] != 0x00 || p[1] != 0x00 {
		return false
	}
	if p[2] != 0x01 || p[3] != 0x00 {
		return false
	}
	images := uint16(p[4]) | uint16(p[5])<<8 //nolint:mnd
	return images > 0
}

// Iff matches the Interchange File Format image.
// This is a generic wrapper format originally created by Electronic Arts for storing data in chunks.
func Iff(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 12 {
		return false
	}

	const (
		form = "FORM"
		list = "LIST"
		cat  = "CAT "
	)
	id := string(p[:4])
	if id != form && id != list && id != cat {
		return false
	}

	const (
		interleaved     = "ILBM"
		planar          = "PBM "
		amigaContiguous = "ACBM"
		deepImage       = "DEEP"
		rgbn            = "RGBN"
		rgb8            = "RGB8"
	)
	switch string(p[8:12]) {
	case
		interleaved,
		planar,
		amigaContiguous,
		deepImage,
		rgbn,
		rgb8:
		return true
	default:
		return false
	}
}

// Deprecated: use [Real] instead.
func Ivr(r io.ReaderAt) bool {
	return Real(r)
}

// Real matches the Real Player and Real Media video formats.
func Real(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); err != nil || n < 4 {
		return false
	}
	const (
		rmf = ".RMF"
		rec = ".REC"
		ivr = ".IVR"
	)
	s := string(p[:])
	return s == rmf || s == rec || s == ivr
}

// Jpeg matches the JPEG File Interchange Format v1 image.
func Jpeg(r io.ReaderAt) bool {
	return jpeg(r, true)
}

// JpegNoSuffix matches the JPEG File Interchange Format v1 image.
// However, it does not check the final bytes, making it more performant,
// but it is a less accurate method than [Jpeg].
func JpegNoSuffix(r io.ReaderAt) bool {
	return jpeg(r, false)
}

func jpeg(r io.ReaderAt, suffix bool) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 3 {
		return false
	}

	// JPEG files always start with SOI (0xFF, 0xD8) followed by the marker (0xFF)
	if p[0] != 0xFF || p[1] != 0xD8 || p[2] != 0xFF {
		return false
	}

	// optional, JFIF or Exif marker validation
	const marker = 11
	if n >= marker {
		if p[3] == 0xE0 && string(p[6:11]) != "JFIF\x00" {
			return false
		}
		if p[3] == 0xE1 && string(p[6:11]) != "Exif\x00" {
			return false
		}
	}

	if suffix {
		return checkJpegSuffix(r)
	}
	return true
}

func checkJpegSuffix(r io.ReaderAt) bool {
	// sizer attempts to get reader size if implemented by r
	type sizer interface {
		Size() int64
	}

	var size int64
	if s, ok := r.(sizer); ok {
		size = s.Size()
	} else {
		return true
	}

	const minimum = 4
	if size < minimum {
		return false
	}

	// read trailing 130 bytes to handle standard EOI, null padding, or SAUCE records
	var tail [130]byte
	readSize := min(size, int64(len(tail)))

	off := size - readSize
	n, err := r.ReadAt(tail[:readSize], off)
	if err != nil && err != io.EOF {
		return false
	}

	buf := tail[:n]
	// strip trailing null bytes efficiently without reading full file
	for len(buf) > 0 && buf[len(buf)-1] == 0x00 {
		buf = buf[:len(buf)-1]
	}

	const sanity = 2
	if len(buf) < sanity {
		return false
	}
	// check for the standard EOI marker (0xFF, 0xD9)
	if buf[len(buf)-2] == 0xFF && buf[len(buf)-1] == 0xD9 {
		return true
	}
	// check for SAUCE metadata trailer
	// the 128-byte SAUCE record precedes EOF
	if len(buf) >= 130 &&
		buf[len(buf)-128] == 'S' &&
		buf[len(buf)-127] == 'A' &&
		buf[len(buf)-126] == 'U' &&
		buf[len(buf)-125] == 'C' &&
		buf[len(buf)-124] == 'E' {
		// check for the standard EOI marker before the SAUCE record
		return buf[len(buf)-130] == 0xFF && buf[len(buf)-129] == 0xD9
	}

	return false
}

// Jpeg2000 matches the JPEG 2000 image format.
func Jpeg2000(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 2 {
		return false
	}

	// raw codestream (.j2c / .j2k)
	if p[0] == 0xFF && p[1] == 0x4F {
		return n < 4 || (p[2] == 0xFF && p[3] == 0x51)
	}
	// .jp2 container
	const jp2Container = 12
	if n < jp2Container {
		return false
	}
	return p == [12]byte{
		0x00, 0x00, 0x00, 0x0C,
		'j', 'P', ' ', ' ',
		0x0D, 0x0A, 0x87, 0x0A,
	}
}

// Ilbm matches the InterLeaved Bitmap image format.
// Created by Electronic Arts it conforms to the IFF standard.
func Ilbm(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 12 {
		return false
	}

	const (
		form        = "FORM"
		interleaved = "ILBM"
	)
	id := string(p[:4])
	if id != form {
		return false
	}
	return string(p[8:12]) == interleaved
}

// IffAnim matches the Amiga animation format.
// Created by Electronic Arts it conforms to the IFF standard.
func IffAnim(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 12 {
		return false
	}

	const (
		form = "FORM"
		list = "LIST"
		cat  = "CAT "
		anim = "ANIM"
	)
	id := string(p[:4])
	if id != form && id != list && id != cat {
		return false
	}
	return string(p[8:12]) == anim
}

// IffPBM matches the IFF Planar BitMap image format.
// This is likely used by the IBM PC edition of Deluxe Paint II Deluxe,
// otherwise known as Deluxe Paint II v3.
func IffPBM(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 12 {
		return false
	}

	const (
		form   = "FORM"
		cat    = "CAT "
		list   = "LIST"
		planar = "PBM "
	)
	id := string(p[:4])
	if id != form && id != list && id != cat {
		return false
	}
	return string(p[8:12]) == planar
}

// IlbmDecode reads the InterLeaved Bitmap or PBM image dimensions from the reader
// and returns the image width and image height.
//
// Zero values are returned if the dimensions cannot be found or there is an error.
func IlbmDecode(r io.ReaderAt) (int, int) {
	if r == nil {
		return 0, 0
	}

	var p [256]byte
	n, err := r.ReadAt(p[:], 0)
	if err != nil || n < 12 {
		return 0, 0
	}

	const (
		form = "FORM"
		cat  = "CAT "
		list = "LIST"
		bmhd = "BMHD"
	)
	id := string(p[:4])
	if id != form && id != list && id != cat {
		return 0, 0
	}

	const offset = 12
	off := int64(offset)
	for off+8 <= int64(n) {
		const offset = 8
		chunkID := string(p[off : off+4])
		chunkSize := int(uint32(p[off+4])<<24 | uint32(p[off+5])<<16 | uint32(p[off+6])<<8 | uint32(p[off+7]))
		if chunkID == bmhd {
			dataOff := off + offset
			// BMHD chunk must have at least 4 bytes for width and height values
			if dataOff+4 > int64(n) {
				return 0, 0
			}
			w := uint16(p[dataOff])<<offset | uint16(p[dataOff+1])
			h := uint16(p[dataOff+2])<<offset | uint16(p[dataOff+3])
			return int(w), int(h)
		}
		// advance past the chunk header and chunk data
		off += offset + int64(chunkSize)
		if chunkSize%2 != 0 {
			off++
		}
	}

	return 0, 0
}

// M4v matches the QuickTime M4V video format.
func M4v(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [32]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 12 {
		return false
	}
	const ftyp = "ftyp"
	if string(p[4:8]) != ftyp {
		return false
	}
	const (
		m4v  = "M4V "
		m4vh = "M4VH"
		m4vp = "M4VP"
	)
	major := string(p[8:12])
	switch major {
	case m4v, m4vh, m4vp:
		return true
	}
	for off := 16; off+4 <= n; off += 4 {
		if string(p[off:off+4]) == m4v {
			return true
		}
	}
	return false
}

// Mp4 matches the MPEG-4 video format.
func Mp4(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [32]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 12 {
		return false
	}
	const ftyp = "ftyp"
	if string(p[4:8]) != ftyp {
		return false
	}
	if isMP4([4]byte(p[8:12])) {
		return true
	}
	for off := 16; off+4 <= n; off += 4 {
		if isMP4([4]byte(p[off : off+4])) {
			return true
		}
	}

	return false
}

// isMP4 returns true if the 4-byte brand code identifies a MP4 container.
func isMP4(brand [4]byte) bool {
	const (
		iso14496_14 = true
		iso14496_16 = true
		neroDigital = true
		sonyMP4     = true
		adobeMP4    = true
	)
	switch string(brand[:]) {
	case "isom", "iso2", "iso3", "iso4", "iso5", "iso6":
		return iso14496_14
	case "mp41", "mp42", "mp71", "avc1":
		return iso14496_16
	case "ndas", "ndsc", "ndsh", "ndsm", "ndsp", "ndss", "ndxc", "ndxh", "ndxm", "ndxp", "ndxs":
		return neroDigital
	case "MSNV":
		return sonyMP4
	case "f4v ", "f4p ":
		return adobeMP4
	default:
		return false
	}
}

// Mp3 matches the MPEG-1 or 2, Audio Layer III (MP3) format.
// It checks for either an ID3v2 container header or an MPEG Audio sync frame.
func Mp3(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [10]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 2 {
		return false
	}
	// check for the common ID3 v2 header
	if n >= 4 && p[0] == 'I' && p[1] == 'D' && p[2] == '3' {
		const sanity = 0xff
		if p[3] < sanity {
			return true
		}
	}
	return isMP3(p[:n])
}

// isMP3 checks if the buffer starts with a valid MPEG Audio Frame sync word.
func isMP3(p []byte) bool {
	const sanity = 2
	if len(p) < sanity {
		return false
	}

	const byte0 = 0xff
	if p[0] != byte0 {
		return false
	}

	// byte 1's first 3 bits must be '111' (sync bits 11..8)
	const sync = 0xe0
	if (p[1] & sync) != sync {
		return false
	}

	// byte 1, extract audio layer (bits 2..1)
	// 0b11 = Layer I, 0b10 = Layer II, 0b01 = Layer III (MP3), 0b00 = Reserved
	const byte1 = 0x03
	const mp3 = 0x01
	layer := (p[1] >> 1) & byte1
	return layer == mp3
}

// Mpeg matches MPEG-1 and MPEG-2 video formats.
func Mpeg(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	if p[0] != 0x00 || p[1] != 0x00 || p[2] != 0x01 {
		return false
	}
	// check for valid mpeg start codes
	code := p[3]
	return code == 0xb3 || code == 0x00 || (code >= 0xba && code <= 0xbf)
}

// Ogg matches the Ogg Vorbis audio format.
func Ogg(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [6]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 6 {
		return false
	}

	if p[0] != 'O' || p[1] != 'g' || p[2] != 'g' || p[3] != 'S' {
		return false
	}
	const version = 0x00
	if p[4] != version {
		return false
	}
	const bos = 0x02 // beginning of stream
	return (p[5] & bos) != 0
}

// Pcx matches the ZSoft Personal Computer eXchange (PCX) image format.
func Pcx(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 4 {
		return false
	}

	const manufacturer = 0x0a
	if p[0] != manufacturer {
		return false
	}
	version := p[1]
	if version > 5 || version == 1 {
		return false
	}
	const uncompressed = 0x00
	const compressed = 0x01
	encoding := p[2]
	if encoding != uncompressed && encoding != compressed {
		return false
	}
	bpp := p[3] // bits per pixel
	return bpp == 1 || bpp == 2 || bpp == 4 || bpp == 8
}

// Png matches the Portable Network Graphics image format.
func Png(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [8]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 8 {
		return false
	}
	return p == [8]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}
}

// QTMov matches the QuickTime Movie video format.
func QTMov(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 8 {
		return false
	}

	const qt = "qt  "
	atomType := string(p[4:8])
	switch atomType {
	case "moov", "mdat", "wide", "free", "skip":
		return true
	case "ftyp":
		return n >= 12 && string(p[8:12]) == qt
	default:
		return false
	}
}

// RIFF returns true if the first 4 bytes in the reader match 'RIFF'.
func RIFF(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	return p == [4]byte{'R', 'I', 'F', 'F'}
}

// Ripscrip returns true if the reader contains the RIPscrip signature.
// It is a vector graphics format used in BBS systems in the early 1990s.
func Ripscrip(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [3]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 3 {
		return false
	}

	if p[0] != '!' || p[1] != '|' {
		return false
	}
	digit := p[2]
	return digit >= '0' && digit <= '9'
}

// Tiff matches the Tagged Image File Format.
// Conforms to Aldus/Adobe TIFF revision 6.0 and BigTIFF specifications.
func Tiff(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [4]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 4 {
		return false
	}
	littleEndian := p[0] == 'I' && p[1] == 'I'
	if littleEndian {
		// v42 (standard) or v43 (BigTIFF)
		return (p[2] == 0x2A && p[3] == 0x00) || (p[2] == 0x2B && p[3] == 0x00)
	}
	bigEndian := p[0] == 'M' && p[1] == 'M'
	if bigEndian {
		// v42 (standard) or v43 (BigTIFF)
		return (p[2] == 0x00 && p[3] == 0x2A) || (p[2] == 0x00 && p[3] == 0x2B)
	}

	return false
}

// Wave matches the IBM / Microsoft Waveform audio format.
func Wave(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [12]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 12 {
		return false
	}

	if p[0] != 'R' || p[1] != 'I' || p[2] != 'F' {
		return false
	}
	if p[3] != 'F' && p[3] != 'X' {
		return false
	}
	return p[8] == 'W' && p[9] == 'A' && p[10] == 'V' && p[11] == 'E'
}

// Webp matches the Google WebP image format.
func Webp(r io.ReaderAt) bool {
	if r == nil {
		return false
	}
	var p [12]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 12 {
		return false
	}
	if p[0] != 'R' || p[1] != 'I' || p[2] != 'F' || p[3] != 'F' {
		return false
	}
	return p[8] == 'W' && p[9] == 'E' && p[10] == 'B' && p[11] == 'P'
}

// Wmv matches the Microsoft Windows Media video format.
func Wmv(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [16]byte
	n, err := r.ReadAt(p[:], 0)
	if (err != nil && err != io.EOF) || n < 16 {
		return false
	}
	return p == [16]byte{
		0x30, 0x26, 0xb2, 0x75, 0x8e, 0x66, 0xcf, 0x11,
		0xa6, 0xd9, 0x00, 0xaa, 0x00, 0x62, 0xce, 0x6c,
	}
}
