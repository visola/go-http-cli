#!/bin/bash

echo --- Running the build ---
bin/build.sh

echo --- Tagging commit ---
git tag "v1.0.1"
git push --tags
