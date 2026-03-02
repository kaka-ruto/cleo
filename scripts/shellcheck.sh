#!/usr/bin/env bash
set -euo pipefail

STEP="shellcheck"
export STEP
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

require_command shellcheck "Install shellcheck and rerun make shellcheck."
run_logged "Fix shellcheck findings in scripts and hooks." shellcheck scripts/*.sh .githooks/post-push install.sh
