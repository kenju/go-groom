package main

import (
	"os"
	"log"
	"runtime/trace"
	"context"
)

const (
	traceFile = "trace.out"
)

// Tracer track application performance with runtime/trace package.
// Call `(t *Tracer) startTraceTask()` to start tracing.
// Be sure to call `(t *Tracer) endTraceTask()` to finish tracing by releasing resources (e.g. file descriptor)
type Tracer struct {
	logger *Logger
	// fields for tracing
	traceCtx     context.Context
	traceTaskMap map[string]trace.Task
	traceFile    *os.File
}

// NewTracer returns a new instance of Tracer pointer
func NewTracer(logger *Logger) (*Tracer) {
	return &Tracer{
		logger:       logger,
		traceTaskMap: make(map[string]trace.Task),
	}
}

// trace runtime performance by runtime/trace package
//
// NOTE: you MUST call (l *Logger)endTraceTask to finish task
func (t *Tracer) startTraceTask(taskName string) {
	f, err := os.Create(traceFile)
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	t.traceFile = f
	t.logger.Printf("Created %s\n", traceFile)

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	t.logger.Printf("Started tracing '%s' task\n", taskName)

	ctx := context.Background()
	ctx, task := trace.NewTask(ctx, taskName)

	t.traceCtx = ctx
	t.traceTaskMap[taskName] = *task

	t.logger.Printf("[WARNING] Do call (l *Logger) endTraceTask() to finish resources!!")
}

func (t *Tracer) traceLog(category, message string) {
	trace.Log(t.traceCtx, category, message)
}

func (t *Tracer) traceRegion(regionType string, fn func()) {
	trace.WithRegion(t.traceCtx, regionType, fn)
}

func (t *Tracer) endTraceTask(taskName string) {
	task, ok := t.traceTaskMap[taskName]
	if !ok {
		log.Fatalf("not found task for %s", taskName)
	}

	t.logger.Printf("Finished tracing '%s' task\n", taskName)

	// finish resources
	task.End()
	trace.Stop()
	if err := t.traceFile.Close(); err != nil {
		log.Fatalf("failed to close trace output file: %v", err)
	}
}
