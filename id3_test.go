package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

const (
	IDv1File = "id3v1_001_basic.mp3"
	IDv2File = "id3v2_001_basic.mp3"
)

const wantFormat = "Title by Artist (2003)"

func TestMusicID3v1(t *testing.T) {
	t.Parallel()
	t.Log("Test MusicID3v1")
	r, err := os.Open(pathMp3(t, IDv1File))
	be.Err(t, err, nil)
	defer r.Close()
	be.Equal(t, magicnumber.MusicID3v1(r), wantFormat)
	be.Equal(t, magicnumber.MusicID3v2(r), "")
}

func TestMusicID3v2(t *testing.T) {
	t.Parallel()
	t.Log("Test MusicID3v2")
	r, err := os.Open(pathMp3(t, IDv2File))
	be.Err(t, err, nil)
	defer r.Close()
	be.Equal(t, magicnumber.MusicID3v1(r), "")
	be.Equal(t, magicnumber.MusicID3v2(r), wantFormat)
}

func TestSyncsafe(t *testing.T) {
	t.Parallel()
	t.Log("Test Syncsafe")
	want := int64(257)
	be.Equal(t, magicnumber.Syncsafe([]byte{0, 0, 0x02, 0x01}), want)
	want = int64(742)
	be.Equal(t, magicnumber.Syncsafe([]byte{0, 0, 0x05, 0x66}), want)
}
