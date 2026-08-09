package magicnumber_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Defacto2/magicnumber"
	"github.com/nalgeon/be"
)

func windows(t *testing.T, name string) string {
	t.Helper()
	return pathFile(t, filepath.Join("binaries", "windows", name))
}

func TestMSExe(t *testing.T) {
	t.Parallel()
	t.Log("TestMSExe")
	r, err := os.Open(windows(t, "hellojs.com"))
	be.Err(t, err, nil)
	defer r.Close()
	be.True(t, magicnumber.MSExe(r))
}

func TestFindBytesExecutableFreeDOS(t *testing.T) {
	t.Parallel()
	got, err := magicnumber.FindExecutable(nil)
	be.Err(t, err)
	be.Equal(t, got.PE, magicnumber.UnknownPE)
	be.Equal(t, got.NE, magicnumber.NoneNE)

	freedos := [4]string{
		filepath.Join("exe", "EXE.EXE"),
		filepath.Join("exemenu", "exemenu.exe"),
		filepath.Join("press", "PRESS.EXE"),
		filepath.Join("rread", "rread.exe"),
	}
	for n, name := range freedos[:] {
		t.Log(n, name)
		p, err := os.Open(pathFile(t, filepath.Join("binaries", "freedos", name)))
		be.Err(t, err, nil)
		defer p.Close()
		got, err = magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, got.PE, magicnumber.UnknownPE)
		be.Equal(t, got.NE, magicnumber.NoneNE)
		got, err := magicnumber.Program(p)
		be.Err(t, err, nil)
		be.Equal(t, got, magicnumber.MicrosoftExecutable)
	}
}

func TestFindBytesExecutableWinVista(t *testing.T) {
	vista := [3]string{
		"hello.com",
		"hellojs.com",
		"life.com",
	}
	for n, name := range vista[:] {
		t.Log(n, name)
		p, err := os.Open(pathFile(t, filepath.Join("binaries", "windows", name)))
		be.Err(t, err, nil)
		defer p.Close()
		be.Err(t, err, nil)
		w, err := magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, w.PE, magicnumber.AMD64PE)
		be.Equal(t, w.Major, 6)
		be.Equal(t, w.Minor, 0)
		be.Equal(t, w.TimeDateStamp.Year(), 2019)
		be.Equal(t, fmt.Sprint(w), "Windows Vista 64-bit")
		be.Equal(t, w.NE, magicnumber.NoneNE)
		sign, err := magicnumber.Program(p)
		be.Err(t, err, nil)
		be.Equal(t, sign, magicnumber.MicrosoftExecutable)
	}
}

func TestFindBytesExecutableWin3(t *testing.T) {
	winv3 := [3]string{
		filepath.Join("calmir10", "CALMIRA.EXE"),
		filepath.Join("calmir10", "TASKBAR.EXE"),
		filepath.Join("dskutl21", "DISKUTIL.EXE"),
	}
	for n, name := range winv3[:] {
		t.Log(n, name)
		p, err := os.Open(pathFile(t, filepath.Join("binaries", "windows3x", name)))
		be.Err(t, err, nil)
		defer p.Close()
		win, err := magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, win.PE, magicnumber.UnknownPE)
		be.Equal(t, win.NE, magicnumber.Windows286Exe)
		be.Equal(t, win.NE.String(), "Windows for 286 New Executable")
		be.Equal(t, win.Major, 3)
		be.Equal(t, win.Minor, 10)
		be.Equal(t, fmt.Sprint(win), "Windows v3.10 for 286")
		got, err := magicnumber.Program(p)
		be.Err(t, err, nil)
		be.Equal(t, got, magicnumber.MicrosoftExecutable)
	}

	t.Log("32-bit Core Temp.exe")
	p, err := os.Open(pathFile(t, filepath.Join("binaries", "windowsXP", "CoreTempv13", "32bit", "Core Temp.exe")))
	be.Err(t, err, nil)
	defer p.Close()
	w, err := magicnumber.FindExecutable(p)
	be.Err(t, err, nil)
	be.Equal(t, w.PE, magicnumber.Intel386PE)
	be.Equal(t, w.NE, magicnumber.NoneNE)
	be.Equal(t, w.Major, 5)
	be.Equal(t, w.Minor, 0)
	be.Equal(t, fmt.Sprint(w), "Windows 2000 32-bit")
	sign, err := magicnumber.Program(p)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.MicrosoftExecutable)

	t.Log("64-bit Core Temp.exe")
	p, err = os.Open(pathFile(t, filepath.Join("binaries", "windowsXP", "CoreTempv13", "64bit", "Core Temp.exe")))
	be.Err(t, err, nil)
	defer p.Close()
	be.Err(t, err, nil)
	w, err = magicnumber.FindExecutable(p)
	be.Err(t, err, nil)
	be.Equal(t, w.PE, magicnumber.AMD64PE)
	be.Equal(t, w.NE, magicnumber.NoneNE)
	be.Equal(t, w.Major, 5)
	be.Equal(t, w.Minor, 2)
	be.Equal(t, fmt.Sprint(w), "Windows XP Professional x64 Edition 64-bit")
	sign, err = magicnumber.Program(p)
	be.Err(t, err, nil)
	be.Equal(t, sign, magicnumber.MicrosoftExecutable)
}

