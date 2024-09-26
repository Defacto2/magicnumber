package magicnumber_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyFile  = "EMPTY"
	avifFile   = "TEST.avif"
	bmpFile    = "TEST.BMP"
	gifFile    = "TEST.GIF"
	gif2File   = "TEST2.gif"
	ilbmFile   = "TEST.IFF"
	jpegFile   = "TEST.JPEG"
	jpgFile    = "TEST.JPG"
	icoFile    = "favicon.ico"
	pcxFile    = "TEST.PCX"
	pngFile    = "TEST.PNG"
	rtfFile    = "TEST.rtf"
	webpFile   = "TEST.webp"
	asciiFile  = "TEST.ASC"
	ansiFile   = "TEST.ANS"
	txtFile    = "TEST.TXT"
	badFile    = "τεχτƒιℓε.τχτ"
	manualFile = "PKZ204EX.TXT"
	pdfFile    = "TEST.pdf"
	utf16File  = "TEST-U16.txt"
	iso7File   = "TEST-8859-7.txt"
	modFile    = "TEST.mod"
	xmFile     = "TEST.xm"
	itFile     = "TEST.it"
	amigaIFF   = "TEST0.IFF"
	wavFile    = "TEST.wav"
	mp3File    = "TEST.mp3"
	oggFile    = "TEST.ogg"
	wmaFile    = "TEST.wma"
)

func ExampleArchive() {
	f1, err := os.Open(filepath.Join("testdata", "TEST.cab"))
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	f2, err := os.Open(filepath.Join("testdata", "README.md"))
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	sign1, err := magicnumber.Archive(f1)
	if err != nil {
		panic(err)
	}
	fmt.Println(sign1)

	sign2, err := magicnumber.Archive(f2)
	if err != nil {
		panic(err)
	}
	fmt.Println(sign2)
	// Output: Microsoft cabinet
	// binary data
}

func ExampleFind() {
	f, err := os.Open(filepath.Join("testdata", "TEST.cab"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sign := magicnumber.Find(f)
	fmt.Println(sign.String())
	fmt.Println(sign.Title())
	// Output: Microsoft cabinet
	// Microsoft Cabinet
}

func ExampleFindExecutable() {
	f, err := os.Open(filepath.Join("testdata", "binaries", "windows9x", "7za920", "7za.exe"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	win, err := magicnumber.FindExecutable(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(win.String())
	// Output: Windows NT v4.0
}

func uncompress(name string) string {
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		panic("runtime.Caller failed")
	}
	d := filepath.Dir(file)
	x := filepath.Join(d, "testdata", "uncompress", name)
	return x
}

func mp3file(name string) string {
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		panic("runtime.Caller failed")
	}
	d := filepath.Dir(file)
	x := filepath.Join(d, "testdata", "mp3", name)
	return x
}

func imgfile(name string) string {
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		panic("runtime.Caller failed")
	}
	d := filepath.Dir(file)
	x := filepath.Join(d, "testdata", "discimages", name)
	return x
}

func tdfile(name string) string {
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		panic("runtime.Caller failed")
	}
	d := filepath.Dir(file)
	x := filepath.Join(d, "testdata", name)
	return x
}

func TestUnknowns(t *testing.T) {
	t.Parallel()

	data := "some binary data"
	nr := strings.NewReader(data)
	sign, err := magicnumber.Archive(nr)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.Unknown, sign)
	assert.Equal(t, "binary data", sign.String())
	assert.Equal(t, "Binary data", sign.Title())

	b, sign, err := magicnumber.MatchExt(emptyFile, nr)
	require.NoError(t, err)
	assert.False(t, b)
	assert.Equal(t, magicnumber.PlainText, sign)

	r, err := os.Open(uncompress(emptyFile))
	require.NoError(t, err)
	defer r.Close()
	sign = magicnumber.Find(r)
	assert.Equal(t, magicnumber.ZeroByte, sign)
}

