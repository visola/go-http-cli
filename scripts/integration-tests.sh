#!/bin/bash
set -e

# Kill daemon if running
if pgrep go-http-daemon >/dev/null; then go-http-daemon --kill; fi

BASE_DIR=$(pwd)
INTEGRATION_TEST_DIR=$BASE_DIR/test/integration
export EXECUTION_DIR=$INTEGRATION_TEST_DIR/execution

rm -Rf $EXECUTION_DIR
mkdir $EXECUTION_DIR

# Build binaries to be tested
go build -o $EXECUTION_DIR/go-http-daemon ./cmd/go-http-daemon
go build -o $EXECUTION_DIR/http ./cmd/http

# Run integration tests
go test ./test/integration

# Kill daemon after integration tests
go-http-daemon --kill
