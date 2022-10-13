#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

kubectl version

ARTIFACTS_PATH=${ARTIFACTS_PATH:-"${HOME}/e2e-logs"}
mkdir -p "$ARTIFACTS_PATH"

# Install ginkgo
go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo@v2.2.0

# Pre run e2e for extra components
echo "Run pre run e2e"
REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
sh "${REPO_ROOT}"/hack/pre-run-e2e.sh

# Run e2e
echo "Run e2e"
export GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn

set +e
ginkgo -v --race --trace --fail-fast -p --randomize-all ./test/e2e/
TESTING_RESULT=$?

# todo(chengxiangdong): Collect logs
echo "Collected logs at $ARTIFACTS_PATH:"
ls -al "$ARTIFACTS_PATH"

# Post run e2e for delete extra components
echo "Run post run e2e"
sh "${REPO_ROOT}"/hack/post-run-e2e.sh

exit $TESTING_RESULT
