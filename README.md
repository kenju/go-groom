# go-groom

[![CircleCI](https://circleci.com/gh/kenju/go-groom.svg?style=svg)](https://circleci.com/gh/kenju/go-groom)

go-groom run grooming commands against multiple repositories concurrently.

# Install

```sh
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
$ go-groom -script script.sh -target github.com

# all repository under "github.com/golang/*"
$ go-groom -script script.sh -target github.com/golang

# single repository "github.com/golang/go"
$ go-groom -script script.sh -target github.com/golang/go
```

### `-concurrency (-c)`

Specify the number of concurrency to execute.

```sh
# spin up 8 pipeline
$ go-groom \
    -script script.sh \
    -target github.com \
    -concurrency 8
```

# Milestones

- Support timeout (a.k.a. deadlines) for each execution
- Integrate profiling with pprof/trace into CI
- use go generate for pipeline stage building with typed
