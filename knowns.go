package magicnumber

import (
	"io"
	"log/slog"
	"sync"
)

// Archive reads all the bytes from the reader and returns the file type signature if
// the file is a known archive of files or Unknown if the file is not an archive.
func Archive(r io.ReaderAt) (Signature, error) {
	if r == nil {
		return Unknown, ErrNilReader
	}

	for _, item := range archives() {
		if item.matcher(r) {
			return item.sig, nil
		}
	}

	return Unknown, nil
}

// Archives returns all the archive file type signatures.
func Archives() []Signature {
	return packs
}

// DiscImage reads all the bytes from the reader and returns the file type signature if
// the file is a known CD disk image or Unknown if the file is not a disk image.
func DiscImage(r io.ReaderAt) (Signature, error) {
	if r == nil {
		return Unknown, ErrNilReader
	}

	for _, item := range discImages() {
		if item.matcher(r) {
			return item.sig, nil
		}
	}

	return Unknown, nil
}

// DiscImages returns all the CD disk image file type signatures.
func DiscImages() []Signature {
	return discs
}

// ArchivesBBS returns all the archive file type signatures that were
// commonly used in the BBS online era of the 1980s and early 1990s.
// Eventually these were replaced by the universal ZIP format using
// the Deflate and Store compression methods.
func ArchivesBBS() []Signature {
	return bbsPacks
}

// Document reads all the bytes from the reader and returns the file type signature if
// the file is a known document or Unknown if the file is not a document.
func Document(r io.ReaderAt) (Signature, error) {
	if r == nil {
		return Unknown, ErrNilReader
	}

	for _, item := range documents() {
		if item.matcher(r) {
			return item.sig, nil
		}
	}
	switch {
	case Ansi(r):
		return ANSIEscapeText, nil
	case CodePage(r):
		return PlainText, nil
	case Txt(r):
		return PlainText, nil
	default:
		return Unknown, nil
	}
}

// Documents returns all the documentation file type signatures.
func Documents() []Signature {
	return docs
}

// Image reads all the bytes from the reader and returns the file type signature if
// the file is a known image or Unknown if the file is not an image.
func Image(r io.ReaderAt) (Signature, error) {
	if r == nil {
		return Unknown, ErrNilReader
	}

	for _, item := range images() {
		if item.matcher(r) {
			return item.sig, nil
		}
	}

	return Unknown, nil
}

// Images returns all the image file type signatures.
func Images() []Signature {
	return imgs
}

// Program reads all the bytes from the reader and returns the file type signature if
// the file is a known DOS or Windows program or Unknown if the file is not a program.
func Program(r io.ReaderAt) (Signature, error) {
	if r == nil {
		return Unknown, ErrNilReader
	}

	for _, item := range programs() {
		if item.matcher(r) {
			return item.sig, nil
		}
	}

	return Unknown, nil
}

// Programs returns all the program file type signatures for
// Microsoft operating systems, DOS and Windows.
func Programs() []Signature {
	return progs
}

// Text reads the first 512 bytes from the reader and returns the file type signature if
// the file is a known plain text file or Unknown if the file is not a text file.
func Text(r io.ReaderAt) (Signature, error) {
	return TextWithLogger(nil, r)
}

// Deprecated: use [TextWithLogger] instead.
// The io.Writer is unused.
func TextW(_ io.Writer, r io.ReaderAt) (Signature, error) {
	return TextWithLogger(nil, r)
}

// TextWithLogger reads the first 512 bytes from the reader and returns the file type signature if
// the file is a known plain text file or Unknown if the file is not a text file.
func TextWithLogger(sl *slog.Logger, r io.ReaderAt) (Signature, error) {
	if sl == nil {
		sl = slog.New(slog.DiscardHandler)
	}
	if r == nil {
		return Unknown, ErrNilReader
	}

	const msg = "known texts"
	for _, item := range texts() {
		if item.matcher(r) {
			sl.Debug(msg+" finder matched", slog.Int("signature", int(item.sig)), slog.String("title", item.sig.Title()))
			return item.sig, nil
		}
	}
	switch {
	case AnsiWithLogger(sl, r):
		return ANSIEscapeText, nil
	case CodePageWithLogger(sl, r):
		return PlainText, nil
	case TxtWithLogger(sl, r):
		return PlainText, nil
	default:
		sl.Debug(msg + " returned a default, unknown")
		return Unknown, nil
	}
}

// Texts returns all the text file type signatures.
func Texts() []Signature {
	return txts
}

