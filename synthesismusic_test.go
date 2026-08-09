package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func TestSynthMod(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(modFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.MK(r))
	be.Equal(t, magicnumber.Find(r), magicnumber.MusicProTracker)
	be.True(t, !magicnumber.MTM(r))
	const want = `ProTracker 8-channel song, "Defacto2 Test XM"`
	be.Equal(t, magicnumber.MusicTracker(r), want)
}

func TestSynthXM(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(xmFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.XM(r))
	be.Equal(t, magicnumber.MusicExtendedModule, magicnumber.Find(r))
	be.True(t, !magicnumber.IT(r))
	const want = "extended module tracked music"
	be.Equal(t, magicnumber.MusicTracker(r), want)
}

func TestSyncIT(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(itFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.IT(r))
	be.Equal(t, magicnumber.MusicImpulseTracker, magicnumber.Find(r))
	be.True(t, !magicnumber.MK(r))
	const want = `Impulse Tracker song, "Defacto2 IT test file"`
	be.Equal(t, magicnumber.MusicTracker(r), want)
}
