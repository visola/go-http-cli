#/bin/bash

set -e

CURRENT_DIR=$(pwd)

# Make sure we're in the right dir
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR
cd ..

go install
golint ./...
bin/coverage_report.sh

cd $CURRENT_DIR
