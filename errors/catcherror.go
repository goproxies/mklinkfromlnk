package errors

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

func AcceptEnterKeyExit(code int) {
	var s byte
	fmt.Println("Entry 'Enter' to exit...")
	fmt.Scanln(&s)
	os.Exit(code)

}

var BPrintError = false

func Printf(s string, ss ...interface{}) {
	if BPrintError {
		fmt.Printf(s, ss...)
	}
}
func Println(ss ...interface{}) {
	if BPrintError {
		fmt.Println(ss...)
	}
}
func CatchErrors(q bool, s ...interface{}) {
	err := recover()
	if err != nil {
		quitcount++
		PrintStack()
		h := ""
		he := ""
		con := GetGID()
		if len(s) > 0 {
			h = fmt.Sprintf("[%v]", s[0])
			he = fmt.Sprintf("%v", s[1:])
		}
		Printf("[co %d][%s] error:%s\n %s\n", con, h, err, he)
		if q {
			Printf("[co %d][%s] error count:%d\n", con, h, quitcount)
			_maincalled.Do(execFinalCall)
			AcceptEnterKeyExit(1)
		}

	}
}

var quitcount = 0

func GetErrorCount() int {
	return quitcount
}

var BOleDebug = false

func OlePrintf(s string, ss ...interface{}) {
	if BOleDebug {
		fmt.Printf(s, ss...)
	}
}
func OlePrintln(ss ...interface{}) {
	if BOleDebug {
		fmt.Println(ss...)
	}
}
func CoCatchErrors() {
	err := recover()
	if err != nil {
		Printf("[co %d]top error:%v\n", GetGID(), err)
		printStack()
	}
}
func OleCatchErrors(q bool, s ...interface{}) {
	defer CatchErrors(false, "OleCatchErrors")
	err := recover()
	if err != nil {
		quitcount++

		h := ""
		he := ""
		con := GetGID()
		if len(s) > 0 {
			h = fmt.Sprintf("[%v]", s[0])
			he = fmt.Sprintf("%v", s[1:])
		}
		if BOleDebug {
			pc, file, line, ok := runtime.Caller(2)
			if ok {
				pcName := runtime.FuncForPC(pc).Name()
				pcName = filepath.Base(pcName)
				file = filepath.Base(file)
				sk := fmt.Sprintf(" %s:%d %s\n", pcName, line, file)
				fmt.Printf("[co %d][%s] error:%v\n %s\n [panic at]%s\n", con, h, err, he, sk)
			}
			printStack()
		}
		if q {
			Printf("[co %d][%s] error count:%d\n", con, h, quitcount)
			AcceptEnterKeyExit(1)
		}

	}

}

var BPrintStack = false

func PrintStack(mm ...int) {
	if BPrintStack {
		printStack(mm...)
	}
}

func printStack(mm ...int) {
	m := 5
	if len(mm) > 0 {
		m = mm[0]
	}
	if m < 1 || m > 10 {
		m = 5
	}
	r := ""
	for i := 3; i < m+2; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if ok {
			pcName := runtime.FuncForPC(pc).Name()
			pcName = filepath.Base(pcName)
			r += fmt.Sprintf(" [stack]%s:%s:%d\n", pcName, file, line)
		}
	}
	fmt.Println(r)

}
func printStackAll() {

}

var mainfinalcalls []func()

func RegisterFinalCall(f func()) {
	mainfinalcalls = append(mainfinalcalls, f)
}

//exec final call once
var _maincalled = sync.Once{}

func MainCatchErrors(s ...interface{}) {
	defer AcceptEnterKeyExit(0)
	defer CatchErrors(false, "MainCatchErrors")
	defer _maincalled.Do(execFinalCall)

	err := recover()
	if err != nil {
		quitcount++
		PrintStack()
		h := ""
		he := ""
		con := GetGID()
		if len(s) > 0 {
			h = fmt.Sprintf("[%v]", s[0])
			he = fmt.Sprintf("%v", s[1:])
		}
		Printf("[co %d][%s] error:%s\n %s\n", con, h, err, he)
		Printf("[co %d][%s] error count:%d\n", con, h, quitcount)
	}

}
func execFinalCall() {
	for _, f := range mainfinalcalls {
		f()
	}
}
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
