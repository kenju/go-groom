package main

import "fmt"

type Logger struct {
	isDebug bool
}

func NewLogger(isDebug bool) (*Logger) {
	return &Logger{isDebug: isDebug}
}

func (l *Logger) Printf(format string, a ...interface{}) {
	if l.isDebug {
		fmt.Printf(format, a...)
	}
}
