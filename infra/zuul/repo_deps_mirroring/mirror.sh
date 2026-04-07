#!/bin/bash

# NOTE: Call this script via Copybara!
#
# ```
# bazel run //infra/zuul/repo_deps_mirroring:mirror -- WORKFLOW_NAME
# ```

if [[ $# -lt 1 ]]; then
  echo "Please supply an argument with the workflow name to run" >&2
  echo "" >&2
  echo "USAGE: bazel run //infra/zuul/repo_deps_mirroring:mirror -- WORKFLOW_NAME [optional flags passed to copybara]" >&2
fi

WORKFLOW_NAME="$1"
shift 1

COPYBARA_BIN=$(rlocation "_main/tools/copybara/copybara")
CONFIG_FILE=$(rlocation "_main/infra/zuul/repo_deps_mirroring/copy.bara.sky")

"${COPYBARA_BIN}" "${CONFIG_FILE}" "${WORKFLOW_NAME}" --git-committer-email "copybara-bot@avm99963.com" --git-committer-name "Copybara bot" "$@"
