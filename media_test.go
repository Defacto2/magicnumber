package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func TestMediaIcon(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(icoFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ico(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MicrosoftIcon)
}

func TestMediaAVIF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(avifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Avif(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.AV1ImageFile)
}

func TestMediaBMP(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(bmpFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Bmp(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.BMPFileFormat)
}

func TestMediaGif(t *testing.T) {
	t.Parallel()

	r, err := os.Open(uncompress(gifFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gif(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GraphicsInterchangeFormat)

	r, err = os.Open(uncompress(gif2File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gif(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GraphicsInterchangeFormat)
}

func TestMediaIlbm(t *testing.T) {
	t.Parallel()

	r, err := os.Open(uncompress(ilbmFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ilbm(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.InterleavedBitmap)

	r, err = os.Open(uncompress(amigaIFF))
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

	r, err := os.Open(uncompress(jpgFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Jpeg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.JPEGFileInterchangeFormat)

	r, err = os.Open(uncompress(jpegFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Jpeg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.JPEGFileInterchangeFormat)
}

func TestMediaPCX(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pcxFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pcx(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.PersonalComputereXchange)
}

func TestMediaPNG(t *testing.T) {
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

func TestMediaWebp(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(webpFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Webp(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.GoogleWebP)
}

func TestMediaWave(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(wavFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Wave(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.WaveAudioForWindows)
}

func TestMediaMP3(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(mp3File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Mp3(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MPEG1AudioLayer3)
}

func TestMediaOGG(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(oggFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Ogg(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.OggVorbisCodec)
}

func TestMediaWMA(t *testing.T) {
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

func TestMediaFlac(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file("TEST.flac"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Flac(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.FreeLosslessAudioCodec)
}
