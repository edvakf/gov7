# gov7

Go binding for V7 JavaScript engine.

[![Build Status](https://travis-ci.org/edvakf/gov7.svg)](https://travis-ci.org/edvakf/gov7) [![Coverage Status](https://coveralls.io/repos/edvakf/gov7/badge.svg?branch=master&service=github)](https://coveralls.io/github/edvakf/gov7?branch=master)

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
