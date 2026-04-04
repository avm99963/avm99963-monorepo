#!/usr/bin/env bash

STABLE_REVISION=$(git rev-list --count --first-parent HEAD)
if [[ $? != 0 ]]; then
 echo "Cannot construct STABLE_REVISION."
 exit 1
fi
echo "STABLE_REVISION ${STABLE_REVISION}"

echo "STABLE_GIT_COMMIT $(git rev-parse HEAD 2>/dev/null || echo 'unknown')"
echo "STABLE_GIT_DIRTY $(if git diff --quiet 2>/dev/null; then echo 'clean'; else echo 'dirty'; fi)"
# This call is posix-compatible.
echo "FORMATTED_DATE_RFC3339 $(date -u "+%Y-%m-%dT%TZ")"
