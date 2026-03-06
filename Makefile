# Variables
BINARY_NAME=chest
MAIN_BRANCH=main
# Gets the version from Git or defaults to v0.0.0
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0")

.PHONY: all ci quality test build run clean create-feature close-feature release help

all: quality test build

ci: quality test build
	@echo "✅ CI completed successfully!"

quality:
	@echo "[QUALITY] --- Running linter bundle ---"
	golangci-lint run ./...
	@echo "[QUALITY] --- Checking for dependency vulnerabilities ---"
	govulncheck ./...

test:
	@echo "[TEST] --- Running unit tests with Race Detector ---"
	go test -v -race -cover ./...

build:
	@echo "[BUILD] --- Compiling binary (Version: $(VERSION)) ---"
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY_NAME) ./cmd/chest

run: build
	@echo "[RUN] --- Starting $(BINARY_NAME) ---"
	./$(BINARY_NAME)

clean:
	@echo "[CLEAN] --- Removing binaries ---"
	rm -f $(BINARY_NAME)

# --- TRUNK-BASED WORKFLOW COMMANDS ---

## create-feature: Setup a new branch for a feature (usage: make create-feature NAME=vault-logic)
create-feature:
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required (e.g., make create-feature NAME=foo)"; exit 1; fi
	git checkout $(MAIN_BRANCH)
	git pull origin $(MAIN_BRANCH)
	git checkout -b feat/$(NAME)
	@echo "🚀 Feature branch 'feat/$(NAME)' created."

## close-feature: Merge current feature into main and delete branch
close-feature:
	@CURRENT_BRANCH=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$CURRENT_BRANCH" = "$(MAIN_BRANCH)" ]; then echo "Error: You are already on $(MAIN_BRANCH)!"; exit 1; fi; \
	echo "📦 Merging feature $$CURRENT_BRANCH into $(MAIN_BRANCH)..."; \
	git checkout $(MAIN_BRANCH); \
	git pull origin $(MAIN_BRANCH); \
	git merge --no-ff $$CURRENT_BRANCH -m "Merge feature: $$CURRENT_BRANCH"; \
	git branch -d $$CURRENT_BRANCH; \
	@echo "✅ Feature successfully merged and branch removed."

## release: Tag and push a new version (usage: make release VERSION=v1.0.0)
release:
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION is required (e.g., make release VERSION=v1.0.0)"; exit 1; fi
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "$(MAIN_BRANCH)" ]; then echo "Error: Releases must be made from $(MAIN_BRANCH)!"; exit 1; fi
	git pull origin $(MAIN_BRANCH)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "🚀 Version $(VERSION) tagged and pushed!"

help:
	@echo "Chest Makefile - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'