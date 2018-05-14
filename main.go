package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
)

func main() {
  // ex. $GOPATH/src
  dir := filepath.Join(os.Getenv("GOPATH"), "src", "*")
  matches, err := filepath.Glob(dir)
  if err != nil {
    fmt.Printf("error file globbing: %+v\n", err)
  }

  // TODO: change this O(N^3) for loop with goroutine
  // ex. $GOPATH/src/github.com/
  for _, host := range matches {
    users, err := filepath.Glob(filepath.Join(host, "*"))
    if err != nil {
      fmt.Printf("error file globbing: %+v\n", err)
    }

    // ex. $GOPATH/src/github.com/kenju
    for _, user := range users {
      repos, err := filepath.Glob(filepath.Join(user, "*"))
      if err != nil {
        fmt.Printf("error file globbing: %+v\n", err)
      }

      // ex. $GOPATH/src/github.com/kenju/go-groom
      for _, repo := range repos {
        fi, err := os.Stat(repo)
        if err != nil {
          fmt.Printf("error while getting os stat for %+v\n", fi)
        }
        if fi.IsDir() {
          // change dir
          fmt.Printf("[INFO] chdir %s\n", repo)
          err := os.Chdir(repo)
          if err != nil {
            fmt.Printf("error while changeing directory to %s\n", repo)
          }
          // run command
          out, err := exec.Command("git", "checkout", "master").Output()
          if err != nil {
            fmt.Printf("error while executing command at %s\n", repo)
          }
          fmt.Println(string(out))
          // out, err := exec.Command("git", "pull", "--prune")
        } else {
          fmt.Printf("%+v is not a dir. skipping.\n", fi.Name())
        }
      }
    }
  }
}
