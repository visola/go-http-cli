#!/bin/bash

if [[ "$#" == 1 && $1 == "clean" ]]; then
  docker build --no-cache -t go-http-cli-test-build-and-run .
else
  docker build -t go-http-cli-test-build-and-run .
fi

docker run go-http-cli-test-build-and-run http https://www.google.com
