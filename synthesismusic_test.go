package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func TestMod(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(modFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.MK(r))
	be.Equal(t, magicnumber.MusicProTracker, magicnumber.Find(r))
	be.True(t, !magicnumber.MTM(r))
	be.Equal(t, "ProTracker 8-channel song", magicnumber.MusicTracker(r))
}

func TestXM(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(xmFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.XM(r))
	be.Equal(t, magicnumber.MusicExtendedModule, magicnumber.Find(r))
	be.True(t, !magicnumber.IT(r))
	be.Equal(t, "extended module tracked music", magicnumber.MusicTracker(r))
}

func TestIT(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(itFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.IT(r))
	be.Equal(t, magicnumber.MusicImpulseTracker, magicnumber.Find(r))
	be.True(t, !magicnumber.MK(r))
	be.Equal(t, `Impulse Tracker song, "Defacto2 IT test fil"`,
		magicnumber.MusicTracker(r))
}
