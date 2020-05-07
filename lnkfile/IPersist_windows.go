package lnkfile

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

func getClassID(v *IPersist) (pClassID *ole.GUID, err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().GetClassID,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(pClassID)),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
