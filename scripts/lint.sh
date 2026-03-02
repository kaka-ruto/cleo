#!/usr/bin/env bash
set -euo pipefail

STEP=lint
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

require_command go "Install Go and rerun make lint."
run_logged "Run go vet ./... and fix reported issues." go vet ./...
