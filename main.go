package main

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
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

	path, err := filepath.Abs(scriptPath)
	if err != nil {
		fmt.Printf("error while getting absolute path for %s: %+v\n", path, err)
	}

	tu := NewTargetURL(target)

	logger.Printf("**DEBUG mode = true**\n")
	logger.Printf("targetURL: %#v\n", *tu)

	paths := tu.buildTargetPaths()
	logger.Printf("Total paths count: %d\n", len(paths))

	runInAsync(path, paths)
}

type execResult struct {
	Error error
	Out   string
}

func runInAsync(scriptPath string, paths []string) {
	// send a signal to cancel goroutines which are internally invoked inside functions
	done := make(chan interface{})
	defer close(done)

	logger.startTimer()

	// spin up the number of pipelines to the number of available CPU on the machine
	logger.Printf("Spinning up %d pipeline\n", numConcurrency)
	pipelines := make([]<-chan interface{}, numConcurrency)
	targetPathCh := stringArrToCh(done, paths)
	for i := 0; i < numConcurrency; i++ {
		pipelines[i] = commandExecutor(done, targetPathCh, scriptPath)
	}

	// execute commands concurrently in each pipelines
	for result := range take(done, fanIn(done, pipelines...), len(paths)) {
		if result.(execResult).Error != nil {
			fmt.Printf("Error: %v\n", result.(execResult).Error)
		}
		fmt.Println(result.(execResult).Out)
	}

	logger.endTimer()
}

// stage to take values from channels
func take(
	done <-chan interface{},
	valueCh <-chan interface{},
	num int,
) <-chan interface{} {
	takeCh := make(chan interface{})

	go func() {
		defer close(takeCh)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeCh <- <-valueCh:
			}
		}
	}()

	return takeCh
}

// stage to multiplex multiple channels
func fanIn(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedCh := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedCh <- i:
			}
		}
	}

	// select from all the channels
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedCh)
	}()

	return multiplexedCh
}

// stage for converting String array to channel
func stringArrToCh(
	done <-chan interface{},
	arr []string,
) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		for _, v := range arr {
			select {
			case <-done:
				return
			case ch <- v:
			}
		}
	}()

	return ch
}

// stage for executing command at target dir
func commandExecutor(
	done <-chan interface{},
	stringCh <-chan string,
	scriptPath string,
) <-chan interface{} {
	resultCh := make(chan interface{})

	go func() {
		defer close(resultCh)

		for {
			select {
			case <-done:
				return
			case resultCh <- execCommand(<-stringCh, scriptPath):
			}
		}
	}()

	return resultCh
}

// execute command
func execCommand(dir, scriptPath string) execResult {
	cmd := exec.Command(scriptPath)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return execResult{Error: err, Out: ""}
	}
	return execResult{Error: nil, Out: string(out)}
}
