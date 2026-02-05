package magicnumber_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
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
	// binary data or text
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
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.Unknown, sign)
	be.Equal(t, sign.String(), "binary data or text")
	be.Equal(t, sign.Title(), "Binary data or binary text")

	b, sign, err := magicnumber.MatchExt(emptyFile, nr)
	be.Err(t, err, nil)
	be.True(t, !b)
	be.Equal(t, magicnumber.PlainText, sign)

	r, err := os.Open(uncompress(emptyFile))
	be.Err(t, err, nil)
	defer r.Close()
	sign = magicnumber.Find(r)
	be.Equal(t, magicnumber.ZeroByte, sign)
}

func TestFind(t *testing.T) {
	t.Parallel()
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
		if slices.Contains(skip, base) {
			return nil
		}
		f, err := os.Open(path)
		be.Err(t, err, nil)
		defer f.Close()
		sign := magicnumber.Find(f)
		if base == "τεχτƒιℓε.τχτ" {
			be.Equal(t, magicnumber.PlainText, sign)
			return nil
		}

		switch ext {
		case ".COM":
			// do not test as it returns different results based on the file
			return nil
		case ".7Z":
			be.Equal(t, magicnumber.X7zCompressArchive, sign)
		case ".ANS":
			be.Equal(t, magicnumber.ANSIEscapeText, sign)
		case ".ARC":
			// two different signatures used for the same file extension
			be.True(t, slices.Contains([]magicnumber.Signature{
				magicnumber.FreeArc, magicnumber.ARChiveSEA,
			}, sign))
		case ".ARJ":
			be.Equal(t, magicnumber.ArchiveRobertJung, sign)
		case ".AVIF":
			be.Equal(t, magicnumber.AV1ImageFile, sign)
		case ".BAT", ".INI", ".CUE":
			be.Equal(t, magicnumber.PlainText, sign)
		case ".BMP":
			be.Equal(t, magicnumber.BMPFileFormat, sign)
		case ".BZ2":
			be.Equal(t, magicnumber.Bzip2CompressArchive, sign)
		case ".CHM", ".HLP":
			be.Equal(t, magicnumber.WindowsHelpFile, sign)
		case ".DAA":
			be.Equal(t, magicnumber.CDPowerISO, sign)
		case ".EXE", ".DLL":
			be.Equal(t, magicnumber.MicrosoftExecutable, sign)
		case ".GIF":
			be.Equal(t, magicnumber.GraphicsInterchangeFormat, sign)
		case ".GZ":
			be.Equal(t, magicnumber.GzipCompressArchive, sign)
		case ".JPG", ".JPEG":
			be.Equal(t, magicnumber.JPEGFileInterchangeFormat, sign)
		case ".ICO":
			be.Equal(t, magicnumber.MicrosoftIcon, sign)
		case ".IFF":
			be.Equal(t, magicnumber.InterleavedBitmap, sign)
		case ".ISO":
			be.Equal(t, magicnumber.CDISO9660, sign)
		case ".LZH":
			be.Equal(t, magicnumber.YoshiLHA, sign)
		case ".MP3":
			// do not test as it returns different results based on the file's ID3 tag
			return nil
		case ".PCX":
			be.Equal(t, magicnumber.PersonalComputereXchange, sign)
		case ".PNG":
			be.Equal(t, magicnumber.PortableNetworkGraphics, sign)
		case ".RAR":
			signs := []magicnumber.Signature{
				magicnumber.RoshalARchivev5,
				magicnumber.RoshalARchive,
			}
			be.True(t, slices.Contains(signs, sign))
		case ".TAR":
			be.Equal(t, magicnumber.TapeARchive, sign)
		case ".TXT", ".MD", ".NFO", ".ME", ".DIZ", ".ASC", ".CAP", ".DOC":
			signs := []magicnumber.Signature{
				magicnumber.PlainText,
				magicnumber.UTF16Text,
			}
			be.True(t, slices.Contains(signs, sign))
		case ".WEBP":
			be.Equal(t, magicnumber.GoogleWebP, sign)
		case ".XZ":
			be.Equal(t, magicnumber.XZCompressArchive, sign)
		case ".ZIP":
			if base == "EMPTY.ZIP" {
				be.Equal(t, magicnumber.ZeroByte, sign)
				return nil
			}
			zips := []magicnumber.Signature{
				magicnumber.PKWAREZip,
				magicnumber.PKWAREZip64,
				magicnumber.PKWAREZipImplode,
				magicnumber.PKWAREZipReduce,
				magicnumber.PKWAREZipShrink,
			}
			be.True(t, slices.Contains(zips, sign))
		default:
			be.True(t, magicnumber.Unknown != sign)
			fmt.Fprintln(os.Stderr, ext, filepath.Base(path), fmt.Sprint(sign))
		}

		return nil
	})
	be.Err(t, err, nil)
}
