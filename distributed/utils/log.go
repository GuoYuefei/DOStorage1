package utils

import (
	"io"
	"log"
	"os"
)

const (
	None = 1 << iota
	Exception
	Err
	Warning
	Record
	Info
	Debug
)

var Log ILog
var pre = map[int]string{
	Debug:     "[debug] ",
	Record:    "[record] ",
	Info:      "[info] ",
	Warning:   "[warning] ",
	Err:       "[error] ",
	Exception: "[exception] ",
}

func init() {
	Log = NewLog()
}

type ILog interface {
	Printf(p int, format string, a ...interface{})
	Println(p int, a ...interface{})
	Print(p int, a ...interface{})
	SetPriority(pri int)
	GetPriority() int
	SetOut(writer io.Writer)
}

type Logger struct {
	log      *log.Logger
	priority int
}

func NewLog() *Logger {
	l := &Logger{}
	l.log = log.New(os.Stdout, "", log.LstdFlags|log.Lmsgprefix)
	l.priority = Info
	return l
}

func (l *Logger) Printf(p int, format string, a ...interface{}) {
	if p == Exception {
		panic(a)
	}
	if l.priority < p {
		// 优先级过滤， 例 p=Info， l.p=Warning
		return
	}
	l.log.SetPrefix(prefix(p))
	l.log.Printf(format, a...)
}

func (l *Logger) Println(p int, a ...interface{}) {
	if p == Exception {
		panic(a)
	}
	if l.priority < p {
		return
	}
	l.log.SetPrefix(prefix(p))
	l.log.Println(a...)
}

func (l *Logger) Print(p int, a ...interface{}) {
	if p == Exception {
		panic(a)
	}
	if l.priority < p {
		return
	}
	l.log.SetPrefix(prefix(p))
	l.log.Print(a...)
}

func (l *Logger) SetPriority(pri int) {
	l.priority = pri
}

func (l *Logger) GetPriority() int {
	return l.priority
}

func (l *Logger) SetOut(writer io.Writer) {
	l.log.SetOutput(writer)
}

// 根据优先级设置log的前缀
func prefix(p int) string {
	if v, ok := pre[p]; ok {
		return v
	}
	return "none"
}

func FailOnError(err error, message string) {
	if err != nil {
		Log.Println(Err, message)
	}
}

func PanicOnError(err error, message string) {
	if err != nil {
		Log.Println(Exception, message)
	}
}
