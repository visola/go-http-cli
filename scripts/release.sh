#!/bin/bash

echo --- Running the build ---
bin/build.sh

echo --- Tagging commit ---
git tag "v1.0.2"
git push --tags
