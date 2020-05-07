package lnkfile

import (
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type WCHAR uint16
type LPCSTR *WCHAR
type LPCWSTR *WCHAR
type LPWSTR *WCHAR
type LPOLESTR *WCHAR
type LPCOLESTR *WCHAR
type LPWCHAR *WCHAR

const _uint16_slice_limit = 100000

func getUint16ptrLen(uu interface{}) int {
	u := uu.(*uint16)
	n := int(0)
	c := uint16(0)
	for {
		if n > _uint16_slice_limit {
			break
		}
		c = *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + uintptr(n)))
		if c == 0 {
			break
		} else {
			n++
		}
	}
	return n
}
func Cov_LPWCHAR_string(p LPWCHAR, max int) (r string) {
	r = UTF16PtrToString(p, max)
	return
}
func Cov_LPWSTR_string(p LPWSTR) (r string) {
	r = UTF16PtrToString(p, getUint16ptrLen(p))
	return
}

func Cov_string_LPCOLESTR(s string) (r LPCOLESTR) {
	u, _ := syscall.UTF16PtrFromString(s)
	r = (LPCOLESTR)(unsafe.Pointer(u))
	runtime.KeepAlive(u)
	return
}
func Cov_pLPCOLESTR_string(s *LPCOLESTR) (r string) {
	t := *s
	r = UTF16PtrToString(t, getUint16ptrLen(t))
	return
}
func Cov_pLPOLESTR_string(s *LPOLESTR) (r string) {
	t := *s
	r = UTF16PtrToString(t, getUint16ptrLen(t))
	return
}

func Cov_string_LPCSTR(s string) (r LPCSTR) {
	u, _ := syscall.UTF16PtrFromString(s)
	r = (LPCSTR)(unsafe.Pointer(u))
	return
}
func Cov_string_LPCWSTR(s string) (r LPCWSTR) {
	u, _ := syscall.UTF16PtrFromString(s)
	r = (LPCWSTR)(unsafe.Pointer(u))
	return
}

// UTF16PtrToString is like UTF16ToString, but takes *uint16
// as a parameter instead of []uint16.
// max is how many times p can be advanced looking for the null terminator.
// If max is hit, the string is truncated at that point.
func UTF16PtrToString(p *WCHAR, max int) string {
	if p == nil {
		return ""
	}
	// Find NUL terminator.
	end := unsafe.Pointer(p)
	n := 0
	for *(*uint16)(end) != 0 && n < max {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*p))
		n++
	}
	s := (*[(1 << 30) - 1]uint16)(unsafe.Pointer(p))[:n:n]
	return string(utf16.Decode(s))
}

type BOOL int

const TRUE = 1
const FALSE = 0

type DWORD uint32
type WORD uint16
type FILETIME struct {
	dwLowDateTime  DWORD
	dwHighDateTime DWORD
}

const MAX_PATH = 260

type WIN32_FIND_DATAW struct {
	dwFileAttributes   DWORD
	ftCreationTime     FILETIME
	ftLastAccessTime   FILETIME
	ftLastWriteTime    FILETIME
	nFileSizeHigh      DWORD
	nFileSizeLow       DWORD
	dwReserved0        DWORD
	dwReserved1        DWORD
	cFileName          [MAX_PATH]WCHAR
	cAlternateFileName [14]WCHAR
}

type USHORT uint16
type BYTE uint8

const SHITEMID_SIZE = 500

type SHITEMID struct {
	cb   USHORT
	abID [SHITEMID_SIZE]WCHAR // BYTE        abID[];
}
type ITEMIDLIST struct {
	SHITEMID
}

//typedef /* [wire_marshal] */ ITEMIDLIST __unaligned *LPITEMIDLIST;
type LPITEMIDLIST *ITEMIDLIST
type IDLIST_ABSOLUTE ITEMIDLIST
type PIDLIST_ABSOLUTE *IDLIST_ABSOLUTE

func Cov_PIDLIST_ABSOLUTE_string(p PIDLIST_ABSOLUTE) (r string) {
	r = UTF16PtrToString(&p.abID[0], int(p.cb)-1)
	return

}
func Cov_string_PIDLIST_ABSOLUTE(s string) (p PIDLIST_ABSOLUTE) {
	p = &IDLIST_ABSOLUTE{}
	sc, _ := syscall.UTF16FromString(s)
	n := len(sc)
	if n > SHITEMID_SIZE {
		panic("too long,please change SHITEMID_SIZE")
	}
	for i, v := range sc {
		p.abID[i] = WCHAR(v)
	}
	p.cb = USHORT(n*2 + 2)
	return
}
func Cov_string_PCIDLIST_ABSOLUTE(s string) (p PCIDLIST_ABSOLUTE) {
	p = &IDLIST_ABSOLUTE{}
	sc, _ := syscall.UTF16FromString(s)
	n := len(sc)
	if n > SHITEMID_SIZE {
		panic("too long,please change SHITEMID_SIZE")
	}
	for i, v := range sc {
		p.abID[i] = WCHAR(v)
	}
	p.cb = USHORT(n*2 + 2)
	return
}

type PCIDLIST_ABSOLUTE *IDLIST_ABSOLUTE

type HWND uintptr
