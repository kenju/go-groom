package main

import (
	"fmt"
	"os"
	"path/filepath"
	"os/exec"
)

func runGroomCommand(repo string) {
	// change dir
	fmt.Printf("[INFO] chdir %s\n", repo)
	err := os.Chdir(repo)
	if err != nil {
		fmt.Printf("error while changeing directory to %s\n", repo)
	}
	// `git checkout master`
	out, err := exec.Command("git", "checkout", "master").Output()
	if err != nil {
		fmt.Printf("error while `git checkout master` at %s: %+v\n", repo, err)
	}
	fmt.Println(string(out))
	// `git pull --prune`
	out, err = exec.Command("git", "pull", "--prune").Output()
	if err != nil {
		fmt.Printf("error while `git pull --prune` at %s: %+v\n", repo, err)
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

func main() {
	dir := filepath.Join(os.Getenv("GOPATH"), "src")
	matches := flattenWalk(dir) // ex. $GOPATH/src

	// TODO: change this O(N^3) for loop with goroutine
	for _, host := range matches {
		users := flattenWalk(host) // ex. $GOPATH/src/github.com/

		for _, user := range users {
			repos := flattenWalk(user) // ex. $GOPATH/src/github.com/kenju

			// ex. $GOPATH/src/github.com/kenju/go-groom
			for _, repo := range repos {
				fi, err := os.Stat(repo)
				if err != nil {
					fmt.Printf("error while getting os stat for %+v\n", fi)
				}
				if fi.IsDir() {
					runGroomCommand(repo)
				} else {
					fmt.Printf("%+v is not a dir. skipping.\n", fi.Name())
				}
			}
		}
	}
}
