#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="${ROOT_DIR}/miaomiaowux/internal/web/dist"

test -f "${DIST_DIR}/index.html" || {
  echo "missing embedded frontend entry" >&2
  exit 1
}

TOPBAR_COUNT="$(find "${DIST_DIR}/assets" -maxdepth 1 -type f -name 'topbar-*.js' -print | wc -l | tr -d ' ')"
if [[ "${TOPBAR_COUNT}" -ne 1 ]]; then
  echo "expected exactly one topbar bundle, found ${TOPBAR_COUNT}" >&2
  exit 1
fi
TOPBAR_FILE="$(find "${DIST_DIR}/assets" -maxdepth 1 -type f -name 'topbar-*.js' -print -quit)"

if grep -R -F -q --include='*.js' '_0x598148' "${DIST_DIR}"; then
  echo "undefined premium-theme guard _0x598148 is present in web dist" >&2
  exit 1
fi

if grep -E -q "theme(Miaomiaowu|Anime|Premium)|'value':'(miaomiaowu|anime|premium)'" "${TOPBAR_FILE}"; then
  echo "removed theme choices are present in the topbar bundle" >&2
  exit 1
fi

while IFS= read -r asset_path; do
  test -f "${DIST_DIR}/${asset_path#/}" || {
    echo "index.html references missing asset: ${asset_path}" >&2
    exit 1
  }
done < <(grep -Eo '/assets/[^"[:space:]]+\.(js|css)' "${DIST_DIR}/index.html" | sort -u)

if command -v node >/dev/null 2>&1; then
  node --input-type=module --check < "${TOPBAR_FILE}"
fi

echo "embedded frontend checks passed"
