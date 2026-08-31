#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/usr/local/bin"
DATA_DIR="/var/lib/mmwx/data"
CONFIG_DIR="/etc/mmwx"
SERVICE_NAME="mmwx"
PORT_VALUE="12889"
REPO="${MMWX_UPDATE_REPO:-}"
TOKEN="${MMWX_GH_TOKEN:-}"
RELEASE_TAG="${MMWX_RELEASE_TAG:-mmwx}"
BINARY_FILE=""
ACTION="install"

usage() {
  cat <<'USAGE'
Usage:
  sudo ./deploy/install.sh [install] [options]
  sudo ./deploy/install.sh update [options]
  sudo ./deploy/install.sh uninstall

Options:
  --binary PATH       install an already-built panel binary
  --repo OWNER/REPO   download from a GitHub Release repository
  --token TOKEN       token used to read a private GitHub repository
  --tag TAG           release tag (default: mmwx)
  --port PORT         HTTP port (default: 12889)

For a private repository, prefer --binary or provide --repo and --token.
USAGE
}

if [[ $# -gt 0 && "$1" != --* ]]; then
  ACTION="$1"
  shift
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --binary) BINARY_FILE="${2:?missing path}"; shift 2 ;;
    --repo) REPO="${2:?missing owner/repo}"; shift 2 ;;
    --token) TOKEN="${2:?missing token}"; shift 2 ;;
    --tag) RELEASE_TAG="${2:?missing tag}"; shift 2 ;;
    --port) PORT_VALUE="${2:?missing port}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage; exit 2 ;;
  esac
done

[[ "${EUID}" -eq 0 ]] || { echo "Run this script with sudo/root." >&2; exit 1; }
[[ "$(uname -s)" == "Linux" ]] || { echo "Only Linux is supported." >&2; exit 1; }

if [[ "${ACTION}" == "uninstall" ]]; then
  systemctl disable --now "${SERVICE_NAME}" 2>/dev/null || true
  rm -f "/etc/systemd/system/${SERVICE_NAME}.service" "${INSTALL_DIR}/mmwx"
  systemctl daemon-reload
  echo "Program and service removed. Persistent data remains in /var/lib/mmwx."
  exit 0
fi

[[ "${ACTION}" == "install" || "${ACTION}" == "update" ]] || { usage; exit 2; }
[[ "${PORT_VALUE}" =~ ^[0-9]+$ ]] && (( PORT_VALUE >= 1 && PORT_VALUE <= 65535 )) || {
  echo "Invalid port: ${PORT_VALUE}" >&2
  exit 2
}

case "$(uname -m)" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

STAGING_DIR="$(mktemp -d /tmp/mmwx-install.XXXXXX)"
trap 'rm -rf "${STAGING_DIR}"' EXIT
STAGED_BINARY="${STAGING_DIR}/mmwx"

if [[ -n "${BINARY_FILE}" ]]; then
  [[ -f "${BINARY_FILE}" ]] || { echo "Binary not found: ${BINARY_FILE}" >&2; exit 1; }
  cp "${BINARY_FILE}" "${STAGED_BINARY}"
elif [[ -n "${REPO}" ]]; then
  ASSET="mmwx-linux-${ARCH}"
  URL="https://github.com/${REPO}/releases/download/${RELEASE_TAG}/${ASSET}"
  CURL_ARGS=(--fail --silent --show-error --location --retry 3 --output "${STAGED_BINARY}")
  if [[ -n "${TOKEN}" ]]; then
    CURL_ARGS+=(--header "Authorization: Bearer ${TOKEN}")
  fi
  echo "Downloading ${ASSET} from ${REPO}@${RELEASE_TAG}"
  curl "${CURL_ARGS[@]}" "${URL}"
else
  echo "Provide --binary PATH or --repo OWNER/REPO." >&2
  exit 2
fi

chmod 0755 "${STAGED_BINARY}"

getent group mmwx >/dev/null 2>&1 || groupadd --system mmwx
id mmwx >/dev/null 2>&1 || useradd --system --gid mmwx --home-dir /var/lib/mmwx --shell /usr/sbin/nologin mmwx
install -d -m 0750 -o mmwx -g mmwx "${DATA_DIR}"
install -d -m 0750 -o root -g mmwx "${CONFIG_DIR}"

if [[ ! -f "${CONFIG_DIR}/mmwx.env" ]]; then
  if command -v openssl >/dev/null 2>&1; then
    JWT_VALUE="$(openssl rand -hex 32)"
  else
    JWT_VALUE="$(od -An -N32 -tx1 /dev/urandom | tr -d ' \n')"
  fi
  {
    printf 'PORT=%s\n' "${PORT_VALUE}"
    printf 'BIND_HOST=0.0.0.0\n'
    printf 'MMWX_DATA_DIR=%s\n' "${DATA_DIR}"
    printf 'TZ=Asia/Shanghai\n'
    printf 'JWT_SECRET=%s\n' "${JWT_VALUE}"
    printf 'MMWX_UPDATE_REPO=%s\n' "${REPO:-off}"
    printf 'MMWX_AGENT_GITHUB_REPO=%s\n' "${REPO:-off}"
    [[ -n "${TOKEN}" ]] && printf 'MMWX_GH_TOKEN=%s\n' "${TOKEN}"
  } > "${CONFIG_DIR}/mmwx.env"
  chmod 0640 "${CONFIG_DIR}/mmwx.env"
  chown root:mmwx "${CONFIG_DIR}/mmwx.env"
elif [[ "${ACTION}" == "install" ]]; then
  echo "Keeping existing ${CONFIG_DIR}/mmwx.env"
fi

if [[ -f "${INSTALL_DIR}/mmwx" ]]; then
  cp "${INSTALL_DIR}/mmwx" "${STAGING_DIR}/mmwx.previous"
fi
install -m 0755 "${STAGED_BINARY}" "${INSTALL_DIR}/mmwx"
install -m 0644 "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/systemd/mmwx.service" "/etc/systemd/system/${SERVICE_NAME}.service"

systemctl daemon-reload
systemctl enable --now "${SERVICE_NAME}"
if ! systemctl restart "${SERVICE_NAME}"; then
  if [[ -f "${STAGING_DIR}/mmwx.previous" ]]; then
    install -m 0755 "${STAGING_DIR}/mmwx.previous" "${INSTALL_DIR}/mmwx"
    systemctl restart "${SERVICE_NAME}" || true
    echo "Update failed; previous binary restored." >&2
  fi
  exit 1
fi

echo "MMWX is running at http://SERVER_IP:${PORT_VALUE}"
echo "Data: ${DATA_DIR}"
echo "Config: ${CONFIG_DIR}/mmwx.env"
