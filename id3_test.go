package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	IDv1File = "id3v1_001_basic.mp3"
	IDv2File = "id3v2_001_basic.mp3"
)

func TestMusicID3v1(t *testing.T) {
	t.Parallel()
	t.Log("TestMusicID3v1")
	r, err := os.Open(mp3file(IDv1File))
	require.NoError(t, err)
	defer r.Close()
	assert.Equal(t, "Title by Artist (2003)", magicnumber.MusicID3v1(r))
	assert.Equal(t, "", magicnumber.MusicID3v2(r))
}

func TestMusicID3v2(t *testing.T) {
	t.Parallel()
	t.Log("TestMusicID3v2")
	r, err := os.Open(mp3file(IDv2File))
	require.NoError(t, err)
	defer r.Close()
	assert.Equal(t, "", magicnumber.MusicID3v1(r))
	assert.Equal(t, "Title by Artist (2003)", magicnumber.MusicID3v2(r))
}

func TestConvSize(t *testing.T) {
	t.Parallel()
	t.Log("TestConvSize")
	assert.Equal(t, int64(257), magicnumber.ConvSize([]byte{0, 0, 0x02, 0x01}))
	assert.Equal(t, int64(742), magicnumber.ConvSize([]byte{0, 0, 0x05, 0x66}))
}