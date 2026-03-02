#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
OUT_DIR="${2:-dist/release}"

if [[ -z "$VERSION" ]]; then
  echo "usage: $0 <version> [out_dir]"
  exit 1
fi

mkdir -p "$OUT_DIR"
rm -f "$OUT_DIR"/cleo_"$VERSION"_*.tar.gz "$OUT_DIR"/checksums.txt

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
)

for target in "${targets[@]}"; do
  os="${target% *}"
  arch="${target#* }"
  bin="$OUT_DIR/cleo_${VERSION}_${os}_${arch}"
  tarball="${bin}.tar.gz"

  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o "$bin" ./cmd/cleo
  tar -C "$OUT_DIR" -czf "$tarball" "$(basename "$bin")"
  rm -f "$bin"
done

(
  cd "$OUT_DIR"
  shasum -a 256 cleo_"$VERSION"_*.tar.gz > checksums.txt
)

echo "Built assets in $OUT_DIR"