func TestFind(t *testing.T) {
	t.Parallel()
	prob := func(ext, path string) string {
		return fmt.Sprintf("ext: %s, path: %s", ext, path)
	}
	// walk the assets directory
	err := filepath.Walk(tdfile(""), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToUpper(filepath.Ext(path))
		if info.IsDir() || ext == "" {
			return nil
		}
		base := filepath.Base(path)
		skip := []string{"SAMPLE.DAT", "uncompress.bin"}
		for _, s := range skip {
			if s == base {
				return nil
			}
		}
		f, err := os.Open(path)
		require.NoError(t, err)
		defer f.Close()
		sign := magicnumber.Find(f)
		if base == "τεχτƒιℓε.τχτ" {
			assert.Equal(t, magicnumber.PlainText, sign, prob(ext, path))
			return nil
		}

		switch ext {
		case ".COM":
			// do not test as it returns different results based on the file
			return nil
		case ".7Z":
			assert.Equal(t, magicnumber.X7zCompressArchive, sign, prob(ext, path))
		case ".ANS":
			assert.Equal(t, magicnumber.ANSIEscapeText, sign, prob(ext, path))
		case ".ARC":
			// two different signatures used for the same file extension
			assert.Contains(t, []magicnumber.Signature{
				magicnumber.FreeArc,
				magicnumber.ARChiveSEA,
			}, sign, prob(ext, path))
		case ".ARJ":
			assert.Equal(t, magicnumber.ArchiveRobertJung, sign, prob(ext, path))
		case ".AVIF":
			assert.Equal(t, magicnumber.AV1ImageFile, sign, prob(ext, path))
		case ".BAT", ".INI", ".CUE":
			assert.Equal(t, magicnumber.PlainText, sign, prob(ext, path))
		case ".BMP":
			assert.Equal(t, magicnumber.BMPFileFormat, sign, prob(ext, path))
		case ".BZ2":
			assert.Equal(t, magicnumber.Bzip2CompressArchive, sign, prob(ext, path))
		case ".CHM", ".HLP":
			assert.Equal(t, magicnumber.WindowsHelpFile, sign, prob(ext, path))
		case ".DAA":
			assert.Equal(t, magicnumber.CDPowerISO, sign, prob(ext, path))
		case ".EXE", ".DLL":
			assert.Equal(t, magicnumber.MicrosoftExecutable, sign, prob(ext, path))
		case ".GIF":
			assert.Equal(t, magicnumber.GraphicsInterchangeFormat, sign, prob(ext, path))
		case ".GZ":
			assert.Equal(t, magicnumber.GzipCompressArchive, sign, prob(ext, path))
		case ".JPG", ".JPEG":
			assert.Equal(t, magicnumber.JPEGFileInterchangeFormat, sign, prob(ext, path))
		case ".ICO":
			assert.Equal(t, magicnumber.MicrosoftIcon, sign, prob(ext, path))
		case ".IFF":
			assert.Equal(t, magicnumber.InterleavedBitmap, sign, prob(ext, path))
		case ".ISO":
			assert.Equal(t, magicnumber.CDISO9660, sign, prob(ext, path))
		case ".LZH":
			assert.Equal(t, magicnumber.YoshiLHA, sign, prob(ext, path))
		case ".MP3":
			// do not test as it returns different results based on the file's ID3 tag
			return nil
		case ".PCX":
			assert.Equal(t, magicnumber.PersonalComputereXchange, sign, prob(ext, path))
		case ".PNG":
			assert.Equal(t, magicnumber.PortableNetworkGraphics, sign, prob(ext, path))
		case ".RAR":
			assert.Contains(t, []magicnumber.Signature{
				magicnumber.RoshalARchivev5,
				magicnumber.RoshalARchive,
			}, sign, prob(ext, path))
		case ".TAR":
			assert.Equal(t, magicnumber.TapeARchive, sign, prob(ext, path))
		case ".TXT", ".MD", ".NFO", ".ME", ".DIZ", ".ASC", ".CAP", ".DOC":
			assert.Contains(t, []magicnumber.Signature{
				magicnumber.PlainText,
				magicnumber.UTF16Text,
			}, sign, prob(ext, path))
		case ".WEBP":
			assert.Equal(t, magicnumber.GoogleWebP, sign, prob(ext, path))
		case ".XZ":
			assert.Equal(t, magicnumber.XZCompressArchive, sign, prob(ext, path))
		case ".ZIP":
			if base == "EMPTY.ZIP" {
				assert.Equal(t, magicnumber.ZeroByte, sign, prob(ext, path))
				return nil
			}
			zips := []magicnumber.Signature{
				magicnumber.PKWAREZip,
				magicnumber.PKWAREZip64,
				magicnumber.PKWAREZipImplode,
				magicnumber.PKWAREZipReduce,
				magicnumber.PKWAREZipShrink,
			}
			assert.Contains(t, zips, sign, prob(ext, path))
		default:
			assert.NotEqual(t, magicnumber.Unknown, sign, prob(ext, path))
			fmt.Fprintln(os.Stderr, ext, filepath.Base(path), fmt.Sprint(sign))
		}

		return nil
	})
	require.NoError(t, err)
}
