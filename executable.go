package magicnumber

// Package file executable.go contains the functions that parse Microsoft and IBM system executable files.

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// Pklite matches the PKLITE archive format in the byte slice which is a
// compressed executable format for DOS and 16-bit Windows.
func Pklite(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [6]byte
	const off = 30
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 6 {
		return false
	}
	return p == [6]byte{'P', 'K', 'L', 'I', 'T', 'E'}
}

// Pksfx matches the PKSFX archive format in the byte slice which is a
// self-extracting archive format.
func Pksfx(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [5]byte
	const off = 526
	if n, err := r.ReadAt(p[:], off); (err != nil && err != io.EOF) || n < 5 {
		return false
	}
	return p == [5]byte{'P', 'K', 'S', 'F', 'X'}
}

// DosKWAJ returns true if the reader begins with the KWAJ compression signature,
// found in some DOS executables.
func DosKWAJ(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [8]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 8 {
		return false
	}
	return p == [8]byte{'K', 'W', 'A', 'J', 0x88, 0xf0, 0x27, 0xd1}
}

// DosSZDD returns true if the reader begins with the SZDD compression signature.
func DosSZDD(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [8]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 8 {
		return false
	}
	return p == [8]byte{'S', 'Z', 'D', 'D', 0x88, 0xf0, 0x27, 0x33}
}

// MSExe returns true if the reader begins with the Microsoft executable signature.
func MSExe(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [2]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 2 {
		return false
	}

	return (p[0] == 'M' && p[1] == 'Z') || (p[0] == 'Z' && p[1] == 'M')
}

// MSComp returns true if the reader contains the Microsoft Compound File signature.
func MSComp(r io.ReaderAt) bool {
	if r == nil {
		return false
	}

	var p [8]byte
	if n, err := r.ReadAt(p[:], 0); (err != nil && err != io.EOF) || n < 8 {
		return false
	}
	return p == [8]byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}
}

// Windows represents the Windows specific information in the executable header.
type Windows struct {
	TimeDateStamp time.Time          // The time the executable was compiled, only included in PE files
	Major         int                // Major minimum version, for example, Windows 3.0 would be 3
	Minor         int                // Minor minimum version, for example, Windows 3.0 would be 0
	NE            NewExecutable      // The New Executable, a legacy format replaced by the Portable Executable format
	PE            PortableExecutable // The Portable Executable CPU architecture
	PE64          bool               // True if the executable is a 64-bit Portable Executable (PE32+)
}

func (w Windows) String() string {
	const (
		Windows2x   = 2
		WindowsNTv3 = 3
		WindowsNT   = 4
	)
	switch w.NE {
	case DOSv4Exe, OS2Exe:
		return fmt.Sprintf("%s v%d.%d", w.NE, w.Major, w.Minor)
	case UnknownNE:
		return "Unknown NE executable"
	case Windows286Exe:
		if w.Major == Windows2x {
			return fmt.Sprintf("Windows/286 v%d.%d", w.Major, w.Minor)
		}
		return fmt.Sprintf("Windows v%d.%d for 286", w.Major, w.Minor)
	case Windows386Exe:
		if w.Major == Windows2x {
			return fmt.Sprintf("Windows/386 v%d.%d", w.Major, w.Minor)
		}
		return fmt.Sprintf("Windows v%d.%d for 386+", w.Major, w.Minor)
	}
	if w.PE == Intel386PE {
		if w.Major < WindowsNTv3 {
			return "Windows 95/98/ME"
		}
		if w.Major <= WindowsNT {
			return fmt.Sprintf("Windows NT v%d.%d", w.Major, w.Minor)
		}
	}

	os := fmt.Sprintf("Windows NT v%d.%d", w.Major, w.Minor)
	for name, v := range windows {
		major, minor := v[0], v[1]
		if w.Major == major && w.Minor == minor {
			os = name
			break
		}
	}

	return pe(w.PE, w.PE64, os)
}

