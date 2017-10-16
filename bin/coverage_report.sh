#!/bin/bash

# Inspired by: https://mlafeldt.github.io/blog/test-coverage-in-go/

GO_PACKAGES=$(go list ./...)
BUILD_DIR=build

rm -Rf $BUILD_DIR
mkdir $BUILD_DIR

for package in $GO_PACKAGES; do
  output_file="$workdir/$(echo $package | tr / -).cover"
  go test "--coverprofile=$BUILD_DIR/$output_file" $package
done

echo "mode: set" > all.cover
grep -h -v "^mode:" "$BUILD_DIR"/*.cover >> all.cover
mv all.cover $BUILD_DIR/all.cover

go tool cover "-html=$BUILD_DIR/all.cover" -o $BUILD_DIR/coverage.html
