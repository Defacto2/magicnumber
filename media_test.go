package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIcon(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(icoFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Ico(r))
	assert.Equal(t, magicnumber.MicrosoftIcon, magicnumber.Find(r))
}

func TestAVIF(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(avifFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Avif(r))
	assert.Equal(t, magicnumber.AV1ImageFile, magicnumber.Find(r))
}

func TestBMP(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(bmpFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Bmp(r))
	assert.Equal(t, magicnumber.BMPFileFormat, magicnumber.Find(r))
}

func TestGif(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(gifFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Gif(r))
	assert.Equal(t, magicnumber.GraphicsInterchangeFormat, magicnumber.Find(r))
	r, err = os.Open(uncompress(gif2File))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Gif(r))
	assert.Equal(t, magicnumber.GraphicsInterchangeFormat, magicnumber.Find(r))
}

func TestIlbm(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(ilbmFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Ilbm(r))
	assert.Equal(t, magicnumber.InterleavedBitmap, magicnumber.Find(r))

	r, err = os.Open(uncompress(amigaIFF))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Ilbm(r))
	assert.Equal(t, magicnumber.InterleavedBitmap, magicnumber.Find(r))
	x, y := magicnumber.IlbmDecode(r)
	assert.Equal(t, 200, x)
	assert.Equal(t, 144, y)
}

func TestJpeg(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(jpgFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Jpeg(r))
	assert.Equal(t, magicnumber.JPEGFileInterchangeFormat, magicnumber.Find(r))
	r, err = os.Open(uncompress(jpegFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Jpeg(r))
	assert.Equal(t, magicnumber.JPEGFileInterchangeFormat, magicnumber.Find(r))
}

func TestPCX(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pcxFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Pcx(r))
	assert.Equal(t, magicnumber.PersonalComputereXchange, magicnumber.Find(r))
}

func TestPNG(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(pngFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Png(r))
	assert.Equal(t, magicnumber.PortableNetworkGraphics, magicnumber.Find(r))
	sign, err := magicnumber.Image(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.PortableNetworkGraphics, sign)
}

func TestWebp(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(webpFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Webp(r))
	assert.Equal(t, magicnumber.GoogleWebP, magicnumber.Find(r))
}

func TestWave(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(wavFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Wave(r))
	assert.Equal(t, magicnumber.WaveAudioForWindows, magicnumber.Find(r))
}

func TestMP3(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(mp3File))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Mp3(r))
	assert.Equal(t, magicnumber.MPEG1AudioLayer3, magicnumber.Find(r))
}

func TestOGG(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(oggFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Ogg(r))
	assert.Equal(t, magicnumber.OggVorbisCodec, magicnumber.Find(r))
}

func TestWMA(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file(wmaFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Wmv(r))
	assert.Equal(t, magicnumber.MicrosoftWindowsMedia, magicnumber.Find(r))
	sign, err := magicnumber.Video(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.MicrosoftWindowsMedia, sign)
}

func TestFlac(t *testing.T) {
	t.Parallel()
	r, err := os.Open(mp3file("TEST.flac"))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Flac(r))
	assert.Equal(t, magicnumber.FreeLosslessAudioCodec, magicnumber.Find(r))
}
