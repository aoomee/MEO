#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_DIR="${OUTPUT_DIR:-${ROOT_DIR}/build}"

command -v go >/dev/null 2>&1 || { echo "Go 1.26+ is required" >&2; exit 1; }
test -f "${ROOT_DIR}/miaomiaowux/internal/web/dist/index.html" || { echo "missing embedded frontend" >&2; exit 1; }

mkdir -p "${OUTPUT_DIR}"

for arch in amd64 arm64; do
  echo "Building panel linux/${arch}"
  (
    cd "${ROOT_DIR}/miaomiaowux"
    CGO_ENABLED=0 GOOS=linux GOARCH="${arch}" \
      go build -trimpath -ldflags="-s -w" -o "${OUTPUT_DIR}/mmwx-linux-${arch}" ./cmd/server
  )

  echo "Building agent linux/${arch}"
  (
    cd "${ROOT_DIR}/mmw-agent"
    CGO_ENABLED=0 GOOS=linux GOARCH="${arch}" \
      go build -trimpath -ldflags="-s -w" -o "${OUTPUT_DIR}/mmwx-agent-linux-${arch}" ./cmd/mmw-agent
  )
done

if command -v sha256sum >/dev/null 2>&1; then
  (cd "${OUTPUT_DIR}" && sha256sum mmwx-linux-* mmwx-agent-linux-* > SHA256SUMS)
else
  (cd "${OUTPUT_DIR}" && shasum -a 256 mmwx-linux-* mmwx-agent-linux-* > SHA256SUMS)
fi
echo "Artifacts written to ${OUTPUT_DIR}"
