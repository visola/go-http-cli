#!/bin/bash
set -ex

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

$SCRIPT_DIR/clean.sh

if hash go-http-daemon 2>/dev/null; then
  # Kill daemon
  go-http-daemon --kill
fi

$SCRIPT_DIR/update-dependencies.sh
$SCRIPT_DIR/test.sh
$SCRIPT_DIR/integration-tests.sh
$SCRIPT_DIR/generate-docs.sh
$SCRIPT_DIR/package.sh
