# go-groom

go-groom run grooming commands against multiple repositories concurrently.

# Install

```
$ go get github.com/kenju/go-groom
```

# Usage

## Options
### [required] `-script (-s)`

Path any executable script to run in the each repositories.

```sh
$ cat script.sh
#!/usr/bin/env sh
git checkout master
git pull --prune
$ go-groom -script script.sh
```

## [optional] `-host (-h)`, `-organization (-o)`, `-repo (-r)`

WIP

```sh
$ go-groom -script script.sh -host github.com
$ go-groom -script script.sh -organization github.com/golang
$ go-groom -script script.sh -repo github.com/golang/go
```
