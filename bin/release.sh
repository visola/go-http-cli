#!/bin/bash

echo --- Running the build ---
./gradlew test buildPackages

echo --- Tagging commit ---
git tag "v1.0.0"
git push --tags
