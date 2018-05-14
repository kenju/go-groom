package main

import (
  "fmt"
  "os"
  "path/filepath"
)

func main() {
  fmt.Println("Hello, world!")

  // TODO: replace with $GO_PATH
  dir := "/Users/kenjuwagatsuma/.ghq/src/*"

  matches, err := filepath.Glob(dir)
  if err != nil {
    fmt.Printf("error file globbing: %+v\n", err)
  }

  // TODO: change this O(N^3) for loop with goroutine
  for _, host := range matches {
    users, err := filepath.Glob(host+"/*")
    if err != nil {
      fmt.Printf("error file globbing: %+v\n", err)
    }

    for _, user := range users {
      repos, err := filepath.Glob(user+"/*")
      if err != nil {
        fmt.Printf("error file globbing: %+v\n", err)
      }

      for _, repo := range repos {
        fi, err := os.Stat(repo)
        if err != nil {
          fmt.Printf("error while getting os stat for %+v\n", fi)
        }
        if fi.IsDir() {
          // TODO: cw and run git pull
          fmt.Printf("%+v is a repository.\n", repo)
        } else {
          fmt.Printf("%+v is not a dir. skipping.\n", fi.Name())
        }
      }
    }
  }
}
