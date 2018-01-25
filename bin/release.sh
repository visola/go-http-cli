#!/bin/bash

echo --- Running the build ---
./gradlew test buildPackages

echo --- Tagging commit ---
git tag "$(date +'%Y%m%d%H%M%S')-$(git log --format=%h -1)"
git push --tags
