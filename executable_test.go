package magicnumber_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func windows(name string) string {
	return tdfile(filepath.Join("binaries", "windows", name))
}

func TestMSExe(t *testing.T) {
	t.Parallel()
	t.Log("TestMSExe")
	r, err := os.Open(windows("hellojs.com"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.MSExe(r))
}

func TestFindBytesExecutableFreeDOS(t *testing.T) {
	t.Parallel()
	w, err := magicnumber.FindExecutable(nil)
	be.Err(t, err)
	be.Equal(t, magicnumber.UnknownPE, w.PE)
	be.Equal(t, magicnumber.NoneNE, w.NE)

	freedos := []string{
		filepath.Join("exe", "EXE.EXE"),
		filepath.Join("exemenu", "exemenu.exe"),
		filepath.Join("press", "PRESS.EXE"),
		filepath.Join("rread", "rread.exe"),
	}
	for _, v := range freedos {
		p, err := os.Open(tdfile(filepath.Join("binaries", "freedos", v)))
		be.Err(t, err, nil)
		defer p.Close()
		w, err = magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.UnknownPE, w.PE)
		be.Equal(t, magicnumber.NoneNE, w.NE)
		sign, err := magicnumber.Program(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.MicrosoftExecutable, sign)
	}
}

func TestFindBytesExecutableWinVista(t *testing.T) {
	vista := []string{
		"hello.com",
		"hellojs.com",
		"life.com",
	}
	for _, v := range vista {
		p, err := os.Open(tdfile(filepath.Join("binaries", "windows", v)))
		be.Err(t, err, nil)
		defer p.Close()
		be.Err(t, err, nil)
		w, err := magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.AMD64PE, w.PE)
		be.Equal(t, 6, w.Major)
		be.Equal(t, 0, w.Minor)
		be.Equal(t, 2019, w.TimeDateStamp.Year())
		be.Equal(t, "Windows Vista 64-bit", fmt.Sprint(w))
		be.Equal(t, magicnumber.NoneNE, w.NE)
		sign, err := magicnumber.Program(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.MicrosoftExecutable, sign)
	}
}

func TestFindBytesExecutableWin3(t *testing.T) {
	winv3 := []string{
		filepath.Join("calmir10", "CALMIRA.EXE"),
		filepath.Join("calmir10", "TASKBAR.EXE"),
		filepath.Join("dskutl21", "DISKUTIL.EXE"),
	}
	for _, v := range winv3 {
		p, err := os.Open(tdfile(filepath.Join("binaries", "windows3x", v)))
		be.Err(t, err, nil)
		defer p.Close()
		w, err := magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.UnknownPE, w.PE)
		be.Equal(t, magicnumber.Windows286Exe, w.NE)
		be.Equal(t, "Windows for 286 New Executable", w.NE.String())
		be.Equal(t, 3, w.Major)
		be.Equal(t, 10, w.Minor)
		be.Equal(t, "Windows v3.10 for 286", fmt.Sprint(w))
		sign, err := magicnumber.Program(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.MicrosoftExecutable, sign)
	}

	p, err := os.Open(tdfile(filepath.Join("binaries", "windowsXP", "CoreTempv13", "32bit", "Core Temp.exe")))
	be.Err(t, err, nil)
	defer p.Close()
	w, err := magicnumber.FindExecutable(p)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.Intel386PE, w.PE)
	be.Equal(t, magicnumber.NoneNE, w.NE)
	be.Equal(t, 5, w.Major)
	be.Equal(t, 0, w.Minor)
	be.Equal(t, "Windows 2000 32-bit", fmt.Sprint(w))
	sign, err := magicnumber.Program(p)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.MicrosoftExecutable, sign)

	p, err = os.Open(tdfile(filepath.Join("binaries", "windowsXP", "CoreTempv13", "64bit", "Core Temp.exe")))
	be.Err(t, err, nil)
	defer p.Close()
	be.Err(t, err, nil)
	w, err = magicnumber.FindExecutable(p)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.AMD64PE, w.PE)
	be.Equal(t, magicnumber.NoneNE, w.NE)
	be.Equal(t, 5, w.Major)
	be.Equal(t, 2, w.Minor)
	be.Equal(t, "Windows XP Professional x64 Edition 64-bit", fmt.Sprint(w))
	sign, err = magicnumber.Program(p)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.MicrosoftExecutable, sign)
}

func TestFindExecutableWinNT(t *testing.T) {
	win9x := []string{
		filepath.Join("rlowe-encrypt", "DEMOCD.EXE"),
		filepath.Join("rlowe-encrypt", "DISKDVR.EXE"),
		filepath.Join("rlowe-cdrools", "DEMOCD.EXE"),
		filepath.Join("7za920", "7za.exe"),
		filepath.Join("7z1604-extra", "7za.exe"),
	}
	for _, v := range win9x {
		p, err := os.Open(tdfile(filepath.Join("binaries", "windows9x", v)))
		be.Err(t, err, nil)
		defer p.Close()
		w, err := magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, magicnumber.Intel386PE, w.PE)
		be.Equal(t, 4, w.Major)
		be.Equal(t, 0, w.Minor)
		gt := w.TimeDateStamp.Year() > 2000
		be.True(t, gt)
		be.Equal(t, "Windows NT v4.0", fmt.Sprint(w))
		be.Equal(t, magicnumber.NoneNE, w.NE)
	}
}

func TestFindExecutableWin9x(t *testing.T) {
	unknown := []string{
		filepath.Join("rlowe-rformat", "RFORMATD.EXE"),
		filepath.Join("rlowe-encrypt", "DFMINST.COM"),
		filepath.Join("rlowe-encrypt", "UNINST.COM"),
	}
	for _, v := range unknown {
		p, err := os.Open(tdfile(filepath.Join("binaries", "windows9x", v)))
		be.Err(t, err, nil)
		defer p.Close()
		w, _ := magicnumber.FindExecutable(p)
		be.Equal(t, magicnumber.UnknownPE, w.PE)
		be.Equal(t, 0, w.Major)
		be.Equal(t, 0, w.Minor)
		be.Equal(t, 1, w.TimeDateStamp.Year())
		be.Equal(t, "Unknown PE executable", fmt.Sprint(w))
		be.Equal(t, magicnumber.NoneNE, w.NE)
	}

	p, err := os.Open(tdfile(filepath.Join("binaries", "windows9x", "7z1604-extra", "x64", "7za.exe")))
	be.Err(t, err, nil)
	defer p.Close()
	w, err := magicnumber.FindExecutable(p)
	be.Err(t, err, nil)
	be.Equal(t, magicnumber.AMD64PE, w.PE)
	be.Equal(t, 4, w.Major)
	be.Equal(t, 0, w.Minor)
	be.Equal(t, 2016, w.TimeDateStamp.Year())
	be.Equal(t, "Windows NT v4.0 64-bit", fmt.Sprint(w))
	be.Equal(t, magicnumber.NoneNE, w.NE)
}
