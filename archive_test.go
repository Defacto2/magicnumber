package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	zipReduceFile  = "PKZ90B4.ZIP"
	zipShrinkFile  = "PKZ80A1.ZIP"
	zipImplodeFile = "PKZ110EI.ZIP"
	zipStoreFile   = "PKZ204E0.ZIP"
	freeArcFile    = "TESTfree.arc"
	seaFile        = "ARC521P.ARC"
	arjFile        = "ARJ310.ARJ"
	tarFile        = "TAR135.TAR"
	rarv5File      = "RAR624.RAR"
	gzFile         = "TAR135.GZ"
	b2zFile        = "TEST.tar.bz2"
	lhaFile        = "LHA114.LZH"
	x7zFile        = "TEST.7z"
	xzFile         = "TEST.tar.xz"
	cabFile        = "TEST.cab"
	zooFile        = "TEST.zoo"
	rarFile        = "TEST.rar"
	zip64File      = "TEST64.zip"
)

func TestArchive(t *testing.T) {
	t.Parallel()
	t.Log("TestArchive")
	r, err := os.Open(tdfile(seaFile))
	require.NoError(t, err)
	defer r.Close()
	sign, err := magicnumber.Archive(r)
	require.NoError(t, err)
	assert.Equal(t, magicnumber.ARChiveSEA, sign)
	assert.Equal(t, "ARC by SEA", sign.String())
	assert.Equal(t, "Archive by SEA", sign.Title())
}

func TestZipReduce(t *testing.T) {
	t.Parallel()
	t.Log("TestZipReduce")
	r, err := os.Open(tdfile(zipReduceFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.PkShrink(r))
}

func TestZipShrink(t *testing.T) {
	t.Parallel()
	t.Log("TestZipShrink")
	r, err := os.Open(tdfile(zipShrinkFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.PkShrink(r))
}

func TestZipImplode(t *testing.T) {
	t.Parallel()
	t.Log("TestZipImplode")
	r, err := os.Open(tdfile(zipImplodeFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Pkzip(r))
}

func TestZipStore(t *testing.T) {
	t.Parallel()
	t.Log("TestZipStore")
	r, err := os.Open(tdfile(zipStoreFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Pkzip(r))
}

func TestTar(t *testing.T) {
	t.Parallel()
	t.Log("TestTar")
	r, err := os.Open(tdfile(tarFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Tar(r))
}

func TestRarv5(t *testing.T) {
	t.Parallel()
	t.Log("TestRarv5")
	r, err := os.Open(tdfile(rarv5File))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Rarv5(r))
}

func TestGzip(t *testing.T) {
	t.Parallel()
	t.Log("TestGzip")
	r, err := os.Open(tdfile(gzFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Gzip(r))
	r, err = os.Open(tdfile(b2zFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.Gzip(r))
}

func TestBzip2(t *testing.T) {
	t.Parallel()
	t.Log("TestBzip2")
	r, err := os.Open(tdfile(b2zFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Bzip2(r))
}

func TestX7z(t *testing.T) {
	t.Parallel()
	t.Log("TestX7z")
	r, err := os.Open(tdfile(x7zFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.X7z(r))
}

func TestXZ(t *testing.T) {
	t.Parallel()
	t.Log("TestXZ")
	r, err := os.Open(tdfile(xzFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.XZ(r))
}

func TestArcFree(t *testing.T) {
	t.Parallel()
	t.Log("TestArcFree")
	r, err := os.Open(tdfile(freeArcFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.ArcFree(r))
	assert.False(t, magicnumber.ArcSEA(r))
	b, sign, err := magicnumber.MatchExt(freeArcFile, r)
	require.NoError(t, err)
	assert.True(t, b)
	assert.Equal(t, magicnumber.FreeArc, sign)
}

func TestArcSEA(t *testing.T) {
	t.Parallel()
	t.Log("TestArcSEA")
	r, err := os.Open(tdfile(seaFile))
	require.NoError(t, err)
	defer r.Close()
	assert.False(t, magicnumber.ArcFree(r))
	assert.True(t, magicnumber.ArcSEA(r))
	b, sign, err := magicnumber.MatchExt(seaFile, r)
	require.NoError(t, err)
	assert.True(t, b)
	assert.Equal(t, magicnumber.ARChiveSEA, sign)
}

func TestLHA(t *testing.T) {
	t.Parallel()
	t.Log("TestLzhLha")
	r, err := os.Open(tdfile(lhaFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.LzhLha(r))
	sign := magicnumber.Find(r)
	assert.Equal(t, magicnumber.YoshiLHA, sign)
	assert.Equal(t, "LHA by Yoshi", sign.String())
	assert.Equal(t, "Yoshi LHA", sign.Title())
	b, sign, err := magicnumber.MatchExt(lhaFile, r)
	require.NoError(t, err)
	assert.True(t, b)
	assert.Equal(t, magicnumber.YoshiLHA, sign)
}

func TestArj(t *testing.T) {
	t.Parallel()
	t.Log("TestArj")
	r, err := os.Open(tdfile(arjFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Arj(r))
}

func TestCab(t *testing.T) {
	t.Parallel()
	t.Log("TestCab")
	r, err := os.Open(tdfile(cabFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Cab(r))
}

func TestZoo(t *testing.T) {
	t.Parallel()
	t.Log("TestZoo")
	r, err := os.Open(tdfile(zooFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Zoo(r))
}

func TestRar(t *testing.T) {
	t.Parallel()
	t.Log("TestRar")
	r, err := os.Open(tdfile(rarFile))
	require.NoError(t, err)
	defer r.Close()
	assert.True(t, magicnumber.Rar(r))
	sign := magicnumber.Find(r)
	assert.Equal(t, magicnumber.RoshalARchive, sign)
}
