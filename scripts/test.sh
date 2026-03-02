#!/usr/bin/env bash
set -euo pipefail

STEP=test
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

require_command go "Install Go and rerun make test."
run_logged "Run go test ./... locally and fix failing tests." go test -race -count=1 ./...
