# MikroTik Configuration Backup

[![CI](https://github.com/mountain-reverie/mikrotik-configuation-backup/actions/workflows/ci.yml/badge.svg)](https://github.com/mountain-reverie/mikrotik-configuation-backup/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mountain-reverie/mikrotik-configuation-backup)](https://goreportcard.com/report/github.com/mountain-reverie/mikrotik-configuation-backup)
[![codecov](https://codecov.io/gh/mountain-reverie/mikrotik-configuation-backup/branch/main/graph/badge.svg)](https://codecov.io/gh/mountain-reverie/mikrotik-configuation-backup)
[![License](https://img.shields.io/github/license/mountain-reverie/mikrotik-configuation-backup)](LICENSE)
[![Release](https://img.shields.io/github/v/release/mountain-reverie/mikrotik-configuation-backup)](https://github.com/mountain-reverie/mikrotik-configuation-backup/releases)

A robust CLI tool to backup MikroTik RouterOS configurations. This tool connects to MikroTik devices via SSH and exports their configurations to local files for version control and disaster recovery.

## Features

- 🔒 **Secure SSH Connection** - Supports both password and key-based authentication
- 📝 **Full Configuration Export** - Exports complete RouterOS configuration
- 🚀 **Fast & Lightweight** - Written in Go for optimal performance
- 🔄 **CI/CD Ready** - Perfect for automated backup workflows
- 🧪 **Well Tested** - Comprehensive unit and integration tests
- 📦 **Easy Installation** - Multiple installation methods available

## Installation

### Using Go Install

```bash
go install github.com/mountain-reverie/mikrotik-configuation-backup/cmd/mikrotik-backup@latest
```

### Download Binary

Download the latest binary from the [releases page](https://github.com/mountain-reverie/mikrotik-configuation-backup/releases).

### Build from Source

```bash
git clone https://github.com/mountain-reverie/mikrotik-configuation-backup.git
cd mikrotik-configuation-backup
go build -o mikrotik-backup ./cmd/mikrotik-backup
```

## Usage

### Basic Usage

Backup configuration using password authentication:

```bash
mikrotik-backup backup --host 192.168.88.1 --username admin --password mypassword --output backup.rsc
```

Backup configuration using SSH key:

```bash
mikrotik-backup backup --host 192.168.88.1 --username admin --key ~/.ssh/mikrotik_rsa --output backup.rsc
```

### Environment Variables

All flags can be set via environment variables:

```bash
export MIKROTIK_HOST=192.168.88.1
export MIKROTIK_USERNAME=admin
export MIKROTIK_PASSWORD=mypassword
mikrotik-backup backup --output backup.rsc
```

### Command Line Options

```
mikrotik-backup backup [flags]

Flags:
  -H, --host string       MikroTik device hostname or IP (required) [env: MIKROTIK_HOST]
  -p, --port int          SSH port (default: 22) [env: MIKROTIK_PORT]
  -u, --username string   SSH username (default: "admin") [env: MIKROTIK_USERNAME]
  -P, --password string   SSH password [env: MIKROTIK_PASSWORD]
  -k, --key string        Path to SSH private key file [env: MIKROTIK_KEY_FILE]
  -o, --output string     Output file path (default: "backup.rsc")
  -h, --help             Help for backup
```

### Version Information

```bash
mikrotik-backup version
```

## Development

This project follows Go best practices and uses standard Go tooling.

### Prerequisites

- **Go 1.22+** - [Installation guide](https://go.dev/doc/install)
- **golangci-lint** - [Installation guide](https://golangci-lint.run/usage/install/)

### Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/mountain-reverie/mikrotik-configuation-backup.git
   cd mikrotik-configuation-backup
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test -v -race ./...
   ```

4. Build the binary:
   ```bash
   go build -o mikrotik-backup ./cmd/mikrotik-backup
   ```

5. Run linting:
   ```bash
   golangci-lint run
   ```

### Common Commands

**Build:**
```bash
# Build for current platform
go build -o mikrotik-backup ./cmd/mikrotik-backup

# Build with version information
go build -ldflags="-X main.version=v1.0.0 -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o mikrotik-backup ./cmd/mikrotik-backup

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o mikrotik-backup-linux-amd64 ./cmd/mikrotik-backup
```

**Testing:**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run integration tests
go test -tags=integration ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

**Code Quality:**
```bash
# Format code
go fmt ./...

# Run golangci-lint
golangci-lint run

# Run golangci-lint with auto-fix
golangci-lint run --fix

# Run go vet
go vet ./...

# Run staticcheck
staticcheck ./...

# Run security scanner
gosec ./...
```

**Dependencies:**
```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy dependencies (remove unused, add missing)
go mod tidy

# Update all dependencies
go get -u ./...
go mod tidy

# Update specific dependency
go get -u github.com/urfave/cli/v2
```

### Project Structure

```
.
├── cmd/
│   └── mikrotik-backup/     # Main application entry point
│       └── main.go
├── internal/
│   ├── backup/              # Backup service implementation
│   │   ├── backup.go
│   │   ├── backup_test.go
│   │   └── backup_integration_test.go  # Integration tests (use -tags=integration)
│   ├── config/              # Configuration management
│   └── ssh/                 # SSH client implementation
├── .github/
│   └── workflows/           # GitHub Actions CI/CD
│       ├── ci.yml
│       └── release.yml
├── .golangci.yml           # golangci-lint configuration
├── .goreleaser.yml         # GoReleaser configuration
├── .pre-commit-config.yaml # Pre-commit hooks (optional)
└── README.md
```

### Running Tests

**Unit tests:**
```bash
go test -v -race ./...
```

**Integration tests:**
```bash
# Set environment variables
export MIKROTIK_HOST=192.168.88.1
export MIKROTIK_USERNAME=admin
export MIKROTIK_PASSWORD=yourpassword

# Run integration tests
go test -v -tags=integration ./...
```

**All tests with coverage:**
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Code Quality Standards

This project uses comprehensive linting with golangci-lint:

- **60+ linters enabled** - Comprehensive code quality checks
- **gofumpt** - Stricter formatting than gofmt
- **gosec** - Security vulnerability scanning
- **staticcheck** - Advanced static analysis
- **See .golangci.yml** for complete configuration

Run quality checks:
```bash
golangci-lint run
```

### Pre-commit Hooks (Optional)

Install pre-commit hooks to automatically run checks before committing:

1. Install pre-commit: `pip install pre-commit`
2. Install hooks: `pre-commit install`
3. Test: `pre-commit run --all-files`

### Continuous Integration

The project uses GitHub Actions for CI/CD:

- **Lint** - Runs golangci-lint
- **Test** - Runs unit tests on Go 1.22 and 1.23
- **Integration Test** - Runs integration tests
- **Build** - Builds binaries for Linux, macOS, and Windows
- **Security** - Runs gosec security scanner
- **Dependency Review** - Checks for vulnerable dependencies

### Release Process

Releases are automated using GoReleaser:

1. Create and push a version tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Create GitHub release with changelog
   - Upload artifacts

Test release locally:
```bash
goreleaser release --snapshot --clean --skip=publish
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting (`go test ./... && golangci-lint run`)
5. Commit your changes using conventional commits
6. Push to your fork (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Best Practices

This project follows Go best practices:

- ✅ **Go 1.22+** - Uses latest stable Go version
- ✅ **Standard Project Layout** - Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- ✅ **Standard Go Commands** - No Make required, use `go build`, `go test`, etc.
- ✅ **Dependency Injection** - Interfaces for testability
- ✅ **Error Wrapping** - Uses `fmt.Errorf` with `%w` for error chains
- ✅ **Context Propagation** - Proper context usage for cancellation
- ✅ **Table-Driven Tests** - Comprehensive test coverage
- ✅ **Parallel Tests** - Tests run in parallel where possible
- ✅ **Mocking** - Interface-based mocking for unit tests
- ✅ **Build Tags** - Separates integration tests with build tags
- ✅ **urfave/cli** - Modern CLI framework with environment variable support

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [urfave/cli Documentation](https://cli.urfave.org/)
- [MikroTik RouterOS Documentation](https://help.mikrotik.com/docs/)

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Support

- 📖 [Documentation](https://github.com/mountain-reverie/mikrotik-configuation-backup/wiki)
- 🐛 [Issue Tracker](https://github.com/mountain-reverie/mikrotik-configuation-backup/issues)
- 💬 [Discussions](https://github.com/mountain-reverie/mikrotik-configuation-backup/discussions)

## Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) - CLI framework
- [golangci-lint](https://golangci-lint.run/) - Comprehensive Go linting
- MikroTik community for RouterOS documentation

---

Made with ❤️ for the MikroTik community
