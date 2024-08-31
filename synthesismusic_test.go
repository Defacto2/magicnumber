package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMod(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(modFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.MK(r))
	assert.Equal(t, magicnumber.MusicProTracker, magicnumber.Find(r))
	assert.False(t, magicnumber.MTM(r))
	assert.Equal(t, "ProTracker 8-channel song", magicnumber.MusicTracker(r))
}

func TestXM(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(xmFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.XM(r))
	assert.Equal(t, magicnumber.MusicExtendedModule, magicnumber.Find(r))
	assert.False(t, magicnumber.IT(r))
	assert.Equal(t, "extended module tracked music", magicnumber.MusicTracker(r))
}

func TestIT(t *testing.T) {
	t.Parallel()
	r, err := os.Open(uncompress(itFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.IT(r))
	assert.Equal(t, magicnumber.MusicImpulseTracker, magicnumber.Find(r))
	assert.False(t, magicnumber.MK(r))
	assert.Equal(t, `Impulse Tracker song, "Defacto2 IT test fil"`,
		magicnumber.MusicTracker(r))
}
