package main

import (
	"fmt"
	"time"
	"os"
	"log"
	"runtime/trace"
	"context"
)

// Logger handles logging and tracking with timer
type Logger struct {
	isDebug     bool
	startedTime time.Time
	// fields for tracing
	traceCtx     context.Context
	traceTaskMap map[string]trace.Task
	traceFile    *os.File
}

// NewLogger returns a new instance of Logger pointer
func NewLogger(isDebug bool) (*Logger) {
	return &Logger{
		isDebug:      isDebug,
		startedTime:  time.Now(),
		traceTaskMap: make(map[string]trace.Task),
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

// trace runtime performance by runtime/trace package
//
// NOTE: you MUST call (l *Logger)endTraceTask to finish task
func (l *Logger) startTraceTask(taskName string) {
	fileName := "trace.out"
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	l.traceFile = f
	l.Printf("Created %s\n", fileName)

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	l.Printf("Started tracing '%s' task\n", taskName)

	ctx := context.Background()
	ctx, task := trace.NewTask(ctx, taskName)

	l.traceCtx = ctx
	l.traceTaskMap[taskName] = *task

	l.Printf("[WARNING] Do call (l *Logger) endTraceTask() to finish resources!!")
}

func (l *Logger) traceLog(category, message string) {
	trace.Log(l.traceCtx, category, message)
}

func (l *Logger) traceRegion(regionType string, fn func()) {
	trace.WithRegion(l.traceCtx, regionType, fn)
}

func (l *Logger) endTraceTask(taskName string) {
	task, ok := l.traceTaskMap[taskName]
	if !ok {
		log.Fatalf("not found task for %s", taskName)
	}

	l.Printf("Finished tracing '%s' task\n", taskName)

	// finish resources
	task.End()
	trace.Stop()
	if err := l.traceFile.Close(); err != nil {
		log.Fatalf("failed to close trace output file: %v", err)
	}
}
