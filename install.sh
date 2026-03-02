#!/usr/bin/env bash
set -euo pipefail

REPO_URL="${REPO_URL:-https://github.com/cafaye/cleo.git}"
BRANCH="${BRANCH:-master}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
NON_INTERACTIVE="${NON_INTERACTIVE:-0}"

need() {
  command -v "$1" >/dev/null 2>&1
}

confirm() {
  local q="$1"
  if [[ "$NON_INTERACTIVE" == "1" ]]; then
    echo "$q [auto: yes]"
    return 0
  fi
  read -r -p "$q [y/N]: " ans
  [[ "${ans,,}" == "y" || "${ans,,}" == "yes" ]]
}

install_pkg() {
  local pkg="$1"
  local resolved_pkg
  resolved_pkg="$(map_pkg_name "$pkg")"
  if need brew; then
    brew install "$resolved_pkg"
    return
  fi
  if need apt-get; then
    sudo apt-get update && sudo apt-get install -y "$resolved_pkg"
    return
  fi
  if need dnf; then
    sudo dnf install -y "$resolved_pkg"
    return
  fi
  if need yum; then
    sudo yum install -y "$resolved_pkg"
    return
  fi
  echo "No supported package manager found to install $resolved_pkg"
  exit 1
}

map_pkg_name() {
  local pkg="$1"
  if need apt-get && [[ "$pkg" == "go" ]]; then
    echo "golang-go"
    return
  fi
  echo "$pkg"
}

ensure_dep() {
  local bin="$1"
  if need "$bin"; then
    echo "[ok] $bin"
    return
  fi
  echo "[missing] $bin"
  if confirm "Install $bin now?"; then
    install_pkg "$bin"
  else
    echo "Cannot continue without $bin"
    exit 1
  fi
}

echo "==> Cleo one-command install"
ensure_dep git
ensure_dep go
ensure_dep gh

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

echo "==> Cloning cleo"
git clone --depth 1 --branch "$BRANCH" "$REPO_URL" "$tmpdir/cleo"

mkdir -p "$INSTALL_DIR"
echo "==> Building cleo"
(
  cd "$tmpdir/cleo"
  go build -o "$INSTALL_DIR/cleo" ./cmd/cleo
)
chmod +x "$INSTALL_DIR/cleo"

echo "==> Running cleo setup"
if [[ "$NON_INTERACTIVE" == "1" ]]; then
  "$INSTALL_DIR/cleo" setup --yes --non-interactive --skip-auth
else
  "$INSTALL_DIR/cleo" setup
fi

echo "==> Installed: $INSTALL_DIR/cleo"
if ! echo ":$PATH:" | grep -q ":$INSTALL_DIR:"; then
  echo "Add this to your shell profile: export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "Done. Try: cleo --version"
