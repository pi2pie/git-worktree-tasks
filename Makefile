.PHONY: help build go-build go-install go-uninstall clean

# Build directory
DIST_DIR := dist

# Default target
help:
	@echo "git-worktree-tasks — Build and installation targets"
	@echo ""
	@echo "Local Development:"
	@echo "  make build           Build both git-worktree-tasks and gwtt binaries to $(DIST_DIR)/"
	@echo "  make clean           Remove binaries from $(DIST_DIR)/"
	@echo ""
	@echo "Go Developer Installation (requires Go):"
	@echo "  make go-install      Build and install both binaries to \$$GOPATH/bin"
	@echo "  make go-uninstall    Remove both binaries from \$$GOPATH/bin"
	@echo ""
	@echo "Development:"
	@echo "  make help            Show this help message"
	@echo ""

# Build both binaries to dist/ directory
.PHONY: build
build: $(DIST_DIR)
	@echo "Building git-worktree-tasks..."
	go build -o $(DIST_DIR)/git-worktree-tasks ./
	@echo "✓ Built git-worktree-tasks"
	@echo ""
	@echo "Building gwtt..."
	go build -o $(DIST_DIR)/gwtt ./
	@echo "✓ Built gwtt"
	@echo ""
	@echo "Both binaries are ready in $(DIST_DIR)/"

# Create dist directory if it doesn't exist
$(DIST_DIR):
	@mkdir -p $(DIST_DIR)

# Alias for build
.PHONY: go-build
go-build: build

# Install binaries to $GOPATH/bin (requires Go)
.PHONY: go-install
go-install:
	@bash ./scripts/go-install.sh

# Uninstall binaries from $GOPATH/bin
.PHONY: go-uninstall
go-uninstall:
	@bash ./scripts/go-uninstall.sh

# Clean up binaries in dist/ directory
.PHONY: clean
clean:
	@echo "Cleaning up binaries from $(DIST_DIR)/..."
	@rm -rf $(DIST_DIR)
	@echo "✓ Cleaned"

# Default when just running 'make'
.DEFAULT_GOAL := help
