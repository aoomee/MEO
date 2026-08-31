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

MINIMAL_CSS="${DIST_DIR}/assets/minimal.css"
for REQUIRED_RULE in '.anime-starfield' '.nav-icon-anime' 'theme-minimal'; do
  if ! grep -F -q -- "${REQUIRED_RULE}" "${MINIMAL_CSS}"; then
    echo "minimal theme is missing legacy visual guard: ${REQUIRED_RULE}" >&2
    exit 1
  fi
done
if ! grep -F -q 'lockMinimalTheme' "${DIST_DIR}/index.html"; then
  echo "index.html is missing the minimal-theme lock" >&2
  exit 1
fi

# The active dashboard and settings chunks must not re-enable legacy visual
# modes from stale localStorage/server state. Keep the PRO probe feature, but
# reject only the removed theme selectors/classes in the chunks users execute.
for ACTIVE_FILE in "${DIST_DIR}/assets/index-ZdYIAvvZ.js" "${DIST_DIR}/assets/settings-"*.js; do
  if [[ -f "${ACTIVE_FILE}" ]] && grep -E -q "theme-(anime|premium|pixel)|'value':'(anime|premium|pixel)'" "${ACTIVE_FILE}"; then
    echo "legacy theme selector found in active bundle: ${ACTIVE_FILE}" >&2
    exit 1
  fi
done

SYSTEM_SETTINGS_FILE="${DIST_DIR}/assets/system-settings-"*.js
for ACTIVE_FILE in ${SYSTEM_SETTINGS_FILE}; do
  if [[ -f "${ACTIVE_FILE}" ]] && grep -E -q "Premium|theme-(anime|premium|pixel)|[\"']pixel[\"']" "${ACTIVE_FILE}"; then
    echo "legacy/premium theme text found in system settings bundle: ${ACTIVE_FILE}" >&2
    exit 1
  fi
done

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
