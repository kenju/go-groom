# go-groom

[![CircleCI](https://circleci.com/gh/kenju/go-groom.svg?style=svg)](https://circleci.com/gh/kenju/go-groom)

go-groom run grooming commands against multiple repositories concurrently.

# Install

```
$ go get -u github.com/kenju/go-groom
```

# Usage

## Options
### `-script (-s)`

Path any executable script to run in the each repositories.

```sh
$ cat script.sh
#!/usr/bin/env sh
git checkout master
git pull --prune

$ go-groom -script script.sh
```

### `-target (-t)`

Specify target repository.


```sh
# all repository under "github.com/**/*"
$ go-groom -script script.sh -host github.com

# all repository under "github.com/golang/*"
$ go-groom -script script.sh -organization github.com/golang

# single repository "github.com/golang/go"
$ go-groom -script script.sh -repo github.com/golang/go
```
