#!/usr/bin/env bash
set -euo pipefail

STEP=fmt
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

require_command go "Install Go and rerun make fmt."

go_file_list="$(mktemp)"
find . -type f -name '*.go' -not -path './.git/*' | LC_ALL=C sort >"${go_file_list}"
trap 'rm -f "${go_file_list}"' EXIT
log_file="${ARTIFACT_DIR}/${STEP}.log"
ensure_artifact_dir

if [[ ! -s "${go_file_list}" ]]; then
  : >"${log_file}"
  event success "log=${log_file} message=\"no_go_files\""
  exit 0
fi

event start "command=\"gofmt -l\" log=${log_file}"
unformatted="$(xargs gofmt -l <"${go_file_list}")"
printf '%s\n' "${unformatted}" >"${log_file}"

if [[ -n "${unformatted}" ]]; then
  cat "${log_file}" >&2
  event failure "log=${log_file} hint=\"Run gofmt -w on listed files and retry.\""
  exit 1
fi

event success "log=${log_file}"