// Video reads all the bytes from the reader and returns the file type signature if
// the file is a known video or Unknown if the file is not a video.
func Video(r io.ReaderAt) (Signature, error) {
	if r == nil {
		return Unknown, ErrNilReader
	}

	for _, item := range videos() {
		if item.matcher(r) {
			return item.sig, nil
		}
	}

	return Unknown, nil
}

// Videos returns all the video file type signatures.
func Videos() []Signature {
	return vids
}

type sigMatcher struct {
	sig     Signature
	matcher Matcher
}

var archives = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, pack := range packs {
		if m, exists := find[pack]; exists {
			matchers = append(matchers, sigMatcher{sig: pack, matcher: m})
		}
	}
	return matchers
})

var packs = []Signature{ //nolint:gochecknoglobals
	PKWAREZipShrink,
	PKWAREZipReduce,
	PKWAREZipImplode,
	PKWAREZip64,
	PKWAREZip,
	PKWAREMultiVolume,
	PKLITE,
	PKSFX,
	TapeARchive,
	RoshalARchive,
	RoshalARchivev5,
	GzipCompressArchive,
	Bzip2CompressArchive,
	X7zCompressArchive,
	XZCompressArchive,
	ZStandardArchive,
	FreeArc,
	NoGatePAK, // NoGatePAK must go before ARChiveSEA
	ARChiveSEA,
	YoshiLHA,
	ZooArchive,
	ArchiveRobertJung,
	MicrosoftCABinet,
}

var bbsPacks = []Signature{ //nolint:gochecknoglobals
	PKWAREZipShrink,
	PKWAREZipReduce,
	PKWAREZipImplode,
	ARChiveSEA,
	YoshiLHA,
	ZooArchive,
	ArchiveRobertJung,
	NoGatePAK,
}

var discImages = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, disc := range discs {
		if m, exists := find[disc]; exists {
			matchers = append(matchers, sigMatcher{sig: disc, matcher: m})
		}
	}
	return matchers
})

var discs = []Signature{ //nolint:gochecknoglobals
	CDISO9660,
	CDNero,
	CDPowerISO,
	CDAlcohol120,
}

var documents = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, doc := range docs {
		if m, exists := find[doc]; exists {
			matchers = append(matchers, sigMatcher{sig: doc, matcher: m})
		}
	}
	return matchers
})

var docs = []Signature{ //nolint:gochecknoglobals
	WindowsHelpFile,
	PortableDocumentFormat,
	RichTextFormat,
	UTF8Text,
	UTF16Text,
	UTF32Text,
}

var images = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, img := range imgs {
		if m, exists := find[img]; exists {
			matchers = append(matchers, sigMatcher{sig: img, matcher: m})
		}
	}
	return matchers
})

var imgs = []Signature{ //nolint:gochecknoglobals
	AV1ImageFile,
	JPEGFileInterchangeFormat,
	JPEG2000,
	PortableNetworkGraphics,
	GraphicsInterchangeFormat,
	GoogleWebP,
	TaggedImageFileFormat,
	BMPFileFormat,
	PersonalComputereXchange,
	InterleavedBitmap,
	MicrosoftIcon,
	RIPscrip,
	ElectronicArtsAnim,
	ElectronicArtsIFF,
}

var programs = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, prog := range progs {
		if m, exists := find[prog]; exists {
			matchers = append(matchers, sigMatcher{sig: prog, matcher: m})
		}
	}
	return matchers
})

var progs = []Signature{ //nolint:gochecknoglobals
	MicrosoftExecutable,
	MicrosoftDOSKWAJ,
	MicrosoftDOSSZDD,
	MicrosoftCompoundFile,
}

var texts = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, txt := range txts {
		if m, exists := find[txt]; exists {
			matchers = append(matchers, sigMatcher{sig: txt, matcher: m})
		}
	}
	return matchers
})

var txts = []Signature{ //nolint:gochecknoglobals
	UTF8Text,
	UTF16Text,
	UTF32Text,
	ANSIEscapeText,
	PlainText,
}

var videos = sync.OnceValue(func() []sigMatcher { //nolint:gochecknoglobals
	find := defaultFinder()
	var matchers []sigMatcher
	for _, vid := range vids {
		if m, exists := find[vid]; exists {
			matchers = append(matchers, sigMatcher{sig: vid, matcher: m})
		}
	}
	return matchers
})

var vids = []Signature{ //nolint:gochecknoglobals
	MPEG4,
	QuickTimeMovie,
	QuickTimeM4V,
	MicrosoftAudioVideoInterleave,
	MicrosoftWindowsMedia,
	MPEG,
	FlashVideo,
	RealPlayer,
}
