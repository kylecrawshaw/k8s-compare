.PHONY: build clean test test-unit test-ginkgo test-coverage test-watch deps help

# Default target
build:
	@echo "🔨 Building k8s-compare..."
	go build -o k8s-compare ./src
	@echo "✅ Build complete! Run with: ./k8s-compare"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f k8s-compare
	@echo "✅ Clean complete!"

# Test the build
test: build
	@echo "🧪 Testing binary..."
	./k8s-compare --help

# Run unit tests
test-unit:
	@echo "🧪 Running unit tests..."
	cd src && go test -v ./...

# Run tests with Ginkgo
test-ginkgo:
	@echo "🧪 Running Ginkgo tests..."
	cd src && ginkgo -v

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	cd src && go test -v -coverprofile=coverage.out ./...
	cd src && go tool cover -html=../coverage.out -o ../coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Run tests in watch mode
test-watch:
	@echo "🧪 Running tests in watch mode..."
	cd src && ginkgo watch -v

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

run: build
	./k8s-compare $@

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the k8s-compare binary"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Build and test the binary"
	@echo "  test-unit    - Run unit tests"
	@echo "  test-ginkgo  - Run Ginkgo tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  test-watch   - Run tests in watch mode"
	@echo "  deps         - Install/update dependencies"
	@echo "  run          - Run the k8s-compare binary"
	@echo "  help         - Show this help message" 