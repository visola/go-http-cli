#!/bin/bash
export GOPATH=/home/appuser/goworkspace
export PATH=$PATH:/usr/lib/go-1.9/bin:$GOPATH/bin

echo Who am I?
whoami

mkdir -pv $GOPATH/src/github.com/visola
mkdir -pv $GOPATH/bin
mkdir -pv $GOPATH/pkg

echo $PATH
go env

echo Fetching dependencies...
./gradlew --no-daemon updateDependencies updateLinter

echo Installing...
./gradlew --no-daemon install

echo Command installed correctly?
which http

echo Run the full build
./gradlew --no-daemon build
