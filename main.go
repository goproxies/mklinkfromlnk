package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/goproxies/mklinkfromlnk/dll"
	"github.com/goproxies/mklinkfromlnk/errors"
	"github.com/goproxies/mklinkfromlnk/lnkfile"
	. "github.com/goproxies/mklinkfromlnk/types"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	h          bool
	d          string
	skip       string
	r          bool
	s          bool
	v          bool
	vv         bool
	c          int
	e          bool
	oledebug   bool
	printerror bool
	printstack bool

	checksymlink bool
	skipre       *regexp.Regexp
	bskip        bool
	version      = "1.0.3"
	wait         sync.WaitGroup
	checksame    bool
)

func init() {
	flag.BoolVar(&h, "h", false, "print help")
	flag.StringVar(&d, "d", ".", "which directory to search")
	flag.StringVar(&skip, "skip", "", "skip directory pattern")
	flag.BoolVar(&r, "r", false, "create with relative path")
	flag.BoolVar(&s, "s", false, "show parsing dir")
	flag.BoolVar(&v, "v", false, "show infos verbosely")
	flag.BoolVar(&vv, "vv", false, "show infos more verbosely")
	flag.IntVar(&c, "c", 50, "minimum 1,maximum numbers of coroutines ")
	flag.BoolVar(&e, "e", false, "delete *.lnk")
	flag.BoolVar(&oledebug, "oledebug", false, "show ole errors")
	flag.BoolVar(&printerror, "printerror", false, "print errors")
	flag.BoolVar(&printstack, "printstack", false, "print stack when some errors occured")
	flag.Usage = usage
	flag.BoolVar(&checksymlink, "checksymlink", false, "convert '/' to '\\' in symlink ")
	flag.BoolVar(&checksame, "checksame", false, "correct target with the name of symlink")
}
func main() {
	//catch errors
	defer errors.MainCatchErrors("mklinkfromlnk")
	flag.Parse()
	if c < 1 {
		c = 1
	}
	if flag.NFlag() == 0 || h {
		flag.Usage()
		errors.AcceptEnterKeyExit(1)
	}
	if vv {
		v = true
		printstack = true
		oledebug = true
	}
	if v {
		s = true
		printerror = true
	}
	if printstack {
		printerror = true
	}
	if printerror {
		errors.BPrintError = true
	}
	if oledebug {
		printstack = true
		printerror = true
		errors.BOleDebug = true
		errors.BPrintStack = true
		errors.BPrintError = true
	}
	bskip = skip != ""
	skipre = regexp.MustCompile(skip)
	var paths = NewHashSet()
	ap, err := filepath.Abs(d)
	if err == nil {
		paths.Add(ap)
	} else {
		errors.Printf("error:get abs directory[%s]\n", d)
	}
	for _, g := range flag.Args() {
		ab, _ := filepath.Abs(g)
		fi, err := os.Lstat(ab)
		if err == nil {
			if !fi.IsDir() {
				VPrintf("%s is a file,skip", ab)
				continue
			}
			paths.Add(ab)
		} else {
			VPrintf("error:get abs directory[%s],skip\n", g)
		}
	}
	paths_o := paths.ToStringSlice()
	VPrintln("intial paths:", paths_o)
	if checksymlink {
		checksymlinkfrom(&paths_o)
	} else if checksame {
		checksamefrom(&paths_o)
	} else {
		mklinkfrom(&paths_o)
	}

	defer wait.Wait()
}
func checksymlinkfrom(dirs *[]string) {
	defer errors.CatchErrors(false, "checksymlinkfrom")
	ct := make(chan *[]SymlinkFileInfo, c)
	//deal with symlink files in the same directory first
	for _, d := range *dirs {
		lns := GetSymlinkFiles(d)
		if len(*lns) == 0 {
			continue
		}
		wait.Add(1)
		ct <- lns
		go func() {
			defer wait.Done()
			defer errors.CoCatchErrors()
			lns := <-ct
			VPrintln("SYMLINKFILES:", *lns)
			CheckSymlinkFiles(lns)
		}()

	}
	//then the subdirs
	for _, d := range *dirs {
		checksymlinkfrom(GetSubdirsSymlink(d))
	}
}
func checksamefrom(dirs *[]string) {
	defer errors.CatchErrors(false, "checksamefrom")
	ct := make(chan *[]SymlinkFileInfo, c)
	//deal with symlink files in the same directory first
	for _, d := range *dirs {
		lns := GetSymlinkFiles(d)
		if len(*lns) == 0 {
			continue
		}
		wait.Add(1)
		ct <- lns
		go func() {
			defer wait.Done()
			defer errors.CoCatchErrors()
			lns := <-ct
			VPrintln("SYMLINKFILES:", *lns)
			CheckSameSymlinkFiles(lns)
		}()

	}
	//then the subdirs
	for _, d := range *dirs {
		checksamefrom(GetSubdirsSymlink(d))
	}
}
func mklinkfrom(dirs *[]string) {
	defer errors.CatchErrors(false, "mklinkfrom")
	ct := make(chan *[]string, c)
	//deal with link files in the same directory first
	for _, d := range *dirs {
		lns := GetLnkFiles(d)
		if len(*lns) == 0 {
			continue
		}
		wait.Add(1)
		ct <- lns
		go func() {
			defer wait.Done()
			defer errors.CoCatchErrors()
			lns := <-ct
			VPrintln("LNKFILES:", *lns)
			CreateSymlinkFromLnk(lns)
		}()

	}
	//then the subdirs
	for _, d := range *dirs {
		mklinkfrom(GetSubdirs(d))
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `mklinkfromlnk version: %s
Usage: mklinkfromlnk [-h] [-e] [-r] [-s]  [-c maximum-numbers-of-coroutines] [-r use-relative-path]   
    [-printerror] [-printstack [-printerror]] [-v [-s -printerror]] [-oledebug [-printstack -printerror]] [-vv [-v -printstack -oledebug]]
    [-checksymlink]
    [-checksame]
    [-d directory] [directory...]
Options:
`, version)
	flag.PrintDefaults()
}

