package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	DaaFile = "uncompress.daa"
	ISOFile = "uncompress.iso"
	MdfFile = "uncompress.bin"
)

func TestDaa(t *testing.T) {
	t.Parallel()
	t.Log("TestDaa")
	r, err := os.Open(imgfile(DaaFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Daa(r))
	sign := magicnumber.Find(r)
	assert.Equal(t, magicnumber.CDPowerISO, sign)
	assert.Equal(t, "CD, PowerISO", sign.String())
	assert.Equal(t, "CD PowerISO", sign.Title())
	b, sign, err := magicnumber.MatchExt(DaaFile, r)
	require.NoError(t, err)
	assert.True(t, b)
	assert.Equal(t, magicnumber.CDPowerISO, sign)
	sign, err = magicnumber.DiscImage(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.CDPowerISO, sign)
}

func TestCDISO(t *testing.T) {
	t.Parallel()
	t.Log("TestCDISO")
	r, err := os.Open(imgfile(ISOFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.ISO(r))
	sign := magicnumber.Find(r)
	assert.Equal(t, magicnumber.CDISO9660, sign)
	assert.Equal(t, "CD, ISO 9660", sign.String())
	assert.Equal(t, "CD ISO 9660", sign.Title())
	b, sign, err := magicnumber.MatchExt(ISOFile, r)
	require.NoError(t, err)
	assert.True(t, b)
	assert.Equal(t, magicnumber.CDISO9660, sign)
}

func TestMdf(t *testing.T) {
	t.Parallel()
	t.Log("TestMdf")
	r, err := os.Open(imgfile(MdfFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Mdf(r))
	sign := magicnumber.Find(r)
	assert.Equal(t, magicnumber.CDAlcohol120, sign)
	b, sign, err := magicnumber.MatchExt(DaaFile, r)
	require.NoError(t, err)
	assert.False(t, b)
	assert.Equal(t, magicnumber.CDAlcohol120, sign)
}
