package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	recursivelyRunGroomCommand()
	log.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
}

// TODO: get script from STDIN and execute instead
func recursivelyRunGroomCommand() {
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
					go func(r string, wg *sync.WaitGroup) {
						defer wg.Done()
						runGroomCommand(r)
					}(repo, waitgroup)
				}
			}
		}
	}

	waitgroup.Wait()
}

func runCommand(name, dir string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s: %s", dir, fmt.Sprint(err))
	}
	return string(out)
}

func runGroomCommand(repo string) {
	// TODO: update log format
	fmt.Printf("[%s] %s", repo, runCommand("git", repo, "rev-parse", "--show-toplevel"))
	fmt.Printf("[%s] %s", repo, runCommand("git", repo, "rev-parse", "--abbrev-ref", "HEAD"))
	fmt.Printf("[%s] %s", repo, runCommand("git", repo, "checkout", "master"))
	fmt.Printf("[%s] %s", repo, runCommand("git", repo, "pull", "--prune"))
}

func flattenWalk(path string) []string {
	matches, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		fmt.Printf("error file globbing: %+v\n", err)
	}
	return matches
}