// makeCmdLine builds a command line out of args by escaping "special"
// characters and joining the arguments with spaces.
func makeCmdLine(args []string) string {
	var s string
	for _, v := range args {
		if s != "" {
			s += " "
		}
		s += EscapeArg(v)
	}
	return s
}

// EscapeArg rewrites command line argument s as prescribed
// in https://msdn.microsoft.com/en-us/library/ms880421.
// This function returns "" (2 double quotes) if s is empty.
// Alternatively, these transformations are done:
// - every back slash (\) is doubled, but only if immediately
//   followed by double quote (");
// - every double quote (") is escaped by back slash (\);
// - finally, s is wrapped with double quotes (arg -> "arg"),
//   but only if there is space or tab inside s.
func EscapeArg(s string) string {
	if len(s) == 0 {
		return "\"\""
	}
	n := len(s)
	hasSpace := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			n++
		case ' ', '\t':
			hasSpace = true
		}
	}
	if hasSpace {
		n += 2
	}
	if n == len(s) {
		return s
	}

	qs := make([]byte, n)
	j := 0
	if hasSpace {
		qs[j] = '"'
		j++
	}
	slashes := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		default:
			slashes = 0
			qs[j] = s[i]
		case '\\':
			slashes++
			qs[j] = s[i]
		case '"':
			for ; slashes > 0; slashes-- {
				qs[j] = '\\'
				j++
			}
			qs[j] = '\\'
			j++
			qs[j] = s[i]
		}
		j++
	}
	if hasSpace {
		for ; slashes > 0; slashes-- {
			qs[j] = '\\'
			j++
		}
		qs[j] = '"'
		j++
	}
	return string(qs[:j])
}

