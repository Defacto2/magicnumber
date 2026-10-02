package magicnumber_test

import (
	"errors"
	"fmt"
	"log"
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

const Format = "magic number find false positive, got %s (%d) not %s (%d)"

// _exampleReadme matches the README.md example and is used to lint and validate the syntax.
func _exampleReadme() { //nolint:unused
	w := os.Stdout
	file, err := os.Open("example.exe")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Option 1.
	result := magicnumber.Find(file)
	fmt.Fprintf(w, "File type: %s\n", result)

	// Option 2.
	valid, result, err := magicnumber.MatchExt("example.exe", file)
	if err != nil {
		fmt.Fprintln(os.Stderr)
	}
	fmt.Fprintf(w, "File type: %s\n", result)
	fmt.Fprintf(w, "File extension valid: %v\n", valid)
}

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

var ErrCaller = errors.New("runtime caller failed")

const testdata = "testdata"

func pathUncompress(tb testing.TB, name string) string {
	tb.Helper()
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		tb.Fatal(ErrCaller)
	}
	d := filepath.Dir(file)
	const uncompress = "uncompress"
	x := filepath.Join(d, testdata, uncompress, name)
	return x
}

func pathMp3(tb testing.TB, name string) string {
	tb.Helper()
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		tb.Fatal(ErrCaller)
	}
	d := filepath.Dir(file)
	const mp3 = "mp3"
	return filepath.Join(d, testdata, mp3, name)
}

func pathDisc(tb testing.TB, name string) string {
	tb.Helper()
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		tb.Fatal(ErrCaller)
	}
	d := filepath.Dir(file)
	const discimages = "discimages"
	return filepath.Join(d, testdata, discimages, name)
}

func pathFile(tb testing.TB, name string) string {
	tb.Helper()
	_, file, _, usable := runtime.Caller(0)
	if !usable {
		tb.Fatal(ErrCaller)
	}
	d := filepath.Dir(file)
	return filepath.Join(d, testdata, name)
}

func TestUnknowns(t *testing.T) {
	t.Parallel()

	s := "some binary data"
	r := strings.NewReader(s)
	got, err := magicnumber.Archive(r)
	be.Err(t, err, nil)
	be.Equal(t, got, magicnumber.Unknown)
	be.Equal(t, got.String(), "binary data or text")
	be.Equal(t, got.Title(), "Binary data or binary text")

	b, got, err := magicnumber.MatchExt(emptyFile, r)
	be.Err(t, err, nil)
	be.True(t, !b)
	be.Equal(t, got, magicnumber.PlainText)

	f, err := os.Open(pathUncompress(t, emptyFile))
	be.Err(t, err, nil)
	defer f.Close()
	got = magicnumber.Find(f)
	be.Equal(t, got, magicnumber.ZeroByte)
}

func TestLastSignature(t *testing.T) {
	t.Parallel()

	const got = magicnumber.LastSignature

	// test for panic conditions when modifying the list of titles and names
	be.Equal(t, got.String(), "XBIN binary text")
	be.Equal(t, got.Title(), "XBIN extended binary text")
}

func TestFind(t *testing.T) {
	t.Parallel()
	// walk the assets directory
	err := filepath.Walk(pathFile(t, ""), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToUpper(filepath.Ext(path))
		if info.IsDir() || ext == "" {
			return nil
		}
		base := filepath.Base(path)
		skip := [2]string{"SAMPLE.DAT", "uncompress.bin"}
		if slices.Contains(skip[:], base) {
			return nil
		}
		f, err := os.Open(path) //nolint:gosec
		be.Err(t, err, nil)
		defer f.Close()
		got := magicnumber.Find(f)
		if base == "τεχτƒιℓε.τχτ" {
			be.Equal(t, got, magicnumber.PlainText)
			return nil
		}

		t.Log(ext, got, path, info.Name(), info.Size())
		switch ext {
		case ".COM":
			// do not test as it returns different results based on the file
			return nil
		case ".7Z":
			be.Equal(t, got, magicnumber.X7zCompressArchive)
		case ".ANS":
			be.Equal(t, got, magicnumber.ANSIEscapeText)
		case ".ARC":
			// two different signatures used for the same file extension
			s := [2]magicnumber.Signature{
				magicnumber.FreeArc, magicnumber.ARChiveSEA,
			}
			be.True(t, slices.Contains(s[:], got))
		case ".ARJ":
			be.Equal(t, got, magicnumber.ArchiveRobertJung)
		case ".AVIF":
			be.Equal(t, got, magicnumber.AV1ImageFile)
		case ".BAT", ".INI", ".CUE":
			be.Equal(t, got, magicnumber.PlainText)
		case ".BMP":
			be.Equal(t, got, magicnumber.BMPFileFormat)
		case ".BZ2":
			be.Equal(t, got, magicnumber.Bzip2CompressArchive)
		case ".CHM", ".HLP":
			be.Equal(t, got, magicnumber.WindowsHelpFile)
		case ".DAA":
			be.Equal(t, got, magicnumber.CDPowerISO)
		case ".EXE", ".DLL":
			be.Equal(t, got, magicnumber.MicrosoftExecutable)
		case ".GIF":
			be.Equal(t, got, magicnumber.GraphicsInterchangeFormat)
		case ".GZ":
			be.Equal(t, got, magicnumber.GzipCompressArchive)
		case ".JPG", ".JPEG":
			be.Equal(t, got, magicnumber.JPEGFileInterchangeFormat)
		case ".ICO":
			be.Equal(t, got, magicnumber.MicrosoftIcon)
		case ".IFF":
			be.Equal(t, got, magicnumber.InterleavedBitmap)
		case ".ISO":
			be.Equal(t, got, magicnumber.CDISO9660)
		case ".LZH":
			be.Equal(t, got, magicnumber.YoshiLHA)
		case ".MP3":
			// do not test as it returns different results based on the file's ID3 tag
			return nil
		case ".PAK":
			be.Equal(t, got, magicnumber.NoGatePAK)
		case ".PCX":
			be.Equal(t, got, magicnumber.PersonalComputereXchange)
		case ".PNG":
			be.Equal(t, got, magicnumber.PortableNetworkGraphics)
		case ".RAR":
			signs := [2]magicnumber.Signature{
				magicnumber.RoshalARchivev5,
				magicnumber.RoshalARchive,
			}
			be.True(t, slices.Contains(signs[:], got))
		case ".TAR":
			be.Equal(t, got, magicnumber.TapeARchive)
		case ".TXT", ".MD", ".NFO", ".ME", ".DIZ", ".ASC", ".CAP", ".DOC":
			signs := [2]magicnumber.Signature{
				magicnumber.PlainText,
				magicnumber.UTF16Text,
			}
			be.True(t, slices.Contains(signs[:], got))
		case ".WEBP":
			be.Equal(t, got, magicnumber.GoogleWebP)
		case ".XZ":
			be.Equal(t, got, magicnumber.XZCompressArchive)
		case ".ZIP":
			if base == "EMPTY.ZIP" {
				be.Equal(t, got, magicnumber.ZeroByte)
				return nil
			}
			zips := [5]magicnumber.Signature{
				magicnumber.PKWAREZip,
				magicnumber.PKWAREZip64,
				magicnumber.PKWAREZipImplode,
				magicnumber.PKWAREZipReduce,
				magicnumber.PKWAREZipShrink,
			}
			be.True(t, slices.Contains(zips[:], got))
		default:
			be.True(t, magicnumber.Unknown != got)
			fmt.Fprintln(os.Stderr, ext, filepath.Base(path), fmt.Sprint(got))
		}

		return nil
	})
	be.Err(t, err, nil)
}
