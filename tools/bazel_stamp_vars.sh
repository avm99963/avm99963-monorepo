#!/usr/bin/env bash

echo "STABLE_REVISION $(git rev-list --count --first-parent HEAD)"
