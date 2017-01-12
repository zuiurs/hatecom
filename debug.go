package main

import (
	"fmt"
	"log"
	"os"
)

const (
	//DebugPrefix = "debug: "
	DebugPrefix = ""
)

type DebugLogger struct {
	*log.Logger
	Mode bool // debug mode enable or disable
}

var debug = DebugLogger{
	Logger: log.New(os.Stdout, DebugPrefix, log.Lshortfile),
	Mode:   true,
}

func (d *DebugLogger) Printf(format string, v ...interface{}) {
	if d.Mode {
		d.Output(2, fmt.Sprint(v...))
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