func TestFindExecutableWinNT(t *testing.T) {
	win9x := [5]string{
		filepath.Join("rlowe-encrypt", "DEMOCD.EXE"),
		filepath.Join("rlowe-encrypt", "DISKDVR.EXE"),
		filepath.Join("rlowe-cdrools", "DEMOCD.EXE"),
		filepath.Join("7za920", "7za.exe"),
		filepath.Join("7z1604-extra", "7za.exe"),
	}
	for n, name := range win9x[:] {
		t.Log(n, name)
		p, err := os.Open(pathFile(t, filepath.Join("binaries", "windows9x", name)))
		be.Err(t, err, nil)
		defer p.Close()
		w, err := magicnumber.FindExecutable(p)
		be.Err(t, err, nil)
		be.Equal(t, w.PE, magicnumber.Intel386PE)
		be.Equal(t, w.Major, 4)
		be.Equal(t, w.Minor, 0)
		gt := w.TimeDateStamp.Year() > 2000
		be.True(t, gt)
		be.Equal(t, fmt.Sprint(w), "Windows NT v4.0")
		be.Equal(t, w.NE, magicnumber.NoneNE)
	}
}

func TestFindExecutableWin9x(t *testing.T) {
	unknowns := [3]string{
		filepath.Join("rlowe-rformat", "RFORMATD.EXE"),
		filepath.Join("rlowe-encrypt", "DFMINST.COM"),
		filepath.Join("rlowe-encrypt", "UNINST.COM"),
	}
	for n, name := range unknowns[:] {
		t.Log(n, name)
		p, err := os.Open(pathFile(t, filepath.Join("binaries", "windows9x", name)))
		be.Err(t, err, nil)
		defer p.Close()
		w, _ := magicnumber.FindExecutable(p)
		be.Equal(t, w.PE, magicnumber.UnknownPE)
		be.Equal(t, w.Major, 0)
		be.Equal(t, w.Minor, 0)
		be.Equal(t, w.TimeDateStamp.Year(), 1)
		be.Equal(t, fmt.Sprint(w), "Unknown PE executable")
		be.Equal(t, w.NE, magicnumber.NoneNE)
	}

	name := "7za.exe"
	t.Log(name)
	r, err := os.Open(pathFile(t, filepath.Join("binaries", "windows9x", "7z1604-extra", "x64", name)))
	be.Err(t, err, nil)
	defer r.Close()
	got, err := magicnumber.FindExecutable(r)
	be.Err(t, err, nil)
	be.Equal(t, got.PE, magicnumber.AMD64PE)
	be.Equal(t, 4, got.Major, 4)
	be.Equal(t, 0, got.Minor, 0)
	be.Equal(t, 2016, got.TimeDateStamp.Year(), 2016)
	be.Equal(t, fmt.Sprint(got), "Windows NT v4.0 64-bit")
	be.Equal(t, got.NE, magicnumber.NoneNE)
}
