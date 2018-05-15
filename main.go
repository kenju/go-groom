package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"flag"
	"log"
	"runtime"
)

func main() {
	// TODO: filter target host/org/user for running command
	// read flag options
	var scriptPath string
	flag.StringVar(&scriptPath, "script", "script.sh", "script file to execute")
	flag.StringVar(&scriptPath, "s", "script.sh", "script file to execute")
	flag.Parse()

	// get abstract script path
	path, err := filepath.Abs(scriptPath);
	if err != nil {
		fmt.Printf("error while getting absolute path for %s: %+v\n", path, err)
	}

	// core
	recursivelyRunGroomCommand(path)

	logDebug()
}

func logDebug() {
	log.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
}

func recursivelyRunGroomCommand(scriptPath string) {
	dir := filepath.Join(os.Getenv("GOPATH"), "src")
	matches := flattenWalk(dir) // ex. $GOPATH/src

	waitgroup := &sync.WaitGroup{}

	for _, host := range matches { // ex. $GOPATH/src/github.com/
		users := flattenWalk(host)

		for _, user := range users { // ex. $GOPATH/src/github.com/kenju
			repos := flattenWalk(user)

			for _, repo := range repos { // ex. $GOPATH/src/github.com/kenju/go-groom
				fi, err := os.Stat(repo)
				if err != nil {
					fmt.Printf("error while getting os stat for %+v\n", fi)
				}
				if fi.IsDir() {
					waitgroup.Add(1)
					// TODO: visualize go script process
					go func(r, p string, wg *sync.WaitGroup) {
						defer wg.Done()
						runGroomCommand(r, p)
					}(repo, scriptPath, waitgroup)
				}
			}
		}
	}

	waitgroup.Wait()
}

func runGroomCommand(dir, scriptPath string) {
	cmd := exec.Command(scriptPath)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s: %s", dir, fmt.Sprint(err))
	}
	fmt.Println(string(out))
}

func flattenWalk(path string) []string {
	matches, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		fmt.Printf("error file globbing: %+v\n", err)
	}
	return matches
}
