package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
)

var logger *Logger
var numConcurrency int

func main() {
	// read flag options
	var scriptPath, target string
	var debug bool
	// -script (-s)
	flag.StringVar(&scriptPath, "script", "script.sh", "script file to execute")
	flag.StringVar(&scriptPath, "s", "script.sh", "script file to execute")
	// -target (-t)
	flag.StringVar(&target, "target", "", "target URL to execute")
	flag.StringVar(&target, "t", "", "target URL to execute")
	// -concurrency (-c)
	flag.IntVar(&numConcurrency, "concurrency", runtime.NumCPU(), "concurrency to execute")
	flag.IntVar(&numConcurrency, "c", runtime.NumCPU(), "concurrency to execute")
	// -debug
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()

	logger = NewLogger(debug)

	scriptPathAbs, err := filepath.Abs(scriptPath)
	if err != nil {
		fmt.Printf("error while getting absolute path for %s: %+v\n", scriptPathAbs, err)
	}

	tu := NewTargetURL(target)

	logger.Printf("**DEBUG mode = true**\n")
	logger.Printf("targetURL: %#v\n", *tu)

	paths := tu.buildTargetPaths()
	logger.Printf("Total paths count: %d\n", len(paths))

	runInAsync(scriptPathAbs, paths)
}

