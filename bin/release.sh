#!/bin/bash

echo --- Running the build ---
./gradlew test buildPackages

echo --- Tagging commit ---
git tag "v0.9.2"
git push --tags
