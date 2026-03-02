#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${STEP:-}" ]]; then
  echo "STEP must be set before sourcing scripts/lib.sh" >&2
  exit 1
fi

ARTIFACT_DIR="${ARTIFACT_DIR:-artifacts}"
EVENT_VERSION="1"

ensure_artifact_dir() {
  mkdir -p "${ARTIFACT_DIR}"
}

event() {
  local state="$1"
  shift
  local ts
  ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  printf 'CLEO_EVENT v=%s ts=%s step=%s state=%s %s\n' "${EVENT_VERSION}" "${ts}" "${STEP}" "${state}" "$*"
}

require_command() {
  local command_name="$1"
  local hint="$2"
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    event failure "reason=missing_command command=${command_name} hint=\"${hint}\""
    exit 1
  fi
}

run_logged() {
  local hint="$1"
  shift

  local log_file="${ARTIFACT_DIR}/${STEP}.log"
  ensure_artifact_dir

  event start "command=\"$*\" log=${log_file}"
  if "$@" >"${log_file}" 2>&1; then
    cat "${log_file}"
    event success "command=\"$*\" log=${log_file}"
    return 0
  fi

  cat "${log_file}" >&2
  event failure "command=\"$*\" log=${log_file} hint=\"${hint}\""
  return 1
}
