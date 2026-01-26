#!/usr/bin/env bash
# uninstall.sh - Remove gwtt installed from release assets
#
# Usage: ./scripts/uninstall.sh [install_dir]
#   install_dir: Optional. Defaults to current directory

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

INSTALL_DIR="${1:-}"

log_info() { echo -e "${BLUE}ℹ${NC} $1"; }
log_ok() { echo -e "${GREEN}✓${NC} $1"; }
log_warn() { echo -e "${YELLOW}⚠${NC} $1"; }
log_err() { echo -e "${RED}✗${NC} $1"; }

UNAME_S="$(uname -s 2>/dev/null || echo unknown)"
case "$UNAME_S" in
    Darwin) OS="darwin" ;;
    Linux) OS="linux" ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT) OS="windows" ;;
    *)
        log_err "Unsupported OS: $UNAME_S"
        exit 1
        ;;
esac

if [ -z "$INSTALL_DIR" ]; then
    INSTALL_DIR="$PWD"
fi

if [ ! -d "$INSTALL_DIR" ]; then
    log_err "Install directory does not exist: $INSTALL_DIR"
    exit 1
fi

BIN_NAME="gwtt"
if [ "$OS" = "windows" ]; then
    BIN_NAME="gwtt.exe"
fi

BIN_PATH="$INSTALL_DIR/$BIN_NAME"
if [ -f "$BIN_PATH" ]; then
    rm "$BIN_PATH"
    log_ok "Removed $BIN_PATH"
else
    log_warn "Binary not found: $BIN_PATH"
fi

log_ok "Uninstallation complete"
