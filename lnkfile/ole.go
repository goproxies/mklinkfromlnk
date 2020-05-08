package lnkfile

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/goproxies/mklinkfromlnk/errors"
	"path/filepath"
	"strings"
)

var (

	// CLSID_ShellLink
	//[ uuid(00021401-0000-0000-C000-000000000046) ] coclass ShellLink { interface IShellLinkW; }
	CLSID_ShellLink *ole.GUID

	//uuid(000214F9-0000-0000-C000-000000000046) IShellLinkW
	IID_IShellLink *ole.GUID

	//uuid(0000010b-0000-0000-C000-000000000046)
	IID_IPersistFile *ole.GUID
	oIPersistFile    *IPersistFile
	oIShellLink      *IShellLink
)

func init() {
	defer errors.OleCatchErrors(false, "Init COM")
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	CLSID_ShellLink = ole.NewGUID("00021401-0000-0000-C000-000000000046")
	IID_IShellLink = ole.NewGUID("000214F9-0000-0000-C000-000000000046")
	IID_IPersistFile = ole.NewGUID("0000010b-0000-0000-C000-000000000046")
	u, err := ole.CreateInstance(CLSID_ShellLink, IID_IShellLink)
	errors.OlePrintln("create IID_IShellLink from CLSID_ShellLink")
	if err != nil {
		panic(err)
	}
	oIShellLink = GetIShellLinkFromIUnkown(u)
	errors.RegisterFinalCall(func() { oIShellLink.Release() })
	u, err = oIShellLink.QueryInterface(IID_IPersistFile)
	errors.OlePrintln("query IID_IPersistFile from IShellLink")
	if err != nil {
		panic(err)
	}
	oIPersistFile = GetIPersistFileFromIUnkown(u)
	errors.RegisterFinalCall(func() { oIPersistFile.Release() })
}
func GetLnkTarget(lnk string) string {
	defer errors.OleCatchErrors(false, "GetLnkTarget")
	err := oIPersistFile.Load(lnk, STGM_READ)
	if err != nil {
		errors.OlePrintln("oIPersistFile.Load error")
		panic(err)
	}
	t, err := oIShellLink.GetPath(SLGP_RAWPATH)
	if err != nil {
		errors.OlePrintln("oIShellLink.GetPath error")
		s := fmt.Sprintf("error from GetPath:%v\n", err)
		// find path in description
		des, err := oIShellLink.GetDescription()
		errors.OlePrintln("oIShellLink.GetDescription", des, err)
		des = strings.TrimSpace(des)

		if des != "" {
			//check slash
			des = strings.Replace(des, "/", "\\", -1)
			t, err = filepath.Abs(filepath.Join(filepath.Dir(lnk), des))
			if err != nil {
				s = fmt.Sprintf("%s error from GetDescription\n %v\n", s, err)
			} else {
				return t
			}
		}
		panic(s)

	}
	return t

}
