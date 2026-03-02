#!/usr/bin/env bash
set -euo pipefail

NON_INTERACTIVE="${NON_INTERACTIVE:-0}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
SCAN_ROOTS="${SCAN_ROOTS:-$HOME/Code}"

has() {
  command -v "$1" >/dev/null 2>&1
}

style() {
  local color="$1"
  shift
  if has gum; then
    gum style --foreground "$color" "$*"
    return
  fi
  echo "$*"
}

title() {
  local text="$1"
  if has gum; then
    gum style --border rounded --padding "1 2" "$text"
    return
  fi
  echo "==> $text"
}

confirm() {
  local prompt="$1"
  if [[ "$NON_INTERACTIVE" == "1" ]]; then
    style 212 "$prompt [auto: yes]"
    return 0
  fi
  if has gum; then
    gum confirm "$prompt"
    return $?
  fi
  read -r -p "$prompt [y/N]: " ans
  [[ "${ans,,}" == "y" || "${ans,,}" == "yes" ]]
}

uniq_paths() {
  awk 'NF && !seen[$0]++'
}

collect_binaries() {
  {
    command -v cleo 2>/dev/null || true
    echo "$INSTALL_DIR/cleo"
  } | uniq_paths
}

remove_file_if_exists() {
  local path="$1"
  if [[ -e "$path" ]]; then
    rm -f "$path"
    style 82 "Removed $path"
  else
    style 240 "Not found: $path"
  fi
}

discover_config_files() {
  local roots
  roots="${SCAN_ROOTS//,/ }"
  for root in $roots; do
    [[ -d "$root" ]] || continue
    find "$root" -type f -name "cleo.yml" -print 2>/dev/null
  done | uniq_paths
}

remove_config_files() {
  local files=("$@")
  if [[ "${#files[@]}" -eq 0 ]]; then
    style 240 "No cleo.yml files found in scan roots: $SCAN_ROOTS"
    return
  fi
  style 212 "Found ${#files[@]} cleo.yml file(s):"
  printf '%s\n' "${files[@]}"
  if confirm "Delete all listed cleo.yml files?"; then
    for file in "${files[@]}"; do
      rm -f "$file"
    done
    style 82 "Deleted ${#files[@]} cleo.yml file(s)."
  else
    style 240 "Skipped cleo.yml cleanup."
  fi
}

title "Cleo Uninstall"

if ! confirm "Remove Cleo binary from your system?"; then
  style 196 "Uninstall canceled."
  exit 1
fi

while IFS= read -r bin_path; do
  remove_file_if_exists "$bin_path"
done < <(collect_binaries)

if confirm "Also remove repo-level cleo.yml files from SCAN_ROOTS ($SCAN_ROOTS)?"; then
  mapfile -t configs < <(discover_config_files)
  remove_config_files "${configs[@]}"
else
  style 240 "Skipped repo config cleanup."
fi

style 82 "Cleo uninstall complete."
