package lnkfile

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

//uuid(000214F9-0000-0000-C000-000000000046)
type IShellLink struct {
	ole.IUnknown
}
type IShellLinkVtbl struct {
	ole.IUnknownVtbl
	GetPath             uintptr
	GetIDList           uintptr
	SetIDList           uintptr
	GetDescription      uintptr
	SetDescription      uintptr
	GetWorkingDirectory uintptr
	SetWorkingDirectory uintptr
	GetArguments        uintptr
	SetArguments        uintptr
	GetHotkey           uintptr
	SetHotkey           uintptr
	GetShowCmd          uintptr
	SetShowCmd          uintptr
	GetIconLocation     uintptr
	SetIconLocation     uintptr
	SetRelativePath     uintptr
	Resolve             uintptr
	SetPath             uintptr
}

func (v *IShellLink) VTable() *IShellLinkVtbl {
	return (*IShellLinkVtbl)(unsafe.Pointer(v.RawVTable))
}
func (v *IShellLink) GetPath(fFlags uint32) (pszFile string, err error) {
	pszFile, err = getPath(v, fFlags)
	return
}
func (v *IShellLink) GetIDList() (ppidl string, err error) {
	ppidl, err = getIDList(v)
	return
}
func (v *IShellLink) SetIDList(t string) (err error) {
	pidl := Cov_string_PCIDLIST_ABSOLUTE(t)
	err = setIDList(v, pidl)
	return
}
func (v *IShellLink) GetDescription() (pszName string, err error) {
	pszName, err = getDescription(v)
	return
}
func (v *IShellLink) SetDescription(s string) (err error) {
	pszName := Cov_string_LPCWSTR(s)
	err = setDescription(v, pszName)
	return
}
func (v *IShellLink) GetWorkingDirectory() (pszDir string, err error) {
	pszDir, err = getWorkingDirectory(v)
	return
}
func (v *IShellLink) SetWorkingDirectory(s string) (err error) {
	pszDir := Cov_string_LPCWSTR(s)
	err = setWorkingDirectory(v, pszDir)
	return
}
func (v *IShellLink) GetArguments(cch int32) (pszArgs string, err error) {
	pszArgs, err = getArguments(v)
	return
}
func (v *IShellLink) SetArguments(s string) (err error) {
	pszArgs := Cov_string_LPCWSTR(s)
	err = setArguments(v, pszArgs)
	return
}
func (v *IShellLink) GetHotkey() (pwHotkey *uint16, err error) {
	pwHotkey, err = getHotkey(v)
	return
}
func (v *IShellLink) SetHotkey(wHotkey uint16) (err error) {
	err = setHotkey(v, wHotkey)
	return
}
func (v *IShellLink) GetShowCmd() (piShowCmd *int, err error) {
	piShowCmd, err = getShowCmd(v)
	return
}
func (v *IShellLink) SetShowCmd(iShowCmd int) (err error) {
	err = setShowCmd(v, iShowCmd)
	return
}
func (v *IShellLink) GetIconLocation() (pszIconPath string, piIcon *int, err error) {
	pszIconPath, piIcon, err = getIconLocation(v)
	return
}
func (v *IShellLink) SetIconLocation(s string, iIcon int32) (err error) {
	pszIconPath := Cov_string_LPCWSTR(s)
	err = setIconLocation(v, pszIconPath, iIcon)
	return
}
func (v *IShellLink) SetRelativePath(s string, dwReserved uint32) (err error) {
	pszPathRel := Cov_string_LPCWSTR(s)
	err = setRelativePath(v, pszPathRel, dwReserved)
	return
}
func (v *IShellLink) Resolve(hwnd HWND, fFlags uint32) (err error) {
	err = resolve(v, hwnd, fFlags)
	return
}
func (v *IShellLink) SetPath(s string) (err error) {
	pszFile := Cov_string_LPCWSTR(s)
	err = setPath(v, pszFile)
	return
}

func (v *IShellLink) QueryInterface(iid *ole.GUID) (unk *ole.IUnknown, err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().QueryInterface,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(&unk)))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}

// IShellLink::GetPath fFlags
//SLGP_FLAGS
const (
	SLGP_SHORTPATH        = 0x0001
	SLGP_UNCPRIORITY      = 0x0002
	SLGP_RAWPATH          = 0x0004
	SLGP_RELATIVEPRIORITY = 0x0008
)

func GetIShellLinkFromIUnkown(u *ole.IUnknown) *IShellLink {
	return (*IShellLink)(unsafe.Pointer(u))
}
