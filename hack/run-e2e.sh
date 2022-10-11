#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ARTIFACTS_PATH=${ARTIFACTS_PATH:-"${HOME}/e2e-logs"}
mkdir -p "$ARTIFACTS_PATH"

# Install ginkgo
go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@latest

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${REPO_ROOT}"

# Run e2e
export GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn

set +e
ginkgo -v --race --trace --fail-fast -p --randomize-all ./test/e2e/
TESTING_RESULT=$?

# Collect logs
echo "Collected logs at $ARTIFACTS_PATH:"
ls -al "$ARTIFACTS_PATH"

exit $TESTING_RESULT
