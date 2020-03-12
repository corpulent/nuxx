# Nuxx CLI

## Install gvm

```shell script
$ bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
$ source "$HOME/.gvm/scripts/gvm"
$ gvm install go1.13.6
$ gvm use go1.13.6 [--default]
```

# Update .bashrc and source it

```shell script
#!/bin/bash

GVM_BIN="$HOME/.gvm/scripts/gvm"

if [[ -f "$GVM_BIN" ]]; then
    source "$GVM_BIN"
fi

export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=${PATH}:$GOBIN
```

# Requirements

```shell script
$ go get -u -d github.com/spf13/cobra/cobra
$ go get -u -d github.com/magefile/mage
```

# Installing and testing

```shell script

# install from github source
$ go install nuxx-cli

# install from project directory
$ go install

# run from project directory
$ go run main.go
```
