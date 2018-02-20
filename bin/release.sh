#!/bin/bash

echo --- Running the build ---
./gradlew test buildPackages

echo --- Tagging commit ---
git tag "v1.0.1"
git push --tags
