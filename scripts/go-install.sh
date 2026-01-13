#!/bin/bash
# go-install.sh - Install git-worktree-tasks and gwtt binaries to $GOPATH/bin
#
# This script is for Go developers who have Go installed.
# It builds both the git-worktree-tasks and gwtt binaries and installs them to $GOPATH/bin.
#
# Usage: ./scripts/go-install.sh [install_path]
#   install_path: Optional. Defaults to $(go env GOBIN) or $GOPATH/bin
#
# Example:
#   ./scripts/go-install.sh                    # Install to default location
#   ./scripts/go-install.sh ~/.local/bin       # Install to custom location

set -e

# Color output helpers
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    echo "This script requires Go to be installed. Visit: https://golang.org/doc/install"
    exit 1
fi

# Determine installation directory
if [ -n "$1" ]; then
    INSTALL_DIR="$1"
else
    # Try to use GOBIN first, fall back to GOPATH/bin
    INSTALL_DIR=$(go env GOBIN)
    if [ -z "$INSTALL_DIR" ]; then
        GOPATH=$(go env GOPATH)
        if [ -z "$GOPATH" ]; then
            echo -e "${RED}✗ Could not determine GOPATH${NC}"
            echo "Please set GOPATH or provide an installation directory as an argument"
            exit 1
        fi
        INSTALL_DIR="$GOPATH/bin"
    fi
fi

# Create installation directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}Creating installation directory: $INSTALL_DIR${NC}"
    mkdir -p "$INSTALL_DIR"
fi

# Check if installation directory is writable
if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${RED}✗ Installation directory is not writable: $INSTALL_DIR${NC}"
    exit 1
fi

echo -e "${YELLOW}Building binaries...${NC}"

# Build git-worktree-tasks
echo "  Building git-worktree-tasks..."
if go build -o "$INSTALL_DIR/git-worktree-tasks" ./; then
    echo -e "  ${GREEN}✓${NC} git-worktree-tasks built successfully"
else
    echo -e "${RED}✗ Failed to build git-worktree-tasks${NC}"
    exit 1
fi

# Build gwtt
echo "  Building gwtt..."
if go build -o "$INSTALL_DIR/gwtt" ./; then
    echo -e "  ${GREEN}✓${NC} gwtt built successfully"
else
    echo -e "${RED}✗ Failed to build gwtt${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "Both binaries have been installed to: $INSTALL_DIR"
echo ""
echo "Next steps:"
echo "1. Verify $INSTALL_DIR is in your \$PATH:"
echo "   echo \$PATH | grep $INSTALL_DIR"
echo ""
echo "2. Choose how to enable the 'gwtt' shorthand:"
echo ""
echo "   Option A: Shell Alias (Recommended)"
echo "   Add this to your shell config (~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish):"
echo ""
echo "   Bash/Zsh:"
echo "     alias gwtt=\"git-worktree-tasks\""
echo ""
echo "   Fish:"
echo "     alias gwtt git-worktree-tasks"
echo ""
echo "   Option B: Manual Symlink"
echo "     ln -s $INSTALL_DIR/git-worktree-tasks $INSTALL_DIR/gwtt"
echo ""
echo "   Option C: Use full command name"
echo "     Just use 'git-worktree-tasks' directly (no alias needed)"
echo ""
echo "Test your installation:"
echo "   git-worktree-tasks --version"
echo "   gwtt --version  (if you set up an alias or symlink)"
