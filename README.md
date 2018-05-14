# go-groom

go-groom run grooming commands against multiple repositories concurrently.

# Usage

```
$ groom
```

Run the following commands to the all repository under `$GOPATH/src`.

TODO: this command should be configurable.

```sh
# update local master branch to the latest
git checkout master
git pull --prune

# remove merged branch
git branch --merged | grep -v '*' | xargs -I % git branch -d %
```
