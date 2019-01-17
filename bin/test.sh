#!/bin/bash
set -e

# Code Climate tool requires the file to be named c.out and to be in the project root
COVERAGE_OUTPUT=c.out
TEMP_COVERAGE=build/temp_cover.out
HTML_REPORT=build/coverage.html

echo "mode: set" > $COVERAGE_OUTPUT

if [ -f $TEMP_COVERAGE ]; then
  rm $TEMP_COVERAGE
fi

PACKAGES=$(go list ./... | grep -v /vendor/)
for package in ${PACKAGES}
do
  if [[ "$package" == *"integrationtests"* ]]; then
    continue
  fi

  go test -coverprofile=$TEMP_COVERAGE $package
  if [ -f $TEMP_COVERAGE ]; then
    cat $TEMP_COVERAGE | grep -v "mode:" | sort -r >> $COVERAGE_OUTPUT
    rm $TEMP_COVERAGE
  fi
done

if [ -f $HTML_REPORT ]; then
  rm $HTML_REPORT
fi

go tool cover -html=$COVERAGE_OUTPUT -o $HTML_REPORT