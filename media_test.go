package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func TestIcon(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(icoFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ico(r))
	be.Equal(t, magicnumber.MicrosoftIcon, magicnumber.Find(r))
}

func TestAVIF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(avifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Avif(r))
	be.Equal(t, magicnumber.AV1ImageFile, magicnumber.Find(r))
}

func TestBMP(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(bmpFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Bmp(r))
	be.Equal(t, magicnumber.BMPFileFormat, magicnumber.Find(r))
}

func TestGif(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gif(r))
	be.Equal(t, magicnumber.GraphicsInterchangeFormat, magicnumber.Find(r))
	r, err = os.Open(uncompress(gif2File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gif(r))
	be.Equal(t, magicnumber.GraphicsInterchangeFormat, magicnumber.Find(r))
}

func TestIlbm(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ilbmFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ilbm(r))
	be.Equal(t, magicnumber.InterleavedBitmap, magicnumber.Find(r))

	r, err = os.Open(uncompress(amigaIFF))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ilbm(r))
	be.Equal(t, magicnumber.InterleavedBitmap, magicnumber.Find(r))
	x, y := magicnumber.IlbmDecode(r)
	be.Equal(t, 200, x)
	be.Equal(t, 144, y)
}

func TestJpeg(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(jpgFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Jpeg(r))
	be.Equal(t, magicnumber.JPEGFileInterchangeFormat, magicnumber.Find(r))
	r, err = os.Open(uncompress(jpegFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Jpeg(r))
	be.Equal(t, magicnumber.JPEGFileInterchangeFormat, magicnumber.Find(r))
}

func TestPCX(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pcxFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pcx(r))
	be.Equal(t, magicnumber.PersonalComputereXchange, magicnumber.Find(r))
}

func TestPNG(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pngFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Png(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PortableNetworkGraphics)
	sign, err := magicnumber.Image(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.PortableNetworkGraphics)
}

func TestWebp(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(webpFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Webp(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GoogleWebP)
}

func TestWave(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(wavFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Wave(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.WaveAudioForWindows)
}

func TestMP3(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(mp3File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Mp3(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MPEG1AudioLayer3)
}

func TestOGG(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(oggFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ogg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.OggVorbisCodec)
}

func TestWMA(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(wmaFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Wmv(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MicrosoftWindowsMedia)
	sign, err := magicnumber.Video(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.MicrosoftWindowsMedia)
}

func TestFlac(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file("TEST.flac"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Flac(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.FreeLosslessAudioCodec)
}
