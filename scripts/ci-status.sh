#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: ci-status.sh [--latest] [--commit <sha>] [--logs] [--watch] [--interval <seconds>] [--limit <n>]
USAGE
}

require_tools() {
  command -v gh >/dev/null 2>&1 || { echo "gh CLI is required." >&2; exit 1; }
  if ! command -v jq >/dev/null 2>&1; then
    echo "jq is required. Install with: brew install jq (macOS) or sudo apt-get install -y jq (Debian/Ubuntu)." >&2
    exit 1
  fi
  gh auth status >/dev/null 2>&1 || { echo "gh is not authenticated. Run: gh auth login" >&2; exit 1; }
}

latest=false
show_logs=false
watch=false
commit=""
limit=20
interval=15

while [[ $# -gt 0 ]]; do
  case "$1" in
    --latest) latest=true ;;
    --logs) show_logs=true ;;
    --watch) watch=true ;;
    --commit) commit="${2:-}"; shift ;;
    --limit) limit="${2:-20}"; shift ;;
    --interval) interval="${2:-15}"; shift ;;
    --help|-h) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage; exit 1 ;;
  esac
  shift
done

require_tools
repo="$(gh repo view --json nameWithOwner -q '.nameWithOwner')"

find_run() {
  local runs_json
  runs_json="$(gh run list --limit "$limit" --json databaseId,headSha,displayTitle,status,conclusion,workflowName,createdAt,url)"
  if [[ -n "$commit" ]]; then
    runs_json="$(echo "$runs_json" | jq --arg c "$commit" '[.[] | select(.headSha | startswith($c))]')"
  fi
  echo "$runs_json"
}

print_logs() {
  local run_id="$1"
  gh run view "$run_id" --log-failed || gh run view "$run_id" --log
}

runs_json="$(find_run)"
run_id="$(echo "$runs_json" | jq -r '.[0].databaseId // empty')"

if [[ -z "$run_id" && "$watch" == "true" ]]; then
  while [[ -z "$run_id" ]]; do
    sleep "$interval"
    runs_json="$(find_run)"
    run_id="$(echo "$runs_json" | jq -r '.[0].databaseId // empty')"
  done
fi

if [[ -z "$run_id" ]]; then
  echo "No workflow run found." >&2
  exit 1
fi

if [[ "$latest" != "true" ]]; then
  echo "$runs_json" | jq -r '(["ID","SHA","WORKFLOW","STATUS","CONCLUSION","TITLE"] | @tsv),(.[] | [.databaseId, (.headSha[0:8]), .workflowName, .status, (.conclusion // "-"), .displayTitle] | @tsv)' | column -t -s $'\t'
fi

if [[ "$watch" == "true" ]]; then
  while true; do
    run_json="$(gh api "repos/${repo}/actions/runs/${run_id}")"
    status="$(echo "$run_json" | jq -r '.status')"
    conclusion="$(echo "$run_json" | jq -r '.conclusion // empty')"
    echo "run=${run_id} status=${status} conclusion=${conclusion:-pending}"
    if [[ "$status" == "completed" ]]; then
      if [[ "$conclusion" != "success" ]]; then
        [[ "$show_logs" == "true" ]] && print_logs "$run_id"
        exit 1
      fi
      exit 0
    fi
    sleep "$interval"
  done
fi

[[ "$show_logs" == "true" ]] && print_logs "$run_id"
