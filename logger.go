package main

import (
	"fmt"
	"time"
				)

// Logger handles logging and tracking with timer
type Logger struct {
	isDebug     bool
	startedTime time.Time
}

// NewLogger returns a new instance of Logger pointer
func NewLogger(isDebug bool) (*Logger) {
	return &Logger{
		isDebug:      isDebug,
		startedTime:  time.Now(),
	}
}

// Printf is a delegation to fmt.Printf with application-specific context
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
