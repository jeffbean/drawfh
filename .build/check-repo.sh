#!/bin/bash
set -xeo pipefail

bazel run //:gazelle -- update-repos -from_file=go.mod -prune

# this will return 1 if there is any differences.
git diff-index --quiet HEAD