//execute windows command
func ShellOut(command ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmdswitch := []string{"/C"}
	cmd := exec.Command("CMD", append(cmdswitch, command...)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stds := dll.CovConsoleBytesToString(stdout.Bytes())
	if v {
		con := errors.GetGID()
		VPrintln(fmt.Sprintf(" [co %d][ShellOut] execute cmd:%s %s\n", con, cmd.Path, makeCmdLine(cmd.Args)),
			fmt.Sprintf("[co %d][ShellOut] stdout:%s\n", con, stds),
			fmt.Sprintf("[co %d][ShellOut] error:%v\n", con, err),
			fmt.Sprintf("[co %d][ShellOut] stderr:%s\n", con, dll.CovConsoleBytesToString(stderr.Bytes())))
	}
	return stds, err
}
func SShellOut(command ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmdswitch := []string{"/C"}
	cmd := exec.Command("CMD", append(cmdswitch, command...)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stds := dll.CovConsoleBytesToString(stdout.Bytes())
	if s {
		con := errors.GetGID()
		SPrintln(fmt.Sprintf(" [co %d][ShellOut] execute cmd:%s %s\n", con, cmd.Path, makeCmdLine(cmd.Args)),
			fmt.Sprintf("[co %d][ShellOut] stdout:%s\n", con, stds),
			fmt.Sprintf("[co %d][ShellOut] error:%v\n", con, err),
			fmt.Sprintf("[co %d][ShellOut] stderr:%s\n", con, dll.CovConsoleBytesToString(stderr.Bytes())))
	}
	return stds, err
}

//-------------------------------------------------------------
//create symlink

// default sep `\r\n`
// split with every sep
func OutToSlice(o, prefix string, sep ...string) *[]string {
	var seps string
	n := len(sep)
	if n == 0 {
		seps = `[\r\n]+`
	} else {
		ts := make([]string, n)
		for i, v := range sep {
			ts[i] = fmt.Sprintf("(%s)+", v)
		}
		seps = strings.Join(ts, "|")
	}
	o = strings.TrimSpace(o)
	var strs []string
	re := regexp.MustCompile(seps)
	strs = re.Split(o, -1)

	var ts []string
	for _, v := range strs {
		v = strings.TrimSpace(v)
		if v != "" {
			ts = append(ts, v)
		}
	}
	if prefix != "" {
		for i, v := range ts {
			ts[i] = prefix + v
		}
	}
	s := NewHashSet()
	s.AddStringSlice(ts)
	r := s.ToStringSlice()
	return &r
}
func GetSubdirs(parent string) *[]string {
	dirs := make([]string, 0)
	if !strings.HasSuffix(parent, `\`) {
		parent = parent + `\`
	}
	t, err := ShellOut("DIR", parent, "/A:D", "/B")
	if err == nil {
		dirs = *OutToSlice(t, parent)
	}
	VPrintln("SUBDIRS:", dirs)
	if bskip {
		s := NewHashSet()
		for _, d := range dirs {
			if !skipre.MatchString(d) {
				VPrintf("match dir:%s\n", d)
				s.Add(d)
			} else {
				VPrintf("skip dir:%s\n", d)
			}
		}
		dirs = s.ToStringSlice()
	}
	return &dirs
}
func trimLnk(s []string) []string {
	t := make([]string, 0)
	for _, v := range s {
		if strings.HasSuffix(v, ".lnk") {
			t = append(t, v)
		}
	}
	return t
}
func GetLnkFiles(dir string) *[]string {
	SPrintln("PARSING:", dir)
	files := make([]string, 0)
	if !strings.HasSuffix(dir, `\`) {
		dir = dir + `\`
	}

	t, _ := ShellOut("DIR", dir+"*.lnk", "/A:-D", "/B")
	// in subdir ShellOut return err!=nil
	files = *OutToSlice(t, dir)
	files = trimLnk(files)

	return &files
}
func CreateSymlinkFromLnk(lns *[]string) {
	defer errors.CatchErrors(false, "CreateSymlinkFromLnk")
	if len((*lns)) == 0 {
		return
	}
	dir := filepath.Dir((*lns)[0])
	for _, ln := range *lns {
		isfile, target := GetLnkInfo(ln)
		if target != "" {
			VPrintln(ln, "->", target)
			MakeSymlink(dir, ln, isfile, target)
		}
	}
}

func GetLnkInfo(lnk string) (isfile bool, target string) {
	defer errors.CatchErrors(false, "GetLnkInfo")
	target = lnkfile.GetLnkTarget(lnk)
	if target != "" {
		ls, err := os.Lstat(target)
		if err != nil {
			isfile = true
		} else {
			isfile = !ls.IsDir()
		}
	}
	return
}

func MakeSymlink(dir, src string, isfile bool, target string) {
	if r {
		target, _ = filepath.Rel(dir, target)
	}

	//trim `.lnk`
	ru := []rune(src)
	src_path := string(ru[:len(ru)-4])
	fmt.Printf("%s\n ->%s\n", src_path, target)
	//delete non-symlink file/dir first
	if PathExists(src_path) {
		if !isSymlink(src_path) {
			DeleteAllFiles(src_path, isfile)
			if !isfile {
				os.Remove(src_path)
			}
			SPrintf("Delete:%s\n", src_path)
		}
	}
	if isfile {
		ShellOut("mklink", src_path, target)
	} else {
		ShellOut("mklink", "/d", src_path, target)
	}

	if e {
		//delete *.lnk only symlink has created
		if PathExists(src_path) {
			os.Remove(src)
		}
	}
}
func isSymlink(d string) bool {
	l, _ := os.Lstat(d)

	return l.Mode()&os.ModeSymlink != 0

	//return false
}
func DeleteAllFiles(d string, isfile bool) {
	SShellOut("del", "/F", "/S", "/Q", d)
}

//-------------------------------------------------------------
//-checksymlink
func GetSubdirsSymlink(parent string) *[]string {
	dirs := make([]string, 0)
	if !strings.HasSuffix(parent, `\`) {
		parent = parent + `\`
	}
	t, err := ShellOut("DIR", parent, "/A:-LD", "/B")
	if err == nil {
		dirs = *OutToSlice(t, parent)
	}
	VPrintln("SUBDIRS:", dirs)
	if bskip {
		s := NewHashSet()
		for _, d := range dirs {
			if !skipre.MatchString(d) {
				VPrintf("match dir:%s\n", d)
				s.Add(d)
			} else {
				VPrintf("skip dir:%s\n", d)
			}
		}
		dirs = s.ToStringSlice()
	}
	return &dirs
}
func CheckSymlinkFiles(lns *[]SymlinkFileInfo) {
	defer errors.CatchErrors(false, "GetSymlinkFiles")
	if len((*lns)) == 0 {
		return
	}
	for _, ln := range *lns {
		target := ln.Target
		if target != "" {
			if strings.Contains(target, "/") {
				VPrintln("TRANSFORM:", "->", target)
				target = strings.Replace(target, "/", "\\", -1)
				CheckSymlink(ln.Src, ln.Isfile, target)
			}
		}
	}

}
func CheckSymlink(src string, isfile bool, target string) {
	fmt.Printf("%s\n ->%s\n", src, target)
	os.Remove(src)
	if isfile {
		ShellOut("mklink", src, target)
	} else {
		ShellOut("mklink", "/d", src, target)
	}
}

type SymlinkFileInfo struct {
	Isfile bool
	Src    string
	Target string
}

var symlink_line_re = regexp.MustCompile(`(?m:^.*<SYMLINK(?P<isDir>D?)>\s+\b(?P<src>.*)\b\s+\[(?P<target>[^\]]+)\].*$)`)

func OutToSymlinkFileInfo(t, dir string) *[]SymlinkFileInfo {
	defer errors.CatchErrors(false, "OutToSymlinkFileInfo")
	sl := OutToSlice(t, "")
	r := make([]SymlinkFileInfo, 0)
	for _, k := range *sl {
		if strings.Contains(k, "[") {
			m := symlink_line_re.FindStringSubmatch(k)
			r = append(r, SymlinkFileInfo{m[1] != "D", dir + m[2], m[3]})
		}
	}
	return &r
}
func GetSymlinkFiles(dir string) *[]SymlinkFileInfo {
	SPrintln("PARSING:", dir)
	files := make([]SymlinkFileInfo, 0)
	if !strings.HasSuffix(dir, `\`) {
		dir = dir + `\`
	}

	t, _ := ShellOut("DIR", dir, "/A:L")
	// in subdir ShellOut return err!=nil
	files = *OutToSymlinkFileInfo(t, dir)
	return &files
}

//-------------------------------------------------------------
//-checksame

func CheckSameSymlinkFiles(lns *[]SymlinkFileInfo) {
	defer errors.CatchErrors(false, "CheckSameSymlinkFiles")
	if len((*lns)) == 0 {
		return
	}
	for _, ln := range *lns {
		target := ln.Target
		base_src := filepath.Base(ln.Src)
		base_target := filepath.Base(target)
		if base_src != base_target {
			VPrintln("CHECKSAME:", "->", target)
			target = target[:len(target)-len(base_target)] + base_src
			CheckSameSymlink(ln.Src, ln.Isfile, target)
		}
	}

}

func CheckSameSymlink(src string, isfile bool, target string) {
	fmt.Printf("%s\n ->%s\n", src, target)
	os.Remove(src)
	if isfile {
		ShellOut("mklink", src, target)
	} else {
		ShellOut("mklink", "/d", src, target)
	}
}

//-------------------------------------------------------------
//assist functions

func PathExists(path string) bool {
	_, err := os.Lstat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func VPrintf(s string, ss ...interface{}) {
	if v {
		fmt.Printf(s, ss...)
	}
}
func VPrintln(s ...interface{}) {
	if v {
		fmt.Println(s...)
	}
}
func SPrintf(a string, aa ...interface{}) {
	if s {
		fmt.Printf(a, aa...)
	}
}
func SPrintln(a ...interface{}) {
	if s {
		fmt.Println(a...)
	}
}
