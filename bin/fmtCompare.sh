#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

project_dir=$DIR/..
pushd $project_dir > /dev/null

project_dir=$(pwd)
build_dir=$project_dir/build
output_dir=$build_dir/fmt

# Clear output directory
rm -Rf $output_dir/*

# Ensure all directories exist
dirs=("$build_dir" "$output_dir" "$output_dir/config" "$output_dir/config/yaml")
for dir in "${dirs[@]}"; do
  if [ ! -d $dir ]; then
    mkdir $dir
  fi
done

# gofmt all files
for file in $(find . -name '*.go'); do
  gofmt -s $file > $output_dir/$file
done

# Compare files
diff $output_dir/main.go main.go
diff -r $output_dir/config/ config/

rm -Rf $output_dir

popd > /dev/null
