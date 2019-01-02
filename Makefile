NAME := groom
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
	-X 'main.revision=$(REVISION)'

OK_COLOR    = \033[0;32m
ERROR_COLOR = \033[0;31m
WARN_COLOR  = \033[0;33m
NO_COLOR    = \033[m

OK_STRING    = "[OK]"
ERROR_STRING = "[ERROR]"
WARN_STRING  = "[WARNING]"

## Build binaries and run in debug mode
run-debug: build
	go run -race *.go \
		-script script.sh \
		-target github.com/kenju \
		-concurrency 8 \
		-debug

## Build binaries and run in production mode
run: build
	./go-groom \
		-script script.sh \
		-target github.com/kenju \
		-concurrency 8

## Setup
setup:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help
	go get github.com/rakyll/gotest


## Update dependencies
update:
	GO111MODULE=on go mod tidy

## Run tests
test: setup
	if gotest ./... -v; then \
		echo "$(OK_COLOR)$(OK_STRING) go test succeeded$(NO_COLOR)"; \
	else \
		echo "$(ERROR_COLOR)$(ERROR_STRING) go test failed$(NO_COLOR)n"; \
	fi

## Lint
lint:
	go vet $$(go list)
	for pkg in $$(go list); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

## Format source codes
fmt:
	goimports -w main.go

## Build binaries
build:
	go build -ldflags "$(LDFLAGS)"

## View traced file generated by runtime/trace
view-trace:
	go tool trace trace.out

## Update CHANGELOG.md
changelog:
	auto-changelog -t keepachangelog
	git commit -am "update CHANGELOG"

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help
