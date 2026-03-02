#!/usr/bin/env bash
set -euo pipefail

STEP="clean"
export STEP
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

ensure_artifact_dir
run_logged "Unable to clean artifacts directory." bash -lc 'find artifacts -type f -name "*.log" ! -name "clean.log" -delete'
