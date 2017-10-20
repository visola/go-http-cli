#/bin/bash

set -e

# Make sure we're in the right dir
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd $DIR/.. > /dev/null

echo "Removing debug files..."
find . -name '*.test' -exec rm {} \;

echo "Compiling..."
go install

echo "Linting..."
golint -set_exit_status ./...

echo "Checking formatting..."
bin/fmtCompare.sh

echo "Testing..."
bin/coverage_report.sh

popd > /dev/null
