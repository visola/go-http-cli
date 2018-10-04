#!/bin/bash
set -ex

cd integration-tests
go build .
./integration-tests
