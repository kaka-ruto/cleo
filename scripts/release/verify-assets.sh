#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
OUT_DIR="${2:-dist/release}"
BINARY_NAME="${BINARY_NAME:-cleo}"

if [[ -z "$VERSION" ]]; then
  echo "usage: $0 <version> [out_dir]"
  exit 1
fi

required=(
  "${BINARY_NAME}_${VERSION}_linux_amd64.tar.gz"
  "${BINARY_NAME}_${VERSION}_linux_arm64.tar.gz"
  "${BINARY_NAME}_${VERSION}_darwin_amd64.tar.gz"
  "${BINARY_NAME}_${VERSION}_darwin_arm64.tar.gz"
  "checksums.txt"
)

for file in "${required[@]}"; do
  if [[ ! -f "$OUT_DIR/$file" ]]; then
    echo "missing artifact: $OUT_DIR/$file"
    exit 1
  fi
done

(
  cd "$OUT_DIR"
  shasum -a 256 -c checksums.txt
)

echo "Artifacts verified in $OUT_DIR"
