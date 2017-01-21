package log

import (
	"fmt"
	"io"
	"log"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type DebugLogger struct {
	*log.Logger
	Mode bool // debug mode enable or disable
}

func New(out io.Writer, prefix string, flag int, mode bool) *DebugLogger {
	return &DebugLogger{
		Logger: log.New(out, prefix, flag),
		Mode:   mode,
	}
}

func (d *DebugLogger) Printf(format string, v ...interface{}) {
	if d.Mode {
		d.Output(2, fmt.Sprintf(format, v...))
	}
}

func (d *DebugLogger) Println(v ...interface{}) {
	if d.Mode {
		d.Output(2, fmt.Sprintln(v...))
	}
}

// (l *Logger) Output() ではなく、(l *Logger) Printf() などを使用した場合、
// コールスタックが1段深くなるので debug.go の情報が出力される。
// そのため Output() を直接呼ぶ方法を使用する。
