#!/bin/bash
set -xeo pipefail

readonly BASE_DIR=$(git rev-parse --show-toplevel)

gcloud --project=jeffbeandev builds submit --config="${BASE_DIR}/cloudbuild.yaml" "${BASE_DIR}"
