package main

import (
		"fmt"
			"path/filepath"
	"os"
	"strings"
)

type targetURL struct {
	host       string
	user       string
	repository string
}

func newTargetURL(target string) *targetURL {
	split := strings.Split(target, "/")
	var tu *targetURL
	if len(split) == 1 {
		tu = &targetURL{split[0], "", ""}
	} else if len(split) == 2 {
		tu = &targetURL{split[0], split[1], ""}
	} else if len(split) == 3 {
		tu = &targetURL{split[0], split[1], split[2]}
	}
	return tu
}

func (tu *targetURL) buildTargetPaths() []string {
	dir := filepath.Join(os.Getenv("GOPATH"), "src")
	hosts := tu.buildHosts(dir)

	var paths []string
	for _, host := range hosts { // ex. $GOPATH/src/github.com/
		users := tu.buildUsers(host)

		for _, user := range users { // ex. $GOPATH/src/github.com/kenju
			repos := tu.buildRepos(user)

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


func (tu *targetURL) buildHosts(dir string) []string {
	var hosts []string
	if len(tu.host) > 0 {
		hosts = []string{filepath.Join(dir, tu.host)}
	} else {
		hosts = flattenWalk(dir)
	}
	return hosts
}

func (tu *targetURL) buildUsers(host string) []string {
	var users []string
	if len(tu.user) > 0 {
		users = []string{filepath.Join(host, tu.user)}
	} else {
		users = flattenWalk(host)
	}
	return users
}

func (tu *targetURL) buildRepos(user string) []string {
	var repos []string
	if len(tu.repository) > 0 {
		repos = []string{filepath.Join(user, tu.repository)}
	} else {
		repos = flattenWalk(user)
	}
	return repos
}

func flattenWalk(path string) []string {
	matches, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		fmt.Printf("error file globbing: %+v\n", err)
	}
	return matches
}

