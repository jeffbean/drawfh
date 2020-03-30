#!/bin/bash
set -xeo pipefail
# this will return 1 if there is any differences.
git diff-index --quiet HEAD
