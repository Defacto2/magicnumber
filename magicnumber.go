// Package magicnumber contains the magic number matchers for identifying file types that
// are expected to be handled by the Defacto2 server application. Magic numbers are not
// always accurate and should be used as hints combined with other checks such as
// file extension matching.
//
// Usually, the magic number is the first few bytes of a file that uniquely identify the file type.
// But a number of document formats also check the final few bytes of a file.
//
// The sources for the magic numbers byte values are from the following:
//   - [Gary Kessler's File Signatures Table]
//   - [Just Solve the File Format Problem]
//   - [OSDev Wiki]
//   - [Wikipedia]
//
// [Gary Kessler's File Signatures Table]: https://www.garykessler.net/library/file_sigs.html
// [Just Solve the File Format Problem]: http://fileformats.archiveteam.org/wiki/Electronic_File_Formats
// [OSDev Wiki]: https://wiki.osdev.org]
// [Wikipedia]: https://en.wikipedia.org/wiki/List_of_file_signatures
package magicnumber

import (
	"errors"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

var ErrNilReader = errors.New("nil reader")

// Signature represents a file type signature.
type Signature int

const (
	ZeroByte Signature = iota - 2
	Unknown
	ElectronicArtsIFF
	AV1ImageFile
	JPEGFileInterchangeFormat
	JPEG2000
	PortableNetworkGraphics
	GraphicsInterchangeFormat
	GoogleWebP
	TaggedImageFileFormat
	BMPFileFormat
	PersonalComputereXchange
	InterleavedBitmap
	MicrosoftIcon
	RIPscrip
	MPEG4
	QuickTimeMovie
	QuickTimeM4V
	MicrosoftAudioVideoInterleave
	MicrosoftWindowsMedia
	MPEG
	FlashVideo
	RealPlayer
	MusicalInstrumentDigitalInterface
	MPEG1AudioLayer3
	MPEGAdvancedAudioCoding
	OggVorbisCodec
	FreeLosslessAudioCodec
	WaveAudioForWindows
	MusicExtendedModule
	MusicMultiTrackModule
	MusicImpulseTracker
	MusicProTracker
	PKWAREZipShrink
	PKWAREZipReduce
	PKWAREZipImplode
	PKWAREZip64
	PKWAREZip
	PKWAREMultiVolume
	PKLITE
	PKSFX
	TapeARchive
	RoshalARchive
	RoshalARchivev5
	GzipCompressArchive
	Bzip2CompressArchive
	X7zCompressArchive
	XZCompressArchive
	ZStandardArchive
	FreeArc
	ARChiveSEA
	YoshiLHA
	ZooArchive
	ArchiveRobertJung
	MicrosoftCABinet
	MicrosoftDOSKWAJ
	MicrosoftDOSSZDD
	MicrosoftExecutable
	MicrosoftCompoundFile
	CDISO9660
	CDNero
	CDPowerISO
	CDAlcohol120
	WindowsHelpFile
	PortableDocumentFormat
	RichTextFormat
	UTF8Text
	UTF16Text
	UTF32Text
	ANSIEscapeText
	PlainText
)

func (sign Signature) String() string { //nolint:funlen
	switch {
	case sign <= ZeroByte:
		return "0-byte data"
	case sign == Unknown:
		return "binary data"
	case sign > PlainText:
		return "error"
	}
	return [...]string{
		"IFF image",
		"AV1 image",
		"JPEG image",
		"JPEG 2000 image",
		"PNG image",
		"GIF image",
		"WebP image",
		"TIFF image",
		"BMP image",
		"PCX image",
		"BMP image",
		"Microsoft icon",
		"RIPscrip",
		"MPEG-4 video",
		"QuickTime video",
		"QuickTime video",
		"AVI video",
		"Windows Media video",
		"MPEG-4 video",
		"Flash video",
		"RealPlayer video",
		"MIDI audio",
		"MP3 audio",
		"ACC audio",
		"Ogg audio",
		"FLAC audio",
		"Wave audio",
		"Tracker music extended mod",
		"Tracker music multi-track mod",
		"Tracker music Impulse mod",
		"Tracker music ProTracker mod",
		"pkzip shrunk archive",
		"pkzip reduced archive",
		"pkzip imploded archive",
		"zip64 archive",
		"zip archive",
		"multivolume zip",
		"pklite compressed",
		"self-extracting zip",
		"Tape archive",
		"RAR archive",
		"RAR v5+ archive",
		"Gzip archive",
		"Bzip2 archive",
		"7z archive",
		"XZ archive",
		"ZST archive",
		"FreeARC",
		"ARC by SEA",
		"LHA by Yoshi",
		"Zoo archive",
		"ARJ archive",
		"Microsoft cabinet",
		"MS-DOS KWAJ",
		"MS-DOS SZDD",
		"MS-DOS executable",
		"Microsoft compound fFile",
		"CD, ISO 9660",
		"CD, Nero",
		"CD, PowerISO",
		"CD, Alcohol 120",
		"Windows help",
		"PDF document",
		"rich text",
		"UTF-8 text",
		"UTF-16 text",
		"UTF-32 text",
		"ANSI text",
		"plain text",
	}[sign]
}

func (sign Signature) Title() string { //nolint:funlen
	switch {
	case sign <= ZeroByte:
		return "Zero-byte data"
	case sign == Unknown:
		return "Binary data"
	case sign > PlainText:
		return "Error"
	}
	return [...]string{
		"Electronic Arts IFF",
		"AV1 Image File",
		"JPEG File Interchange Format",
		"JPEG 2000",
		"Portable Network Graphics",
		"Graphics Interchange Format",
		"Google WebP",
		"Tagged Image File Format",
		"Bitmap image file",
		"Personal Computer eXchange",
		"Interleaved Bitmap",
		"Microsoft Icon",
		"RIPscrip vector graphic",
		"MPEG-4 video",
		"QuickTime Movie",
		"QuickTime M4V",
		"Microsoft Audio Video Interleave",
		"Microsoft Windows Media",
		"MPEG-4 video",
		"Flash Video",
		"RealPlayer",
		"Musical Instrument Digital Interface",
		"MPEG-1 Audio Layer 3",
		"MPEG Advanced Audio Coding",
		"Ogg Vorbis Codec",
		"Free Lossless Audio Codec",
		"Wave Audio for Windows",
		"Tracker music extended module",
		"Tracker music multi-track module",
		"Tracker music Impulse module",
		"Tracked music ProTracker module",
		"Shrunked pkzip archive",
		"Reduced pkzip archive",
		"Imploded pkzip archive",
		"PKWARE zip64 archive",
		"Zip archive",
		"Zip multi-Volume archive",
		"PKLITE compressed executable",
		"PKSFX self-extracting archive",
		"Tape Archive",
		"Roshal Archive",
		"Roshal Archive v5",
		"Gzip compress archive",
		"Bzip2 compress archive",
		"7z compress archive",
		"XZ compress archive",
		"ZStandard archive",
		"FreeArc",
		"Archive by SEA",
		"Yoshi LHA",
		"Zoo Archive",
		"Archive by Robert Jung",
		"Microsoft Cabinet",
		"Microsoft DOS KWAJ",
		"Microsoft DOS SZDD",
		"Microsoft executable",
		"Microsoft compound file",
		"CD ISO 9660",
		"CD Nero",
		"CD PowerISO",
		"CD Alcohol 120",
		"Windows Help File",
		"Portable Document Format",
		"Rich Text Format",
		"UTF-8 text",
		"UTF-16 text",
		"UTF-32 text",
		"ANSI escaped text",
		"Plain text",
	}[sign]
}

// Extension is a map of file type signatures to file extensions.
type Extension map[Signature][]string

// Ext returns a map of file type signatures to common file extensions.
func Ext() *Extension { //nolint:funlen
	exts := Extension{
		ElectronicArtsIFF:                 []string{".iff"},
		AV1ImageFile:                      []string{".avif"},
		JPEGFileInterchangeFormat:         []string{".jpg", ".jpeg"},
		JPEG2000:                          []string{".jp2", ".j2k", ".jpf", ".jpx", ".jpm", ".mj2"},
		PortableNetworkGraphics:           []string{".png"},
		GraphicsInterchangeFormat:         []string{".gif"},
		GoogleWebP:                        []string{".webp"},
		TaggedImageFileFormat:             []string{".tif", ".tiff"},
		BMPFileFormat:                     []string{".bmp"},
		PersonalComputereXchange:          []string{".pcx"},
		InterleavedBitmap:                 []string{".ilbm", ".iff"},
		MicrosoftIcon:                     []string{".ico"},
		RIPscrip:                          []string{".rip"},
		MPEG4:                             []string{".mp4"},
		QuickTimeMovie:                    []string{".mov"},
		QuickTimeM4V:                      []string{".m4v"},
		MicrosoftAudioVideoInterleave:     []string{".avi"},
		MicrosoftWindowsMedia:             []string{".wmv"},
		MPEG:                              []string{".mpg", ".mpeg"},
		FlashVideo:                        []string{".flv"},
		RealPlayer:                        []string{".rv", ".rm", ".rmvb"},
		MusicalInstrumentDigitalInterface: []string{".mid", ".midi"},
		MPEG1AudioLayer3:                  []string{".mp3"},
		MPEGAdvancedAudioCoding:           []string{".aac", ".mp3"},
		OggVorbisCodec:                    []string{".ogg"},
		FreeLosslessAudioCodec:            []string{".flac"},
		WaveAudioForWindows:               []string{".wav"},
		MusicExtendedModule:               []string{".xm"},
		MusicMultiTrackModule:             []string{".mtm"},
		MusicImpulseTracker:               []string{".it"},
		MusicProTracker:                   []string{".mod"},
		PKWAREZipShrink:                   []string{".zip"},
		PKWAREZipReduce:                   []string{".zip"},
		PKWAREZipImplode:                  []string{".zip"},
		PKWAREZip64:                       []string{".zip"},
		PKWAREZip:                         []string{".zip"},
		PKWAREMultiVolume:                 []string{".zip"},
		PKLITE:                            []string{".zip"},
		PKSFX:                             []string{".zip"},
		TapeARchive:                       []string{".tar"},
		RoshalARchive:                     []string{".rar"},
		RoshalARchivev5:                   []string{".rar"},
		GzipCompressArchive:               []string{".gz"},
		Bzip2CompressArchive:              []string{".bz2"},
		X7zCompressArchive:                []string{".7z"},
		XZCompressArchive:                 []string{".xz"},
		ZStandardArchive:                  []string{".zst"},
		FreeArc:                           []string{".arc"},
		ARChiveSEA:                        []string{".arc"},
		YoshiLHA:                          []string{".lzh", ".lha"},
		ZooArchive:                        []string{".zoo"},
		ArchiveRobertJung:                 []string{".arj"},
		MicrosoftCABinet:                  []string{".cab"},
		MicrosoftDOSKWAJ:                  []string{".com"},
		MicrosoftDOSSZDD:                  []string{".exe"},
		MicrosoftExecutable:               []string{".exe"},
		MicrosoftCompoundFile:             []string{".exe"},
		CDISO9660:                         []string{".iso"},
		CDNero:                            []string{".nri"},
		CDPowerISO:                        []string{".daa"},
		CDAlcohol120:                      []string{".mdf"},
		WindowsHelpFile:                   []string{".hlp"},
		PortableDocumentFormat:            []string{".pdf"},
		RichTextFormat:                    []string{".rtf"},
		UTF8Text:                          []string{".txt"},
		UTF16Text:                         []string{".txt"},
		UTF32Text:                         []string{".txt"},
		ANSIEscapeText:                    []string{".ans"},
		PlainText:                         []string{".txt"},
	}
	return &exts
}

// Matcher is a function that matches a byte slice to a file type.
type Matcher func(io.ReaderAt) bool

// Finder is a map of file type signatures to matchers.
type Finder map[Signature]Matcher

// New returns a new Finder with all the matchers.
//
// ANSIEscapeText and PlainText are not included as they need to be
// checked separately and in a specific order.
func New() *Finder { //nolint:funlen
	finds := Finder{
		ElectronicArtsIFF:                 Iff,
		AV1ImageFile:                      Avif,
		JPEGFileInterchangeFormat:         Jpeg,
		JPEG2000:                          Jpeg2000,
		PortableNetworkGraphics:           Png,
		GraphicsInterchangeFormat:         Gif,
		GoogleWebP:                        Webp,
		TaggedImageFileFormat:             Tiff,
		BMPFileFormat:                     Bmp,
		PersonalComputereXchange:          Pcx,
		InterleavedBitmap:                 Ilbm,
		MicrosoftIcon:                     Ico,
		RIPscrip:                          Ripscrip,
		MPEG4:                             Mp4,
		QuickTimeMovie:                    QTMov,
		QuickTimeM4V:                      M4v,
		MicrosoftAudioVideoInterleave:     Avi,
		MicrosoftWindowsMedia:             Wmv,
		MPEG:                              Mpeg,
		FlashVideo:                        Flv,
		RealPlayer:                        Ivr,
		MusicalInstrumentDigitalInterface: Midi,
		MPEG1AudioLayer3:                  Mp3,
		MPEGAdvancedAudioCoding:           AAC,
		OggVorbisCodec:                    Ogg,
		FreeLosslessAudioCodec:            Flac,
		WaveAudioForWindows:               Wave,
		MusicExtendedModule:               XM,
		MusicMultiTrackModule:             MTM,
		MusicImpulseTracker:               IT,
		MusicProTracker:                   MK,
		PKWAREZipShrink:                   PkShrink,
		PKWAREZipReduce:                   PkReduce,
		PKWAREZipImplode:                  PkImplode,
		PKWAREZip64:                       Zip64,
		PKWAREZip:                         Pkzip,
		PKWAREMultiVolume:                 PkzipMulti,
		PKLITE:                            Pklite,
		PKSFX:                             Pksfx,
		TapeARchive:                       Tar,
		RoshalARchive:                     Rar,
		RoshalARchivev5:                   Rarv5,
		GzipCompressArchive:               Gzip,
		Bzip2CompressArchive:              Bzip2,
		X7zCompressArchive:                X7z,
		XZCompressArchive:                 XZ,
		ZStandardArchive:                  ZStd,
		FreeArc:                           ArcFree,
		ARChiveSEA:                        ArcSEA,
		YoshiLHA:                          LzhLha,
		ZooArchive:                        Zoo,
		ArchiveRobertJung:                 Arj,
		MicrosoftCABinet:                  Cab,
		MicrosoftDOSKWAJ:                  DosKWAJ,
		MicrosoftDOSSZDD:                  DosSZDD,
		MicrosoftExecutable:               MSExe,
		MicrosoftCompoundFile:             MSComp,
		CDISO9660:                         ISO,
		CDNero:                            Nri,
		CDPowerISO:                        Daa,
		CDAlcohol120:                      Mdf,
		WindowsHelpFile:                   Hlp,
		PortableDocumentFormat:            Pdf,
		RichTextFormat:                    Rtf,
		UTF8Text:                          Utf8,
		UTF16Text:                         Utf16,
		UTF32Text:                         Utf32,
	}
	return &finds
}

// MatchExt determines if the reader matches the file type signature expected
// from the extension of the filename. It returns true if the file type matches and
// a found signature is always returned.
//
// A PNG encoded image using the filename TEST.PNG will return true
// and the PortableNetworkGraphics signature.
// A PNG encoded image using the filename TEST.JPG will return false
// and the PortableNetworkGraphics signature.
func MatchExt(filename string, r io.ReaderAt) (bool, Signature, error) {
	if Empty(r) {
		return false, Unknown, ErrNilReader
	}
	ext := strings.ToLower(filepath.Ext(filename))
	finds := New()
	for signature, exts := range *Ext() {
		if !slices.Contains(exts, ext) {
			continue
		}
		for find, matcher := range *finds {
			if matcher(r) && find == signature {
				return true, find, nil
			}
		}
	}
	return false, Find(r), nil
}

// Find returns the file type signature from the byte slice.
func Find(r io.ReaderAt) Signature {
	if Empty(r) {
		return ZeroByte
	}
	matchers := *New()
	for sign, matcher := range matchers {
		if matcher(r) {
			return sign
		}
	}
	switch {
	case Ansi(r):
		return ANSIEscapeText
	case CodePage(r), Txt(r):
		return PlainText
	default:
		return Unknown
	}
}

// Empty returns true if the reader is empty.
func Empty(r io.ReaderAt) bool {
	if r == nil {
		return true
	}
	p := make([]byte, 1)
	sr := io.NewSectionReader(r, 0, 1)
	if n, err := sr.Read(p); err != nil || n < 1 {
		return true
	}
	return false
}
