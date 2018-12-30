package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func buildHosts(target targetURL, dir string) []string {
	var hosts []string
	if len(target.host) > 0 {
		hosts = []string{filepath.Join(dir, target.host)}
	} else {
		hosts = flattenWalk(dir)
	}
	return hosts
}

func buildUsers(target targetURL, host string) []string {
	var users []string
	if len(target.user) > 0 {
		users = []string{filepath.Join(host, target.user)}
	} else {
		users = flattenWalk(host)
	}
	return users
}

func buildRepos(target targetURL, user string) []string {
	var repos []string
	if len(target.repository) > 0 {
		repos = []string{filepath.Join(user, target.repository)}
	} else {
		repos = flattenWalk(user)
	}
	return repos
}

func buildTargetPaths(target targetURL) []string {
	dir := filepath.Join(os.Getenv("GOPATH"), "src")
	hosts := buildHosts(target, dir)

	var paths []string
	for _, host := range hosts { // ex. $GOPATH/src/github.com/
		users := buildUsers(target, host)

		for _, user := range users { // ex. $GOPATH/src/github.com/kenju
			repos := buildRepos(target, user)

			for _, repo := range repos { // ex. $GOPATH/src/github.com/kenju/go-groom
				fi, err := os.Stat(repo)
				if err != nil {
					fmt.Printf("error while getting os stat for %+v\n", fi)
				}
				if fi.IsDir() {
					paths = append(paths, repo)
				}
			}
		}
	}

	return paths
}

func recursivelyRunGroomCommand(scriptPath string, target targetURL) {
	done := make(chan interface{})
	defer close(done)

	paths := buildTargetPaths(target)

	fmt.Printf("Total paths count: %d\n", len(paths))

	for result := range runInAsync(done, scriptPath, paths...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("out: %s", result.Out)
	}
}

type Result struct {
	Error error
	Out   string
}

func runInAsync(done <-chan interface{}, scriptPath string, paths ...string) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		for _, path := range paths {
			out, err := runGroomCommand(path, scriptPath)
			result := Result{Error: err, Out: out}

			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()

	return results
}

func runGroomCommand(dir, scriptPath string) (string, error) {
	cmd := exec.Command(scriptPath)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func flattenWalk(path string) []string {
	matches, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		fmt.Printf("error file globbing: %+v\n", err)
	}
	return matches
}
