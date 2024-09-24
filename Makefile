# Project-specific variables
BINARY_NAME = go-visualizer
SRC_DIR = ./cmd/visualizer
BUILD_DIR = ./bin
GO_FILES = $(shell find . -name '*.go' -not -path "./vendor/*")
BINARY = $(BUILD_DIR)/$(BINARY_NAME)

# Default target
all: build

# Build the binary only if source files are newer than the binary
$(BINARY): $(GO_FILES)
	@echo "Building the Go visualizer project..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(SRC_DIR)/main.go

# Target for building
build: $(BINARY)

# Run the binary (forces a build first if needed)
run: build
	@echo "Running the Go visualizer..."
	$(BINARY) ./mux/

# Format the Go source files
fmt:
	@echo "Formatting Go files..."
	gofmt -w $(GO_FILES)

# Tidy dependencies and update go.sum
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Clean the build directory
clean:
	@echo "Cleaning build files..."
	rm -rf $(BUILD_DIR)

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Lint the project (requires golangci-lint installed)
lint:
	@echo "Running lint checks..."
	golangci-lint run ./...

# Install golangci-lint
install-linter:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

# Install all dependencies
install:
	@echo "Installing dependencies..."
	go mod download

# Generate graphs (sample task, modify according to your project logic)
graph: $(BINARY)
	@echo "Generating graph visualization..."
	$(BINARY) --generate-graph

.PHONY: build run fmt tidy clean test lint