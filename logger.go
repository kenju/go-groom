package main

import (
	"fmt"
	"time"
)

type Logger struct {
	isDebug     bool
	startedTime time.Time
}

func NewLogger(isDebug bool) (*Logger) {
	return &Logger{
		isDebug:     isDebug,
		startedTime: time.Now(),
	}
}

func (l *Logger) Printf(format string, a ...interface{}) {
	if l.isDebug {
		fmt.Printf(format, a...)
	}
}

func (l *Logger) startTimer() {
	l.startedTime = time.Now()
}

func (l *Logger) endTimer() {
	l.Printf("Execution took: %v\n", time.Since(l.startedTime))
}
