#!/bin/bash
set -ex

build_and_zip() {
    # $1 -> operating system
    # $2 -> architecture
    # $3 -> OS alias, used in the output file name
    # $4 -> Optional extension with ".", e.g.: .exe
    PACKAGE_DIR=build/$3_$2
    mkdir -p $PACKAGE_DIR

    # build http cli
    PACKAGE_FILE=$PACKAGE_DIR/http$4
    GOOS=$1 GOARCH=$2 go build -o $PACKAGE_FILE ./binaries/http

    # build go-http-daemon
    PACKAGE_FILE=$PACKAGE_DIR/go-http-daemon$4
    GOOS=$1 GOARCH=$2 go build -o $PACKAGE_FILE ./binaries/go-http-daemon

    # build go-http-completion
    PACKAGE_FILE=$PACKAGE_DIR/go-http-completion$4
    GOOS=$1 GOARCH=$2 go build -o $PACKAGE_FILE ./binaries/go-http-completion

    # Copy license and readme
    cp LICENSE $PACKAGE_DIR/
    cp README.md $PACKAGE_DIR/

    zip -j build/$3_$2.zip $PACKAGE_DIR/*
    rm -Rf $PACKAGE_DIR
}

build_and_zip darwin amd64 mac
build_and_zip linux amd64 linux
build_and_zip windows amd64 win .exe
