#!/usr/bin/env sh 

git rev-parse --show-toplevel
git rev-parse --abbrev-ref HEAD

git checkout master
git pull --prune
