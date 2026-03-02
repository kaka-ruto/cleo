#!/usr/bin/env bash
set -euo pipefail

REPO_URL="${REPO_URL:-https://github.com/cafaye/cleo.git}"
BRANCH="${BRANCH:-master}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
NON_INTERACTIVE="${NON_INTERACTIVE:-0}"
REQUIRED_GO_VERSION="${REQUIRED_GO_VERSION:-1.25.1}"
if [[ "$NON_INTERACTIVE" == "1" ]]; then
  export DEBIAN_FRONTEND=noninteractive
fi

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
  echo "$pkg"
}

go_version() {
  if ! need go; then
    echo ""
    return
  fi
  go version | sed -E 's/^go version go([0-9]+\.[0-9]+(\.[0-9]+)?).*/\1/'
}

version_gte() {
  local have="$1"
  local want="$2"
  [[ "$(printf '%s\n%s\n' "$want" "$have" | sort -V | head -n1)" == "$want" ]]
}

install_go_toolchain() {
  local os arch url tmp archive
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  case "$arch" in
    x86_64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
  esac
  url="https://go.dev/dl/go${REQUIRED_GO_VERSION}.${os}-${arch}.tar.gz"
  tmp="$(mktemp -d)"
  archive="${tmp}/go.tgz"
  echo "Installing Go ${REQUIRED_GO_VERSION} from ${url}"
  curl -fsSL "$url" -o "$archive"
  mkdir -p "$HOME/.local"
  rm -rf "$HOME/.local/go"
  tar -C "$HOME/.local" -xzf "$archive"
  rm -rf "$tmp"
  export PATH="$HOME/.local/go/bin:$PATH"
}

ensure_go() {
  local have
  have="$(go_version)"
  if [[ -n "$have" ]] && version_gte "$have" "$REQUIRED_GO_VERSION"; then
    echo "[ok] go ${have}"
    return
  fi
  echo "[missing/outdated] go (need >= ${REQUIRED_GO_VERSION}, found: ${have:-none})"
  if confirm "Install Go ${REQUIRED_GO_VERSION} now?"; then
    install_go_toolchain
  else
    echo "Cannot continue without Go ${REQUIRED_GO_VERSION}+"
    exit 1
  fi
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
ensure_go
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
  "$INSTALL_DIR/cleo" setup --non-interactive
else
  "$INSTALL_DIR/cleo" setup
fi

echo "==> Installed: $INSTALL_DIR/cleo"
if ! echo ":$PATH:" | grep -q ":$INSTALL_DIR:"; then
  echo "Add this to your shell profile: export PATH=\"$INSTALL_DIR:\$PATH\""
fi
if ! echo ":$PATH:" | grep -q ":$HOME/.local/go/bin:"; then
  echo "If Go was installed, add: export PATH=\"$HOME/.local/go/bin:\$PATH\""
fi

echo "Done. Try: cleo --version"
