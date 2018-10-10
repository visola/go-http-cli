#!/bin/bash
set -ex

# Kill daemon if running
if pgrep go-http-daemon; then go-http-daemon --kill; fi

BASE_DIR=$(pwd)
INTEGRATION_TEST_DIR=$BASE_DIR/integration-tests
EXECUTION_DIR=$INTEGRATION_TEST_DIR/execution

rm -Rf $EXECUTION_DIR
mkdir $EXECUTION_DIR

# Build binaries to be tested
go build -o $EXECUTION_DIR/go-http-daemon ./binaries/go-http-daemon
go build -o $EXECUTION_DIR/http ./binaries/http

# Build integration test runner
cd $INTEGRATION_TEST_DIR
go build -o $EXECUTION_DIR/integration-tests .

# Prepare for the test
export PATH=$EXECUTION_DIR
cd $EXECUTION_DIR
./integration-tests