func pe(pe PortableExecutable, pe64 bool, os string) string {
	switch pe {
	case UnknownPE:
		if pe64 {
			return "Unknown PE+ executable"
		}
		return "Unknown PE executable"
	case Intel386PE:
		return os + " 32-bit"
	case AMD64PE:
		return os + " 64-bit"
	case ARMPE:
		return os + " for ARM"
	case ARM64PE:
		return os + " for ARM64"
	case ItaniumPE:
		return os + " for Itanium"
	default:
		// safe fallback for obscure/unlisted architectures
		return fmt.Sprintf("%s (%v)", os, pe)
	}
}

// WindowsName represents the Windows version names and their minimum version numbers.
type WindowsName map[string][2]int

// WindowsNames returns the Windows version names and their minimum version numbers.
// The minimum version numbers are based on the minimum system version required by the executable,
// and not the libraries or system calls in use by the program.
//
// The minimum version numbers were discontinued by Microsoft in Windows 8.1 and
// may not be accurate for modern programs.
func WindowsNames() WindowsName {
	return windows
}

var windows = WindowsName{ //nolint:gochecknoglobals
	"Windows 2000":                        {5, 0},
	"Windows XP":                          {5, 1},
	"Windows XP Professional x64 Edition": {5, 2},
	"Windows Vista":                       {6, 0},
	"Windows 7":                           {6, 1},
	"Windows 8":                           {6, 2},
	"Windows 8.1":                         {6, 3},
	"Windows 10":                          {10, 0},
}

// NewExecutable represents the New Executable file type, a format used by Microsoft and IBM
// from the mid-1980s to improve on the limitations of the MS-DOS MZ executable format.
type NewExecutable int

const (
	NoneNE        NewExecutable = iota - 1 // Not a New Executable
	UnknownNE                              // Unknown New Executable
	OS2Exe                                 // Microsoft IBM OS/2 New Executable
	Windows286Exe                          // Windows requiring an Intel 286 CPU New Executable
	DOSv4Exe                               // MS-DOS v4 New Executable
	Windows386Exe                          // Windows requiring an Intel 386 CPU New Executable
)

func (ne NewExecutable) String() string {
	switch ne {
	case NoneNE:
		return "Not a New Executable"
	case UnknownNE:
		return "Unknown New Executable"
	case OS2Exe:
		return "OS/2 New Executable"
	case Windows286Exe:
		return "Windows for 286 New Executable"
	case DOSv4Exe:
		return "MS-DOS v4 New Executable"
	case Windows386Exe:
		return "Windows for 386+ New Executable"
	}
	return ""
}

// PortableExecutable represents the Portable Executable file type, a format used by Microsoft
// for executables, object code, DLLs, FON Font files, and others. In this implementation, only
// executables for desktop Windows are considered.
type PortableExecutable uint16

const (
	UnknownPE  PortableExecutable = 0x0    // Unknown Portable Executable
	Intel386PE PortableExecutable = 0x14c  // Intel 386 Portable Executable
	AMD64PE    PortableExecutable = 0x8664 // AMD64 Portable Executable
	ARMPE      PortableExecutable = 0x1c0  // ARM Portable Executable
	ARM64PE    PortableExecutable = 0xaa64 // ARM64 Portable Executable
	ItaniumPE  PortableExecutable = 0x200  // Itanium Portable Executable
)

// FindExecutable reads the first 3KB from the reader and returns the specific information contained
// within the executable headers. Both the New Executable and Portable Executable formats are supported,
// which are commonly used by IBM and Microsoft desktop operating systems from PC/MS-DOS to modern Windows.
func FindExecutable(r io.ReaderAt) (Windows, error) {
	const format = "find executable first %d bytes: %w"
	win := Default()
	if r == nil {
		return win, ErrNilReader
	}

	const size = 1024 * 3
	var p [size]byte
	n, err := r.ReadAt(p[:], 0)
	if err != nil && err != io.EOF {
		return win, fmt.Errorf(format, size, err)
	}
	if n == 0 {
		return win, fmt.Errorf(format, size, io.ErrUnexpectedEOF)
	}

	s := p[:n]
	win = NE(s)
	if win.NE == NoneNE {
		win = PE(s)
	}
	return win, nil
}

func Default() Windows {
	return Windows{
		Major:         0,
		Minor:         0,
		TimeDateStamp: time.Time{},
		PE64:          false,
		PE:            UnknownPE,
		NE:            NoneNE,
	}
}

