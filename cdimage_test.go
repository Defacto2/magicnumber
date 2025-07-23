package magicnumber_test

import (
	"os"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
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
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Daa(r))
	sign := magicnumber.Find(r)
	be.Equal(t, magicnumber.CDPowerISO, sign)
	be.Equal(t, "CD, PowerISO", sign.String())
	be.Equal(t, "CD PowerISO", sign.Title())
	b, sign, err := magicnumber.MatchExt(DaaFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, magicnumber.CDPowerISO, sign)
	sign, err = magicnumber.DiscImage(r)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.CDPowerISO, sign)
}

func TestCDISO(t *testing.T) {
	t.Parallel()
	t.Log("TestCDISO")
	r, err := os.Open(imgfile(ISOFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ISO(r))
	sign := magicnumber.Find(r)
	be.Equal(t, magicnumber.CDISO9660, sign)
	be.Equal(t, "CD, ISO 9660", sign.String())
	be.Equal(t, "CD ISO 9660", sign.Title())
	b, sign, err := magicnumber.MatchExt(ISOFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, magicnumber.CDISO9660, sign)
}

func TestMdf(t *testing.T) {
	t.Parallel()
	t.Log("TestMdf")
	r, err := os.Open(imgfile(MdfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Mdf(r))
	sign := magicnumber.Find(r)
	be.Equal(t, magicnumber.CDAlcohol120, sign)
	b, sign, err := magicnumber.MatchExt(DaaFile, r)
	be.Err(t, err, nil)
	be.True(t, !b)
	be.Equal(t, magicnumber.CDAlcohol120, sign)
}
