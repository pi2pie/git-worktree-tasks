#!/bin/bash

# go-uninstall.sh — Remove git-worktree-tasks and gwtt binaries
#
# This script removes both binaries that were installed via go-install.sh
#
# Requirements:
#   - Go 1.25.5 or higher (to determine $GOPATH/bin)
#
# Usage:
#   ./scripts/go-uninstall.sh [install-directory]
#   MAN_DIR=~/.local/share/man/man1 ./scripts/go-uninstall.sh
#
# Examples:
#   ./scripts/go-uninstall.sh                    # Remove from $GOPATH/bin
#   ./scripts/go-uninstall.sh /usr/local/bin     # Remove from custom directory

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_error() {
    echo -e "${RED}✗ Error:${NC} $1" >&2
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in your PATH"
    exit 1
fi

# Determine install directory
if [ -n "$1" ]; then
    INSTALL_DIR="$1"
else
    INSTALL_DIR=$(go env GOBIN)
    if [ -z "$INSTALL_DIR" ]; then
        INSTALL_DIR="$(go env GOPATH)/bin"
    fi
fi

# Determine man page directory (override with MAN_DIR)
if [ -n "$MAN_DIR" ]; then
    MAN_INSTALL_DIR="$MAN_DIR"
else
    INSTALL_PREFIX=$(dirname "$INSTALL_DIR")
    MAN_INSTALL_DIR="$INSTALL_PREFIX/share/man/man1"
fi

print_info "Uninstall directory: $INSTALL_DIR"
print_info "Man directory: $MAN_INSTALL_DIR"

# Check if directory exists
if [ ! -d "$INSTALL_DIR" ]; then
    print_error "Directory does not exist: $INSTALL_DIR"
    exit 1
fi

# Remove git-worktree-tasks binary
if [ -f "$INSTALL_DIR/git-worktree-tasks" ]; then
    echo ""
    print_info "Removing git-worktree-tasks..."
    if rm "$INSTALL_DIR/git-worktree-tasks"; then
        print_success "Removed git-worktree-tasks"
    else
        print_error "Failed to remove git-worktree-tasks"
        exit 1
    fi
else
    print_warning "git-worktree-tasks not found in $INSTALL_DIR"
fi

# Remove gwtt binary
if [ -f "$INSTALL_DIR/gwtt" ]; then
    echo ""
    print_info "Removing gwtt..."
    if rm "$INSTALL_DIR/gwtt"; then
        print_success "Removed gwtt"
    else
        print_error "Failed to remove gwtt"
        exit 1
    fi
else
    print_warning "gwtt not found in $INSTALL_DIR"
fi

# Remove man pages
if [ -d "$MAN_INSTALL_DIR" ]; then
    MAN_GLOB_ROOT="$MAN_INSTALL_DIR/git-worktree-tasks"*.1
    MAN_GLOB_ALIAS="$MAN_INSTALL_DIR/gwtt"*.1
    if ls $MAN_GLOB_ROOT $MAN_GLOB_ALIAS &> /dev/null; then
        echo ""
        print_info "Removing man pages..."
        if rm $MAN_GLOB_ROOT $MAN_GLOB_ALIAS; then
            print_success "Removed man pages"
        else
            print_error "Failed to remove man pages"
            exit 1
        fi
    else
        print_warning "Man pages not found in $MAN_INSTALL_DIR"
    fi
else
    print_warning "Man directory not found: $MAN_INSTALL_DIR"
fi

# Print cleanup instructions
echo ""
echo "========================================"
print_success "Uninstallation complete!"
echo "========================================"
echo ""
echo "Binaries removed from: $INSTALL_DIR"
echo "Man pages removed from: $MAN_INSTALL_DIR"
echo ""
echo "Additional cleanup needed:"
echo ""
echo "If you added a shell alias, remove it:"
echo ""
echo "1. Bash/Zsh users:"
echo "   Edit ~/.bashrc or ~/.zshrc and remove:"
echo "   ${BLUE}alias gwtt=\"git-worktree-tasks\"${NC}"
echo ""
echo "2. Fish users:"
echo "   Edit ~/.config/fish/config.fish and remove:"
echo "   ${BLUE}alias gwtt git-worktree-tasks${NC}"
echo ""
echo "3. If you created a symlink:"
echo "   ${BLUE}rm $INSTALL_DIR/gwtt${NC}"
echo ""
echo "Then reload your shell:"
echo "   ${BLUE}source ~/.bashrc${NC}  (Bash)"
echo "   ${BLUE}source ~/.zshrc${NC}   (Zsh)"
echo "   Restart terminal or ${BLUE}exec fish${NC} (Fish)"
echo ""
