package lnkfile

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

func getPath(v *IShellLink, fFlags uint32) (r string, err error) {
	pszFile := [MAX_PATH]WCHAR{}
	cch := MAX_PATH
	pfd := new(WIN32_FIND_DATAW)
	hr, _, _ := syscall.Syscall6(
		v.VTable().GetPath,
		5,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&pszFile[0])),
		uintptr(cch),
		uintptr(unsafe.Pointer(pfd)),
		uintptr(fFlags),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	r = UTF16PtrToString(&pszFile[0], cch)
	return
}
func getIDList(v *IShellLink) (s string, err error) {
	t := &IDLIST_ABSOLUTE{}
	ppidl := &t
	hr, _, _ := syscall.Syscall(
		v.VTable().GetIDList,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(ppidl)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	s = Cov_PIDLIST_ABSOLUTE_string(*ppidl)
	return
}
func setIDList(v *IShellLink, pidl PCIDLIST_ABSOLUTE) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetIDList,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pidl)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getDescription(v *IShellLink) (t string, err error) {
	pszName := [1000]WCHAR{}
	cch := 1000
	hr, _, _ := syscall.Syscall(
		v.VTable().GetDescription,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&pszName[0])),
		uintptr(cch))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	t = Cov_LPWCHAR_string(&pszName[0], cch)
	return
}
func setDescription(v *IShellLink, pszName LPCWSTR) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetDescription,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszName)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getWorkingDirectory(v *IShellLink) (s string, err error) {
	tp := [MAX_PATH]WCHAR{}
	cch := MAX_PATH
	pszDir := &tp[0]
	hr, _, _ := syscall.Syscall(
		v.VTable().GetWorkingDirectory,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszDir)),

		uintptr(cch))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	s = Cov_LPWSTR_string(pszDir)
	return
}
func setWorkingDirectory(v *IShellLink, pszDir LPCWSTR) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetWorkingDirectory,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszDir)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getArguments(v *IShellLink) (t string, err error) {
	tp := [MAX_PATH]WCHAR{}
	cch := MAX_PATH
	pszArgs := &tp[0]
	hr, _, _ := syscall.Syscall(
		v.VTable().GetArguments,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszArgs)),
		uintptr(cch))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	t = Cov_LPWSTR_string(pszArgs)
	return
}
func setArguments(v *IShellLink, pszArgs LPCWSTR) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetArguments,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszArgs)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getHotkey(v *IShellLink) (pwHotkey *uint16, err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().GetHotkey,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pwHotkey)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func setHotkey(v *IShellLink, wHotkey uint16) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetHotkey,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(wHotkey),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getShowCmd(v *IShellLink) (piShowCmd *int, err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().GetShowCmd,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(piShowCmd)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func setShowCmd(v *IShellLink, iShowCmd int) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetShowCmd,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(iShowCmd),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getIconLocation(v *IShellLink) (t string, piIcon *int, err error) {
	tp := [MAX_PATH]WCHAR{}
	cch := MAX_PATH
	pszIconPath := &tp[0]
	hr, _, _ := syscall.Syscall6(
		v.VTable().GetIconLocation,
		4,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszIconPath)),
		uintptr(cch),
		uintptr(unsafe.Pointer(piIcon)),
		0,
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	t = Cov_LPWSTR_string(pszIconPath)
	return
}
func setIconLocation(v *IShellLink, pszIconPath LPCWSTR, iIcon int32) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetIconLocation,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszIconPath)),
		uintptr(iIcon))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func setRelativePath(v *IShellLink, pszPathRel LPCWSTR, dwReserved uint32) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetRelativePath,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszPathRel)),
		uintptr(dwReserved))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func resolve(v *IShellLink, hwnd HWND, fFlags uint32) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().Resolve,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(hwnd),
		uintptr(fFlags))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func setPath(v *IShellLink, pszFile LPCWSTR) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SetPath,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszFile)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
