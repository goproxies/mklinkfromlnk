// +build windows

package dll

import (
	"github.com/goproxies/mklinkfromlnk/lnkfile"
	"golang.org/x/sys/windows"
	"unsafe"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
)
var (
	procMultiByteToWideChar = modkernel32.NewProc("MultiByteToWideChar")
)

func CovConsoleBytesToString(b []byte) string {
	n := len(b)
	if n == 0 {
		return ""
	}
	c := make([]lnkfile.WCHAR, n)
	//use CP_OEMCP
	hr, _, _ := procMultiByteToWideChar.Call(uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(n),
		uintptr(unsafe.Pointer(&c[0])),
		uintptr(n))
	if hr == 0 {
		return ""
	}
	return lnkfile.UTF16PtrToString(&c[0], n)
}
