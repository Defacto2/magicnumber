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
	r, err := os.Open(pathDisc(t, DaaFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Daa(r))

	sign := magicnumber.Find(r)
	be.Equal(t, sign, magicnumber.CDPowerISO)
	be.Equal(t, sign.String(), "CD, PowerISO")
	be.Equal(t, sign.Title(), "CD PowerISO")

	b, sign, err := magicnumber.MatchExt(DaaFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, sign, magicnumber.CDPowerISO)

	sign, err = magicnumber.DiscImage(r)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.CDPowerISO)
}

func TestCDISO(t *testing.T) {
	t.Parallel()

	t.Log("TestCDISO")
	r, err := os.Open(pathDisc(t, ISOFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.ISO(r))
	sign := magicnumber.Find(r)

	be.Equal(t, sign, magicnumber.CDISO9660)
	be.Equal(t, sign.String(), "CD, ISO 9660")
	be.Equal(t, sign.Title(), "CD ISO 9660")

	b, sign, err := magicnumber.MatchExt(ISOFile, r)
	be.Err(t, err, nil)
	be.True(t, b)
	be.Equal(t, sign, magicnumber.CDISO9660)
}

func TestMdf(t *testing.T) {
	t.Parallel()

	t.Log("TestMdf")
	r, err := os.Open(pathDisc(t, MdfFile))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.Mdf(r))

	sign := magicnumber.Find(r)
	be.Equal(t, sign, magicnumber.CDAlcohol120)
	b, sign, err := magicnumber.MatchExt(DaaFile, r)
	be.Err(t, err, nil)
	be.True(t, !b)
	be.Equal(t, sign, magicnumber.CDAlcohol120)
}
