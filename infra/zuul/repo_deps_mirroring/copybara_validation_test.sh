#!/bin/bash

set -e

COPYBARA_BIN=$(rlocation "_main/tools/copybara/copybara")
CONFIG_FILE=$(rlocation "_main/infra/zuul/repo_deps_mirroring/copy.bara.sky")

"${COPYBARA_BIN}" validate "${CONFIG_FILE}"
