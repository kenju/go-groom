NAME := groom
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
	-X 'main.revision=$(REVISION)'

## Setup
setup:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help

## Install dependencies
deps:
	echo "Not implemented yet"

## Update dependencies
update:
	echo "Not implemented yet"

## Run tests
test:
	echo "Not implemented yet"

## Lint
lint:
	echo "Not implemented yet"

## Build binaries
build:
	go build -ldflags "$(LDFLAGS)"

## Build binaries and run
run: build
	./go-groom

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint help
