# Kubernetes Resource Comparison Tool

A CLI tool for comparing Kubernetes resources between two clusters with modern terminal UI, Google Cloud authentication support, and HTML report generation.

## Features

- 🔍 **Interactive Cluster Selection** - Choose from available kubectl contexts
- 🏠 **Multi-Namespace Support** - Select specific namespaces or compare all
- 📦 **Prioritized Resource Types** - Common resources (pods, services, deployments) shown first
- 🅰️🅱️ **Universal Cluster Labels** - Generic cluster A/B naming for any Kubernetes environment
- ☁️ **Google Cloud Integration** - Automatic GKE authentication with `gcloud auth login`
- 🎨 **Modern Terminal UI** - Beautiful forms with the [Charm](https://charm.sh/) `huh` library
- 📊 **HTML Report Generation** - Automatic HTML reports with rich visualizations
- ⌨️ **Keyboard Shortcuts** - Use arrow keys, space, Ctrl+A, Enter, and Esc for navigation

## Code Structure

The codebase is organized into modular files in the `src/` directory for better maintainability:

- **`src/main.go`** - Entry point and CLI setup
- **`src/types.go`** - Type definitions (ClusterConfig, ComparisonConfig)
- **`src/setup.go`** - Interactive setup and UI functions
- **`src/auth.go`** - Google Cloud authentication handling
- **`src/kubernetes.go`** - Kubernetes client and resource fetching
- **`src/fetcher.go`** - Resource fetching orchestration
- **`src/output.go`** - JSON and HTML file generation
- **`src/html_template.go`** - HTML template and JavaScript functions
- **`src/utils.go`** - Utility helper functions

## Installation

### Prerequisites

- Go 1.19 or later
- kubectl configured with access to your Kubernetes clusters
- For Google Cloud clusters: `gcloud` CLI installed and configured

### Build from Source

```bash
# Clone or download the repository
git clone <repository-url>
cd operator-comparison

# Install dependencies
go mod tidy

# Build the binary (choose one method):

# Method 1: Using Makefile (recommended)
make build

# Method 2: Using go build with path
go build -o k8s-compare ./src

# Method 3: Build from src directory
cd src && go build -o ../k8s-compare *.go && cd ..
```

### Available Make Targets

```bash
make build        # Build the binary
make test         # Build and test the binary
make test-unit    # Run unit tests
make test-ginkgo  # Run Ginkgo tests
make test-coverage# Run tests with coverage report
make test-watch   # Run tests in watch mode
make clean        # Remove build artifacts
make deps         # Install/update dependencies
make help         # Show available targets
```

## Testing

The project includes comprehensive unit tests using Ginkgo and Gomega testing frameworks.

### Running Tests

**Basic test (build and CLI verification):**
```bash
make test
```

**Unit tests with Go's built-in test runner:**
```bash
make test-unit
```

**Unit tests with Ginkgo (recommended):**
```bash
make test-ginkgo
```

**Tests with coverage report:**
```bash
make test-coverage
# Opens coverage.html in src/ directory
```

**Watch mode for development:**
```bash
make test-watch
```

### Test Structure

The test suite covers all major components:

- **`utils_test.go`** - Tests for utility functions (`contains`, `removeFromSlice`)
- **`types_test.go`** - Tests for data structures (`ClusterConfig`, `ComparisonConfig`)
- **`auth_test.go`** - Tests for Google Cloud authentication logic
- **`output_test.go`** - Tests for JSON/HTML report generation
- **`fetcher_test.go`** - Tests for resource fetching orchestration
- **`kubernetes_test.go`** - Tests for Kubernetes client and resource processing
- **`setup_test.go`** - Tests for interactive setup and resource prioritization
- **`html_template_test.go`** - Tests for HTML template generation and validation
- **`main_test.go`** - Test suite bootstrap and configuration

### Test Categories

- **Unit Tests**: Test individual functions and data structures
- **Integration Tests**: Marked as `PIt()` (pending) - require actual Kubernetes clusters or mocking
- **File I/O Tests**: Use temporary directories for safe testing

### Test Coverage

Current test coverage includes:
- ✅ **97 passing unit tests** covering:
  - Authentication logic (`isGoogleCloudContext`, `ensureGCloudAuth`)
  - Resource prioritization (`reorderResourcesByPriority`)
  - Data structures (`ClusterConfig`, `ComparisonConfig`)
  - HTML template generation (`generateHTMLTemplate`)
  - File operations and output generation
  - Utility functions (`contains`, `removeFromSlice`)
  - Resource fetching configuration and validation
- ⏸️ **42 pending integration tests** (require external dependencies like Kubernetes clusters)
- 🎯 Comprehensive coverage of testable business logic and data processing
- 📊 HTML coverage reports available via `make test-coverage`

## Usage

### Basic Usage

```bash
./k8s-compare
```

The tool will guide you through:

1. **📍 Context Selection** - Choose two Kubernetes contexts
2. **🔐 Authentication** - Automatic Google Cloud auth if needed
3. **🏠 Namespace Selection** - Pick namespaces for each cluster  
4. **📦 Resource Types** - Select which resource types to compare

### Command Line Options

```bash
./k8s-compare --help
```

**Flags:**
- `-o, --output-dir` - Output directory for generated files (default: current directory)
- `-i, --interactive` - Run in interactive mode (default: true)

### Output Files

The tool generates several files:

- **cluster-a.json** - Resources from the first cluster
- **cluster-b.json** - Resources from the second cluster  
- **comparison-report-YYYYMMDD-HHMMSS.html** - Interactive HTML comparison report
- **index.html** - Generic HTML comparison tool (for manual file uploads)

## Google Cloud Authentication

The tool automatically detects Google Cloud contexts (GKE, Connect Gateway) and handles authentication:

- **Auto-detection** - Recognizes `gke_`, `connectgateway_` context patterns
- **Auth verification** - Checks current `gcloud auth list` status
- **Interactive login** - Prompts for `gcloud auth login` when needed
- **Timeout handling** - 2-minute timeout for authentication processes
- **Error recovery** - Clear error messages with retry options

### Troubleshooting Authentication

If you encounter authentication issues:

1. **Manual login**: Run `gcloud auth login` before using the tool
2. **Check contexts**: Verify contexts with `kubectl config get-contexts`
3. **Cluster access**: Ensure you have proper RBAC permissions
4. **Network connectivity**: Verify cluster endpoint accessibility

## HTML Report Features

The generated HTML reports include:

### 📊 **Metadata Overview**
- Cluster contexts and namespaces
- Resource counts and types
- Generation timestamp
- Resource type tags

### 📋 **Three-Tab Interface**
- **Overview** - Summary statistics and counts
- **Resource Breakdown** - Resources by type for each cluster
- **Detailed Comparison** - Side-by-side resource differences

### 🔍 **Rich Comparison Features**
- **Individual Resource Collapsing** - Expand/collapse each resource
- **Syntax Highlighted JSON** - Color-coded keys, strings, numbers, booleans
- **Status Badges** - Visual indicators for different/unique resources
- **Smart Field Filtering** - Skips ephemeral fields like `resourceVersion`
- **Horizontal Scrolling** - Handle wide JSON content
- **No Truncation** - Complete data visibility

### 🎨 **Modern UI Design**
- Clean, responsive design
- Emoji icons for visual clarity
- Color-coded differences
- Smooth animations and transitions

## Keyboard Shortcuts

During selection prompts:
- **Arrow Keys** - Navigate options
- **Space** - Toggle selection (multi-select)
- **Ctrl+A** - Select all items
- **Enter** - Confirm selection
- **Esc** - Cancel/go back

## Resource Type Priority

Common Kubernetes resources are prioritized in the selection list:

**High Priority:** pods, services, deployments, configmaps, secrets, ingresses, persistentvolumes, persistentvolumeclaims, nodes

**Standard Priority:** All other available resource types

## Example Workflow

```bash
$ ./k8s-compare

🔍 Kubernetes Cluster Resource Comparison Tool
==============================================

📍 Step 1: Select Kubernetes contexts
? Select Cluster A context: › gke_my-project_us-central1_cluster-prod
? Select Cluster B context: › gke_my-project_us-central1_cluster-staging

🔐 Checking authentication for selected contexts...
✓ Authentication verified for Cluster A
✓ Authentication verified for Cluster B

🏠 Step 2: Select namespaces
? Select namespaces for Cluster A: ›
  ◉ default
  ◉ kube-system  
  ◯ istio-system

🏠 Step 2: Select namespaces  
? Select namespaces for Cluster B: ›
  ◉ default
  ◉ kube-system
  ◯ istio-system

📦 Step 3: Select resource types
? Select resource types to compare: ›
  ◉ pods
  ◉ services
  ◉ deployments
  ◯ configmaps

🚀 Fetching resources from Cluster A...
✓ Found 45 resources

🚀 Fetching resources from Cluster B...  
✓ Found 38 resources

📝 Generating output files...
✓ Generated cluster-a.json
✓ Generated cluster-b.json
✓ Generated comparison-report-20241215-143022.html

🎉 Comparison completed successfully!
📄 Generated files:
   - cluster-a.json
   - cluster-b.json
   - comparison-report-20241215-143022.html
💡 Open the HTML report in your browser to view the comparison
```

## Technical Details

- **Built with Go 1.21+** using Kubernetes client-go libraries
- **Modern Terminal UI** powered by [Charm's huh](https://github.com/charmbracelet/huh) library
- **Dynamic Resource Discovery** via Kubernetes API discovery client
- **Cross-platform** compatibility (Linux, macOS, Windows)
- **Memory efficient** streaming of large resource sets

## Contributing

Contributions welcome! Please feel free to submit issues and pull requests.

## License

MIT License - see LICENSE file for details. 