package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// read flag options
	var scriptPath, target string
	// -script (-s)
	flag.StringVar(&scriptPath, "script", "script.sh", "script file to execute")
	flag.StringVar(&scriptPath, "s", "script.sh", "script file to execute")
	// -target (-t)
	flag.StringVar(&target, "target", "", "target URL to execute")
	flag.StringVar(&target, "t", "", "target URL to execute")
	flag.Parse()

	// get abstract script path
	path, err := filepath.Abs(scriptPath)
	if err != nil {
		fmt.Printf("error while getting absolute path for %s: %+v\n", path, err)
	}

	split := strings.Split(target, "/")
	var tu targetURL
	if len(split) == 1 {
		tu = targetURL{split[0], "", ""}
	} else if len(split) == 2 {
		tu = targetURL{split[0], split[1], ""}
	} else if len(split) == 3 {
		tu = targetURL{split[0], split[1], split[2]}
	}
	log.Printf("targetURL: %#v\n", tu)
	recursivelyRunGroomCommand(path, tu)
}

type targetURL struct {
	host       string
	user       string
	repository string
}

func recursivelyRunGroomCommand(scriptPath string, target targetURL) {
	dir := filepath.Join(os.Getenv("GOPATH"), "src")

	var hosts []string
	if len(target.host) > 0 {
		hosts = []string{filepath.Join(dir, target.host)}
	} else {
		hosts = flattenWalk(dir)
	}

	waitgroup := &sync.WaitGroup{}

	for _, host := range hosts { // ex. $GOPATH/src/github.com/
		var users []string
		if len(target.user) > 0 {
			users = []string{filepath.Join(host, target.user)}
		} else {
			users = flattenWalk(host)
		}

		for _, user := range users { // ex. $GOPATH/src/github.com/kenju
			var repos []string
			if len(target.repository) > 0 {
				repos = []string{filepath.Join(user, target.repository)}
			} else {
				repos = flattenWalk(user)
			}

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
