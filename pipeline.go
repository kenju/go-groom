package main

import (
	"os/exec"
	"sync"
)

type execResult struct {
	Dir        string
	ScriptPath string
	Error      error
	Out        string
}

func runInAsync(scriptPath string, paths []string, logger *Logger) {
	// send a signal to cancel goroutines which are internally invoked inside functions
	done := make(chan interface{})
	defer close(done)

	logger.startTimer()

	// spin up the number of pipelines to the number of available CPU on the machine
	logger.Printf("Spinning up %d pipeline\n", numConcurrency)
	executionPipeline := make([]<-chan interface{}, numConcurrency)
	targetPathCh := stringArrToCh(done, paths)
	for i := 0; i < numConcurrency; i++ {
		executionPipeline[i] = commandExecutor(done, targetPathCh, scriptPath)
	}

	var numError int
	// execute commands concurrently in each pipelines
	pipelines := take(done, fanIn(done, executionPipeline...), len(paths))
	for result := range pipelines {
		logger.Printf("\t" + result.(execResult).Dir + "\n")
		if result.(execResult).Error != nil {
			numError++
			logger.Printf("\t\tError: %v\n", result.(execResult).Error)
		}
		logger.Printf("\t\t" + result.(execResult).Out + "\n")
	}

	logger.endTimer()
	logger.Printf("%d paths, %d error\n", len(paths), numError)
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
		return execResult{Dir: dir, ScriptPath: scriptPath, Error: err, Out: ""}
	}
	return execResult{Dir: dir, ScriptPath: scriptPath, Error: nil, Out: string(out)}
}
