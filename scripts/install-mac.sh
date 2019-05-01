#!/bin/bash
set -ex

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

$SCRIPT_DIR/build.sh

unzip -d build/mac build/mac_amd64.zip
cp build/mac/http /usr/local/bin/
cp build/mac/go-http-completion /usr/local/bin/
cp build/mac/go-http-daemon /usr/local/bin/
