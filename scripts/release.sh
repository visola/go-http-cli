#!/bin/bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

./semantic-release --token $GITHUB_TOKEN --slug visola/go-http-cli --ghr --vf
export VERSION=$(cat .version)

ls -la
cat .version
cat .ghr

$SCRIPT_DIR/package.sh
ghr $(cat .ghr) build/*.zip
