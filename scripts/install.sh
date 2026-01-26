#!/usr/bin/env bash
# install.sh - Install gwtt from GitHub release assets
#
# Usage: ./scripts/install.sh [install_dir]
#   install_dir: Optional. Defaults to current directory

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

OWNER="pi2pie"
REPO="git-worktree-tasks"

INSTALL_DIR=""

usage() {
    echo "Usage: ./scripts/install.sh [install_dir]"
}

log_info() { echo -e "${BLUE}ℹ${NC} $1"; }
log_ok() { echo -e "${GREEN}✓${NC} $1"; }
log_warn() { echo -e "${YELLOW}⚠${NC} $1"; }
log_err() { echo -e "${RED}✗${NC} $1"; }

require_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        log_err "Missing required command: $1"
        exit 1
    fi
}

fetch() {
    local url="$1"
    local dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$dest"
        return
    fi
    if command -v wget >/dev/null 2>&1; then
        wget -qO "$dest" "$url"
        return
    fi
    log_err "Either curl or wget is required."
    exit 1
}

sha256_file() {
    local file="$1"
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$file" | awk '{print $1}'
        return
    fi
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$file" | awk '{print $1}'
        return
    fi
    log_err "Missing checksum tool: sha256sum or shasum."
    exit 1
}

while [ $# -gt 0 ]; do
    case "$1" in
        -h|--help)
            usage
            exit 0
            ;;
        *)
            if [ -z "$INSTALL_DIR" ]; then
                INSTALL_DIR="$1"
                shift
            else
                log_err "Unknown argument: $1"
                usage
                exit 1
            fi
            ;;
    esac
done

UNAME_S="$(uname -s 2>/dev/null || echo unknown)"
UNAME_M="$(uname -m 2>/dev/null || echo unknown)"

case "$UNAME_S" in
    Darwin) OS="darwin" ;;
    Linux) OS="linux" ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT) OS="windows" ;;
    *)
        log_err "Unsupported OS: $UNAME_S"
        exit 1
        ;;
esac

case "$UNAME_M" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        log_err "Unsupported architecture: $UNAME_M"
        exit 1
        ;;
esac

if [ -z "$INSTALL_DIR" ]; then
    INSTALL_DIR="$PWD"
fi

if [ ! -d "$INSTALL_DIR" ]; then
    log_info "Creating install directory: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
fi

if [ ! -w "$INSTALL_DIR" ]; then
    log_err "Install directory is not writable: $INSTALL_DIR"
    exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

RELEASE_JSON="$TMP_DIR/release.json"
API_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
log_info "Fetching latest release metadata..."
fetch "$API_URL" "$RELEASE_JSON"

TAG="$(grep -m1 '"tag_name"' "$RELEASE_JSON" | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/' | tr -d '\r')"
if [ -z "$TAG" ]; then
    log_err "Failed to determine latest release tag."
    exit 1
fi
log_info "Latest release: $TAG"

EXT="tar.gz"
if [ "$OS" = "windows" ]; then
    EXT="zip"
fi

CHECKSUMS_NAME="checksums.txt"

ASSET_URL="$(grep -oE '"browser_download_url": *"[^"]+"' "$RELEASE_JSON" | sed -E 's/.*"([^"]+)".*/\1/' | grep -E "${REPO}_.+_${OS}_${ARCH}\\.${EXT}$" | head -n1)"
CHECKSUMS_URL="$(grep -oE '"browser_download_url": *"[^"]+"' "$RELEASE_JSON" | sed -E 's/.*"([^"]+)".*/\1/' | grep -F "$CHECKSUMS_NAME" | head -n1)"

if [ -z "$ASSET_URL" ]; then
    log_err "Release asset not found for ${OS}_${ARCH}.${EXT}"
    exit 1
fi

if [ -z "$CHECKSUMS_URL" ]; then
    log_err "Release checksums not found: $CHECKSUMS_NAME"
    exit 1
fi

ASSET_NAME="$(basename "$ASSET_URL")"
ASSET_FILE="$TMP_DIR/$ASSET_NAME"
CHECKSUMS_FILE="$TMP_DIR/$CHECKSUMS_NAME"

log_info "Downloading $ASSET_NAME..."
fetch "$ASSET_URL" "$ASSET_FILE"
log_info "Downloading checksums..."
fetch "$CHECKSUMS_URL" "$CHECKSUMS_FILE"

EXPECTED_SUM="$(grep -E " ${ASSET_NAME}\$" "$CHECKSUMS_FILE" | awk '{print $1}' | head -n1)"
if [ -z "$EXPECTED_SUM" ]; then
    log_err "Checksum entry not found for $ASSET_NAME"
    exit 1
fi

ACTUAL_SUM="$(sha256_file "$ASSET_FILE")"
if [ "$EXPECTED_SUM" != "$ACTUAL_SUM" ]; then
    log_err "Checksum verification failed for $ASSET_NAME"
    exit 1
fi
log_ok "Checksum verified"

EXTRACT_DIR="$TMP_DIR/extract"
mkdir -p "$EXTRACT_DIR"

if [ "$EXT" = "zip" ]; then
    require_cmd unzip
    unzip -q "$ASSET_FILE" -d "$EXTRACT_DIR"
else
    require_cmd tar
    tar -xzf "$ASSET_FILE" -C "$EXTRACT_DIR"
fi

BIN_NAME="gwtt"
if [ "$OS" = "windows" ]; then
    BIN_NAME="gwtt.exe"
fi

if [ ! -f "$EXTRACT_DIR/$BIN_NAME" ]; then
    log_err "Binary not found in archive: $BIN_NAME"
    exit 1
fi

cp "$EXTRACT_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
if [ "$OS" != "windows" ]; then
    chmod +x "$INSTALL_DIR/$BIN_NAME"
fi
log_ok "Installed $BIN_NAME to $INSTALL_DIR"

echo ""
log_ok "Installation complete"
echo "Binary: $INSTALL_DIR/$BIN_NAME"
