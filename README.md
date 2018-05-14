# go-groom

go-groom run grooming commands against multiple repositories concurrently.

# Usage

```
$ groom
```

Run the following commands to the all repository under `$GOPATH/src`.

TODO: get command from STDIN and run it.

```sh
# update local master branch to the latest
git checkout master
git pull --prune
```
