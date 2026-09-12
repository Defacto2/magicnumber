package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func TestMediaIcon(t *testing.T) {
	t.Parallel()
	r, err := os.Open(pathUncompress(t, icoFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ico(r))

	got := magicnumber.Find(r)
	const want = magicnumber.MicrosoftIcon
	be.Equal(t, got, want)
	if got != want {
		t.Fatalf(Format, got, got, want, want)
	}
}

func TestMediaAVIF(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, avifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Avif(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.AV1ImageFile)
}

func TestMediaBMP(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, bmpFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Bmp(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.BMPFileFormat)
}

func TestMediaGif(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gif(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GraphicsInterchangeFormat)

	r, err = os.Open(pathUncompress(t, gif2File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gif(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GraphicsInterchangeFormat)
}

func TestMediaIlbm(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, ilbmFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ilbm(r))
	got := magicnumber.Find(r)
	const want = magicnumber.InterleavedBitmap
	be.Equal(t, got, want)
	if got != want {
		t.Fatalf(Format, got, got, want, want)
	}

	r, err = os.Open(pathUncompress(t, amigaIFF))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ilbm(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.InterleavedBitmap)
	x, y := magicnumber.IlbmDecode(r)
	be.Equal(t, x, 200)
	be.Equal(t, y, 144)
}

func TestMediaJpeg(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, jpgFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Jpeg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.JPEGFileInterchangeFormat)

	r, err = os.Open(pathUncompress(t, jpegFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Jpeg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.JPEGFileInterchangeFormat)
}

func TestMediaPCX(t *testing.T) {
	t.Parallel()
	r, err := os.Open(pathUncompress(t, pcxFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pcx(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PersonalComputereXchange)
}

func TestMediaPNG(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, pngFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Png(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PortableNetworkGraphics)
	sign, err := magicnumber.Image(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.PortableNetworkGraphics)
}

func TestMediaWebp(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathUncompress(t, webpFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Webp(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GoogleWebP)
}

func TestMediaWave(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathMp3(t, wavFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Wave(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.WaveAudioForWindows)
}

func TestMediaMP3(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathMp3(t, mp3File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Mp3(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MPEG1AudioLayer3)
}

func TestMediaOGG(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathMp3(t, oggFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ogg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.OggVorbisCodec)
}

func TestMediaWMA(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathMp3(t, wmaFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Wmv(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MicrosoftWindowsMedia)
	sign, err := magicnumber.Video(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.MicrosoftWindowsMedia)
}

func TestMediaFlac(t *testing.T) {
	t.Parallel()

	r, err := os.Open(pathMp3(t, "TEST.flac"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Flac(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.FreeLosslessAudioCodec)
}
