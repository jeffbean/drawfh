#!/bin/bash
set -xeo pipefail

readonly BASE_DIR=$(git rev-parse --show-toplevel)

# "${BASE_DIR}/.build/check-repo.sh"

gcloud --project=jeffbeandev builds submit --config="${BASE_DIR}/.build/cloudbuild.publish.yaml" "${BASE_DIR}"
