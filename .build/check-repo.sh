#!/bin/bash
set -xeo pipefail

go mod tidy

# this will return 1 if there is any differences.
git diff-index --quiet HEAD
