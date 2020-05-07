package lnkfile

import (
	"github.com/go-ole/go-ole"
	"unsafe"
)

//uuid(0000010b-0000-0000-C000-000000000046)
type IPersistFile struct {
	IPersist
}
type IPersistFileVtbl struct {
	IPersistVtbl
	IsDirty       uintptr
	Load          uintptr
	Save          uintptr
	SaveCompleted uintptr
	GetCurFile    uintptr
}

func (v *IPersistFile) VTable() *IPersistFileVtbl {
	return (*IPersistFileVtbl)(unsafe.Pointer(v.RawVTable))
}
func (v *IPersistFile) IsDirty() (err error) {
	err = isDirty(v)
	return
}
func (v *IPersistFile) Load(pszFileName string, dwMode uint32) (err error) {
	err = load(v, Cov_string_LPCOLESTR(pszFileName), dwMode)
	return
}
func (v *IPersistFile) Save(pszFileName string, fRemember BOOL) (err error) {
	err = save(v, Cov_string_LPCOLESTR(pszFileName), fRemember)
	return
}
func (v *IPersistFile) SaveCompleted(pszFileName string) (err error) {
	err = saveCompleted(v, Cov_string_LPCOLESTR(pszFileName))
	return
}
func (v *IPersistFile) GetCurFile() (ppszFileName string, err error) {
	t, err := getCurFile(v)
	ppszFileName = Cov_pLPOLESTR_string(t)
	return
}

// consts used for functions
const (
	// Load dwMode
	/* Storage instantiation modes */
	STGM_DIRECT           = 0x00000000
	STGM_TRANSACTED       = 0x00010000
	STGM_SIMPLE           = 0x08000000
	STGM_READ             = 0x00000000
	STGM_WRITE            = 0x00000001
	STGM_READWRITE        = 0x00000002
	STGM_SHARE_DENY_NONE  = 0x00000040
	STGM_SHARE_DENY_READ  = 0x00000030
	STGM_SHARE_DENY_WRITE = 0x00000020
	STGM_SHARE_EXCLUSIVE  = 0x00000010
	STGM_PRIORITY         = 0x00040000
	STGM_DELETEONRELEASE  = 0x04000000
	STGM_NOSCRATCH        = 0x00100000
	STGM_CREATE           = 0x00001000
	STGM_CONVERT          = 0x00020000
	STGM_FAILIFTHERE      = 0x00000000
	STGM_NOSNAPSHOT       = 0x00200000
	STGM_DIRECT_SWMR      = 0x00400000
)

func GetIPersistFileFromIUnkown(u *ole.IUnknown) *IPersistFile {
	return (*IPersistFile)(unsafe.Pointer(u))
}
