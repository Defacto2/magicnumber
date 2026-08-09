package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
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
	pakFile        = "PAK100.PAK"
)

func TestPak(t *testing.T) {
	t.Parallel()
	t.Log("TestPak")
	r, err := os.Open(pathFile(t, pakFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pak(r))
}

func TestArchive(t *testing.T) {
	t.Parallel()
	t.Log("TestArchive")
	r, err := os.Open(pathFile(t, seaFile))
	be.Err(t, err, nil)
	defer r.Close()
	sign, err := magicnumber.Archive(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.ARChiveSEA, sign)
	be.Equal(t, sign.String(), "ARC by SEA")
	be.Equal(t, sign.Title(), "Archive by SEA")
}

func TestZipReduce(t *testing.T) {
	t.Parallel()
	t.Log("TestZipReduce")
	r, err := os.Open(pathFile(t, zipReduceFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.PkShrink(r))
}

func TestZipShrink(t *testing.T) {
	t.Parallel()
	t.Log("TestZipShrink")
	r, err := os.Open(pathFile(t, zipShrinkFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.PkShrink(r))
}

func TestZipImplode(t *testing.T) {
	t.Parallel()
	t.Log("TestZipImplode")
	r, err := os.Open(pathFile(t, zipImplodeFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pkzip(r))
}

func TestZipStore(t *testing.T) {
	t.Parallel()
	t.Log("TestZipStore")
	r, err := os.Open(pathFile(t, zipStoreFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Pkzip(r))
}

func TestTar(t *testing.T) {
	t.Parallel()
	t.Log("TestTar")
	r, err := os.Open(pathFile(t, tarFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Tar(r))
}

func TestRarv5(t *testing.T) {
	t.Parallel()
	t.Log("TestRarv5")
	r, err := os.Open(pathFile(t, rarv5File))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Rarv5(r))
}

func TestGzip(t *testing.T) {
	t.Parallel()
	t.Log("TestGzip")
	r, err := os.Open(pathFile(t, gzFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Gzip(r))
	r, err = os.Open(pathFile(t, b2zFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.Gzip(r))
}

func TestBzip2(t *testing.T) {
	t.Parallel()
	t.Log("TestBzip2")
	r, err := os.Open(pathFile(t, b2zFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Bzip2(r))
}

func TestX7z(t *testing.T) {
	t.Parallel()
	t.Log("TestX7z")
	r, err := os.Open(pathFile(t, x7zFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.X7z(r))
}

func TestXZ(t *testing.T) {
	t.Parallel()
	t.Log("TestXZ")
	r, err := os.Open(pathFile(t, xzFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.XZ(r))
}

func TestArcFree(t *testing.T) {
	t.Parallel()
	t.Log("Test wanting ArcFree, not ArcSEA")
	r, err := os.Open(pathFile(t, freeArcFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ArcFree(r))
	be.True(t, !magicnumber.ArcSEA(r))
	b, got, err := magicnumber.MatchExt(freeArcFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, got, magicnumber.FreeArc)
}

func TestArcSEA(t *testing.T) {
	t.Parallel()
	t.Log("Test wanting ArcSEA, not ArcFree")
	r, err := os.Open(pathFile(t, seaFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, !magicnumber.ArcFree(r))
	be.True(t, magicnumber.ArcSEA(r))
	b, got, err := magicnumber.MatchExt(seaFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, got, magicnumber.ARChiveSEA)
}

func TestLHA(t *testing.T) {
	t.Parallel()
	t.Log("TestLzhLha")
	r, err := os.Open(pathFile(t, lhaFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.LzhLha(r))
	got := magicnumber.Find(r)
	be.Equal(t, got, magicnumber.YoshiLHA)
	be.Equal(t, got.String(), "LHA by Yoshi")
	be.Equal(t, got.Title(), "Yoshi LHA")
	b, got, err := magicnumber.MatchExt(lhaFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, got, magicnumber.YoshiLHA)
}

func TestArj(t *testing.T) {
	t.Parallel()
	t.Log("TestArj")
	r, err := os.Open(pathFile(t, arjFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Arj(r))
}

func TestCab(t *testing.T) {
	t.Parallel()
	t.Log("TestCab")
	r, err := os.Open(pathFile(t, cabFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Cab(r))
}

func TestZoo(t *testing.T) {
	t.Parallel()
	t.Log("TestZoo")
	r, err := os.Open(pathFile(t, zooFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Zoo(r))
}

func TestRar(t *testing.T) {
	t.Parallel()
	t.Log("TestRar")
	r, err := os.Open(pathFile(t, rarFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Rar(r))
	got := magicnumber.Find(r)
	be.Equal(t, got, magicnumber.RoshalARchive)
}
