#!/usr/bin/env bash
set -euo pipefail

STEP="smoke"
export STEP
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

require_command go "Install Go and rerun make smoke."
run_logged "Run cleo --version and ensure output format is valid." go run ./cmd/cleo --version

log_file="${ARTIFACT_DIR}/${STEP}.log"
if ! grep -Eq '^cleo [^ ]+$' "${log_file}"; then
  event failure "log=${log_file} hint=\"Expected output format: cleo <version>.\""
  exit 1
fi