// NE returns the New Executable file type from the byte slice.
//
// Windows programs that are New Executables are usually for the ancient Windows 2 or 3.x editions.
// Windows v2 came in two versions, Windows 2 (for the 286 CPU) and Windows/386,
// while Windows 3.0+ unified support for both CPUs.
// The New Executable format was replaced by the Portable Executable format in Windows 95/NT.
//
// If a Windows program is detected, the major and minor version numbers are returned,
// for example, a Windows 3.0 requirement would return 3 and 0.
func NE(p []byte) Windows {
	none := Default()

	// DOS header must be at least 62 bytes to contain e_lfanew at 0x3c
	const minDOSHead = 0x3c + 2
	if len(p) < minDOSHead {
		return none
	}
	if p[0] != 'M' || p[1] != 'Z' {
		return none
	}

	const (
		segmentedHeaderIndex = 0x3c     // the location of the segmented header
		executableTypeIndex  = 0x36     // the executable type aka the operating system
		winMinorIndex        = 0x3e     // the location of the Windows minor version
		winMajorIndex        = 0x3f     // the location of the Windows major version
		minNEHead            = 0x3f + 1 // minimum bytes required relative to NE header offset
	)

	offset := int(binary.LittleEndian.Uint16(p[segmentedHeaderIndex:]))
	if len(p) < offset+minNEHead {
		return none
	}
	if p[offset] != 'N' || p[offset+1] != 'E' {
		return none
	}
	newExec := NewExecutable(p[offset+executableTypeIndex])
	switch newExec {
	case Windows286Exe, Windows386Exe, OS2Exe, DOSv4Exe, UnknownNE:
		return Windows{
			Major: int(p[offset+winMajorIndex]),
			Minor: int(p[offset+winMinorIndex]),
			NE:    newExec,
		}
	}

	return none
}

// PE returns the Portable Executable file type from the byte slice.
//
// The [Portable Executable format] is used by Microsoft for executables, object code, DLLs, FON Font files, and others.
// In this implementation, only executables for desktop Windows are considered. The information returned is the
// CPU architecture, the Windows NT version, and the time the executable was compiled.
//
// The major and minor version numbers are not always accurate.
//
// [Portable Executable format]: https://learn.microsoft.com/en-us/windows/win32/debug/pe-format
func PE(p []byte) Windows {
	none := Default()

	// DOS header must be at least 64 bytes to contain e_lfanew at 0x3c
	const minDOSHead = 0x3c + 4
	if len(p) < minDOSHead {
		return none
	}
	if p[0] != 'M' || p[1] != 'Z' {
		return none
	}

	const peHead = 0x3c
	off := int(binary.LittleEndian.Uint32(p[peHead:]))
	const peLen = 4 + 20 + 2
	if off < 0 || len(p) < off+peLen {
		return none
	}
	if p[off] != 'P' || p[off+1] != 'E' || p[off+2] != 0 || p[off+3] != 0 {
		return none
	}

	const index = 4
	coffIndex := off + index
	machine := binary.LittleEndian.Uint16(p[coffIndex:])
	sec := int64(binary.LittleEndian.Uint32(p[coffIndex+4:]))
	compiled := time.Unix(sec, 0)

	const coffHead = 20
	optIndex := coffIndex + coffHead
	magic := binary.LittleEndian.Uint16(p[optIndex:])

	const (
		majorOff = 40
		minorOff = 42
	)
	if len(p) < optIndex+minorOff+2 {
		return none
	}
	major := int(binary.LittleEndian.Uint16(p[optIndex+majorOff:]))
	minor := int(binary.LittleEndian.Uint16(p[optIndex+minorOff:]))

	const pe32Plus = 0x020b
	return Windows{
		Major:         major,
		Minor:         minor,
		TimeDateStamp: compiled,
		PE:            exec(machine),
		PE64:          magic == pe32Plus,
		NE:            NoneNE,
	}
}

func exec(pem uint16) PortableExecutable {
	pe := PortableExecutable(pem)
	switch pe {
	case Intel386PE, AMD64PE, ARMPE, ARM64PE, ItaniumPE:
		return pe
	default:
		return UnknownPE
	}
}
