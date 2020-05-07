package lnkfile

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

func isDirty(v *IPersistFile) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().IsDirty,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func load(v *IPersistFile, pszFileName LPCOLESTR, dwMode uint32) (err error) {
	hr, _, _ := syscall.Syscall(v.VTable().Load,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszFileName)),
		uintptr(dwMode))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func save(v *IPersistFile, pszFileName LPCOLESTR, fRemember BOOL) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().Save,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszFileName)),
		uintptr(fRemember))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func saveCompleted(v *IPersistFile, pszFileName LPCOLESTR) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().SaveCompleted,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pszFileName)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
func getCurFile(v *IPersistFile) (ppszFileName *LPOLESTR, err error) {
	ppszFileName = new(LPOLESTR)
	hr, _, _ := syscall.Syscall(
		v.VTable().GetCurFile,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(ppszFileName)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
