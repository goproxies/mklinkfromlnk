package lnkfile // import "github.com/goproxies/mklinkfromlnk/lnkfile"

import (
	"github.com/go-ole/go-ole"
	"unsafe"
)

//uuid(0000010c-0000-0000-C000-000000000046)
type IPersist struct {
	ole.IUnknown
}
type IPersistVtbl struct {
	ole.IUnknownVtbl
	GetClassID uintptr
}

func (v *IPersist) VTable() *IPersistVtbl {
	return (*IPersistVtbl)(unsafe.Pointer(v.RawVTable))
}
func (v *IPersist) GetClassID() (pClassID *ole.GUID, err error) {
	pClassID, err = getClassID(v)
	return
}

func GetIPersistFromIUnkown(u *ole.IUnknown) *IPersist {
	return (*IPersist)(unsafe.Pointer(u))
}
