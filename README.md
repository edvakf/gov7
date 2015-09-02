# gov7

Go binding for V7 JavaScript engine.

[![Build Status](https://travis-ci.org/edvakf/gov7.svg)](https://travis-ci.org/edvakf/gov7)

## Updating v7

v7 was added as a git subtree like this;

```
$ git remote add -f v7_origin git@github.com:cesanta/v7.git
$ git subtree add --prefix=v7/ --squash v7_origin master
```

To update v7, do `merge --squash`

```
$ git merge --squash -s subtree --no-commit v7_origin
$ git commit
```
