# Project settings
BINARY_NAME = hmr
OUTPUT_DIR = ./bin

# Go settings
GO = go
GO_BUILD_FLAGS = -o $(OUTPUT_DIR)/$(BINARY_NAME)
GO_BUILD_STATIC_FLAGS = -o $(OUTPUT_DIR)/$(BINARY_NAME) -ldflags "-s -w"
GO_RUN_FLAGS = ./main.go

.PHONY: help build build-static run clean

# Default target
.DEFAULT_GOAL := help

help: ## Display this help message
	@echo "Usage: make [target] [args]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make run DIR=./my-app"
	@echo "  make build-static"

build: ## Build the application dynamically
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) build $(GO_BUILD_FLAGS) ./main.go
	@echo "Build completed: $(OUTPUT_DIR)/$(BINARY_NAME)"

build-static: ## Build the application statically
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 $(GO) build $(GO_BUILD_STATIC_FLAGS) ./main.go
	@echo "Statically linked build completed: $(OUTPUT_DIR)/$(BINARY_NAME)"
	@echo "You can now move the $(BINARY_NAME) binary to your bin path. Enjoy!"

run: ## Run the application (provide DIR=path/to/app)
ifndef DIR
	$(error DIR is not set. Usage: make run DIR=path/to/app)
endif
	@$(GO) run $(GO_RUN_FLAGS) $(DIR)

clean: ## Clean up build files
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleaned up build files."
