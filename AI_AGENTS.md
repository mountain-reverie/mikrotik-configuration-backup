# AI Agent Documentation

This document provides comprehensive guidance for AI agents (like Claude Code, GitHub Copilot, ChatGPT, etc.) working on this project. It covers project structure, coding standards, workflows, and best practices.

## Table of Contents

- [Project Overview](#project-overview)
- [Project Structure](#project-structure)
- [Technology Stack](#technology-stack)
- [Coding Standards](#coding-standards)
- [Development Workflow](#development-workflow)
- [Testing Guidelines](#testing-guidelines)
- [CI/CD Workflows](#cicd-workflows)
- [Common Tasks](#common-tasks)
- [Security Considerations](#security-considerations)
- [Important Constraints](#important-constraints)

## Project Overview

**Project Name:** MikroTik Configuration Backup
**Type:** CLI Tool
**Language:** Go (latest stable)
**Purpose:** Backup MikroTik RouterOS configurations via SSH

### Key Features

- SSH-based configuration export (password and key authentication)
- CLI built with `github.com/urfave/cli/v2`
- Environment variable support for all flags
- Comprehensive testing (unit, integration, benchmarks)
- Signed releases with cosign (keyless OIDC)
- SBOM generation for supply chain security
- GitHub Pages for documentation and coverage reports

### Design Principles

1. **Simplicity** - Use standard Go tooling, no Makefiles
2. **Testability** - Interface-based dependency injection
3. **Security** - Signed releases, security scanning, no hardcoded secrets
4. **Automation** - Comprehensive CI/CD with GitHub Actions
5. **Observability** - Coverage reports, benchmarks, test analytics

## Project Structure

```
.
├── cmd/
│   └── mikrotik-backup/          # Main CLI entry point
│       └── main.go               # Uses urfave/cli/v2, calls internal packages
├── internal/                     # Private application code
│   ├── backup/                   # Core backup service
│   │   ├── backup.go             # Service implementation
│   │   ├── backup_test.go        # Unit tests (table-driven)
│   │   └── backup_integration_test.go  # Integration tests (//go:build integration)
│   ├── config/                   # Configuration management
│   └── ssh/                      # SSH client implementation
├── .github/
│   └── workflows/                # GitHub Actions workflows
│       ├── README.md             # Detailed workflow documentation
│       ├── ci.yml                # Reusable CI workflow
│       ├── main.yml              # Main branch automation (auto-tag, pages)
│       ├── pr.yml                # PR checks (coverage diff, Dependabot)
│       ├── release.yml           # GoReleaser with signing
│       └── pages.yml             # GitHub Pages deployment
├── .golangci.yml                 # golangci-lint v2 configuration (60+ linters)
├── .goreleaser.yml               # GoReleaser v2 with cosign + SBOM
├── lefthook.yml                  # Git hooks (Go-native, no Python)
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── README.md                     # User documentation
├── CONTRIBUTING.md               # Contributor guide
└── AI_AGENTS.md                  # This file

```

### File Naming Conventions

- `*_test.go` - Unit tests (in same package)
- `*_integration_test.go` - Integration tests (use `//go:build integration`)
- `main.go` - Only in `cmd/` directories
- No uppercase letters in file names (except markdown)

## Technology Stack

### Core Dependencies

- **CLI Framework:** `github.com/urfave/cli/v2` - NOT Cobra
- **SSH Library:** `golang.org/x/crypto/ssh`
- **Testing:** Standard library `testing` package
- **No external mocking libraries** - Use interface-based mocks

### Development Tools

- **Linter:** golangci-lint v2 (60+ linters enabled)
- **Formatter:** gofumpt (stricter than gofmt)
- **Import Organizer:** goimports with local prefix
- **Security Scanner:** gosec
- **Git Hooks:** lefthook (Go-native, NOT pre-commit/Python)
- **Release:** GoReleaser v2
- **Signing:** cosign (keyless with GitHub OIDC)
- **SBOM:** Syft

### CI/CD Tools

- **GitHub Actions:** All automation
- **gotestsum:** Better test output
- **goteststats:** Test analytics
- **benchstat:** Benchmark comparison
- **CodeQL:** Security analysis

## Coding Standards

### Go Version and Style

- **Required:** Latest stable Go release
- **Policy:** This project always uses the latest stable Go version
- **Check:** See `go.mod` for current version
- **CI Testing:** Tests run on both `stable` and `oldstable` Go releases
- **Auto-Update:** Monthly workflow checks for new Go releases and creates PRs automatically
- **Style Guide:** [Effective Go](https://go.dev/doc/effective_go) + [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **Formatting:** gofumpt (stricter than gofmt)

### Package Structure

```go
// Package backup provides functionality for backing up MikroTik RouterOS configurations.
package backup

import (
    "context"
    "fmt"
    "io"

    // Standard library first
    // External packages second
    // Internal packages last with local prefix
    "github.com/mountain-reverie/mikrotik-configuation-backup/internal/ssh"
)
```

### Naming Conventions

```go
// GOOD
type Service struct {}          // Exported type
type sshClient struct {}        // Unexported type
func NewService() *Service {}   // Constructor
func (s *Service) Execute() {}  // Method

// BAD
type BackupService struct {}    // Redundant package name
type SSH_Client struct {}       // Underscores
func new_service() {}           // Underscores, unexported constructor
```

### Error Handling

```go
// ALWAYS wrap errors with context
if err := client.Connect(ctx, config); err != nil {
    return fmt.Errorf("failed to connect to device: %w", err)
}

// NOT this
if err := client.Connect(ctx, config); err != nil {
    return err  // ❌ No context
}

// NOT this
if err := client.Connect(ctx, config); err != nil {
    return fmt.Errorf("error: %v", err)  // ❌ Use %w, not %v
}
```

### Context Usage

```go
// ALWAYS propagate context
func (s *Service) Execute(ctx context.Context, config Config) error {
    // Pass context to all I/O operations
    if err := s.client.Connect(ctx, config); err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }

    result, err := s.client.ExecuteCommand(ctx, "/export")
    // ...
}
```

### Interface Design

```go
// Use interfaces for dependency injection
type SSHClient interface {
    Connect(ctx context.Context, config Config) error
    ExecuteCommand(ctx context.Context, cmd string) (string, error)
    Close() error
}

// Service depends on interface, not concrete type
type Service struct {
    client SSHClient  // ✅ Testable
}
```

### Version Information

```go
// ❌ DO NOT use global variables for version info
var (
    version = "dev"
    commit  = "none"
)

// ✅ DO use runtime/debug.ReadBuildInfo()
func getBuildInfo() BuildInfo {
    info := BuildInfo{
        Version:   "dev",
        Commit:    "none",
        Date:      "unknown",
        GoVersion: "unknown",
    }

    if buildInfo, ok := debug.ReadBuildInfo(); ok {
        info.GoVersion = buildInfo.GoVersion
        for _, setting := range buildInfo.Settings {
            switch setting.Key {
            case "vcs.revision":
                info.Commit = setting.Value
            case "vcs.time":
                info.Date = setting.Value
            }
        }
    }

    return info
}
```

### Constants vs Magic Numbers

```go
// ❌ BAD - Magic numbers
port := 22
timeout := 30

// ✅ GOOD - Named constants
const (
    defaultSSHPort = 22
    defaultTimeout = 30 * time.Second
)
```

## Development Workflow

### 1. Making Changes

```bash
# Create feature branch
git checkout -b feature/my-feature

# Make changes
# Edit files...

# Run tests
go test -v -race ./...

# Run linter
golangci-lint run

# Commit (conventional commits format)
git commit -m "feat: add new feature"
```

### 2. Running Tests

```bash
# Unit tests only
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Integration tests (requires environment variables)
export MIKROTIK_HOST=192.168.88.1
export MIKROTIK_USERNAME=admin
export MIKROTIK_PASSWORD=password
go test -tags=integration ./...

# Benchmarks
go test -bench=. -benchmem ./...

# Short tests (for quick checks)
go test -short ./...
```

### 3. Building

```bash
# Standard build
go build ./cmd/mikrotik-backup

# With build info
go build -ldflags="-s -w" ./cmd/mikrotik-backup

# Test GoReleaser locally
goreleaser release --snapshot --clean --skip=publish
```

### 4. Local Git Hooks (Optional)

```bash
# Install lefthook
go install github.com/evilmartians/lefthook@latest

# Install tools
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install mvdan.cc/gofumpt@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Install hooks
lefthook install

# Now pre-commit/pre-push hooks run automatically
```

## Testing Guidelines

### Unit Test Structure

Use **table-driven tests** with `t.Parallel()`:

```go
func TestService_Execute(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name        string
        setupMock   func() *mockSSHClient
        wantErr     bool
        errContains string
    }{
        {
            name: "success",
            setupMock: func() *mockSSHClient {
                return &mockSSHClient{
                    executeCommandFunc: func(_ context.Context, cmd string) (string, error) {
                        return "config output", nil
                    },
                }
            },
            wantErr: false,
        },
        {
            name: "connection error",
            setupMock: func() *mockSSHClient {
                return &mockSSHClient{
                    connectFunc: func(_ context.Context, _ Config) error {
                        return errors.New("connection failed")
                    },
                }
            },
            wantErr:     true,
            errContains: "failed to connect",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            client := tt.setupMock()
            service := backup.New(client)
            output := &bytes.Buffer{}

            err := service.Execute(context.Background(), Config{}, output)

            if (err != nil) != tt.wantErr {
                t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
            }

            if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
                t.Errorf("Execute() error = %v, should contain %q", err, tt.errContains)
            }
        })
    }
}
```

### Mock Implementation

```go
// In *_test.go file
type mockSSHClient struct {
    connectFunc        func(context.Context, Config) error
    executeCommandFunc func(context.Context, string) (string, error)
    closeFunc          func() error
}

func (m *mockSSHClient) Connect(ctx context.Context, cfg Config) error {
    if m.connectFunc != nil {
        return m.connectFunc(ctx, cfg)
    }
    return nil
}

func (m *mockSSHClient) ExecuteCommand(ctx context.Context, cmd string) (string, error) {
    if m.executeCommandFunc != nil {
        return m.executeCommandFunc(ctx, cmd)
    }
    return "", nil
}

func (m *mockSSHClient) Close() error {
    if m.closeFunc != nil {
        return m.closeFunc()
    }
    return nil
}
```

### Integration Tests

```go
//go:build integration

package backup_test

import (
    "context"
    "os"
    "testing"
)

func TestBackupIntegration(t *testing.T) {
    host := os.Getenv("MIKROTIK_HOST")
    if host == "" {
        t.Skip("Skipping integration test: MIKROTIK_HOST not set")
    }

    // Integration test implementation
}
```

## CI/CD Workflows

### Workflow Trigger Summary

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `go-version-update.yml` | Monthly (1st) or manual | Check & update Go version |
| `ci.yml` | Reusable (called by others) | Lint, test, build, security |
| `main.yml` | Push to `main` | Auto-tag, release, deploy pages |
| `pr.yml` | Pull requests | CI, coverage diff, Dependabot |
| `release.yml` | Tag push (`v*`) | GoReleaser with signing |
| `pages.yml` | Called by main/pr | Build docs, coverage, benchmarks |

### What Runs in CI

**Every PR and Push:**
- golangci-lint (60+ linters)
- Unit tests (latest stable and previous stable Go releases)
- Integration tests
- Build for all platforms
- gosec security scan
- CodeQL analysis
- Dependency review (PRs only)

**On Main Branch:**
- All CI checks
- Auto-create semantic version tag
- Deploy GitHub Pages

**On Tag Push:**
- Build multi-platform binaries
- Sign with cosign (keyless OIDC)
- Generate SBOM with Syft
- Create GitHub release

### Automation Features

**Dependabot Integration:**
- Auto-merge safe updates (Go modules, GitHub Actions)
- Claude Code auto-fix for failures

**Coverage Tracking:**
- Coverage diff posted on PRs
- HTML reports on GitHub Pages
- SVG badges generated

**Benchmarking:**
- Benchmark comparison with `benchstat`
- Results published to GitHub Pages

See [.github/workflows/README.md](.github/workflows/README.md) for detailed workflow documentation.

## Common Tasks

### Adding a New Command

1. Add command to `cmd/mikrotik-backup/main.go`:

```go
app := &cli.App{
    Commands: []*cli.Command{
        {
            Name:  "backup",
            Usage: "Backup MikroTik configuration",
            // ...
        },
        // Add new command here
        {
            Name:  "restore",
            Usage: "Restore MikroTik configuration",
            Flags: []cli.Flag{
                // Define flags
            },
            Action: func(c *cli.Context) error {
                // Call internal package
                return restore.Execute(c.Context, /* ... */)
            },
        },
    },
}
```

2. Create implementation in `internal/restore/`:
   - `restore.go` - Implementation
   - `restore_test.go` - Unit tests
   - Define interfaces for dependencies

3. Add tests
4. Update documentation

### Adding a New Flag

```go
&cli.StringFlag{
    Name:    "new-flag",
    Usage:   "description of the flag",
    EnvVars: []string{"MIKROTIK_NEW_FLAG"},  // Always add env var support
    Value:   "default",
},
```

### Fixing Linting Issues

```bash
# See what's wrong
golangci-lint run

# Auto-fix what's possible
golangci-lint run --fix

# Common issues:
# - Unused parameters: use _ for intentionally unused params
# - Error wrapping: use %w instead of %v
# - Package comments: add comment starting with "Package name..."
```

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...
go mod tidy

# Update specific dependency
go get -u github.com/urfave/cli/v2
go mod tidy

# Verify dependencies
go mod verify
```

### Creating a Release

**Automatic (Recommended):**
1. Merge PR to `main`
2. GitHub Actions auto-creates tag (patch bump)
3. Tag triggers release workflow
4. Binaries signed and published

**Manual (for major/minor versions):**
```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

## Security Considerations

### What to Avoid

❌ **Never commit secrets**
- No hardcoded passwords, API keys, tokens
- Use environment variables
- `.env` files are gitignored

❌ **No global mutable state**
- No global variables (except constants)
- Pass dependencies explicitly

❌ **No unsafe operations**
- No `unsafe` package without strong justification
- No `#nosec` comments without review

### Security Scanning

All code is scanned by:
- **gosec** - Go security issues
- **CodeQL** - Advanced security analysis
- **golangci-lint** - Security-focused linters
- **Dependency Review** - Vulnerable dependencies

### Signed Releases

All releases are signed with **cosign** using keyless signing:
- GitHub OIDC tokens (no private keys)
- Ephemeral certificates (30-minute lifetime)
- Recorded in Rekor transparency log
- Certificate contains workflow metadata

Users can verify releases:
```bash
cosign verify-blob checksums.txt \
  --certificate checksums.txt.pem \
  --signature checksums.txt.sig \
  --certificate-identity=https://github.com/mountain-reverie/mikrotik-configuation-backup/.github/workflows/release.yml@refs/tags/v1.0.0 \
  --certificate-oidc-issuer=https://token.actions.githubusercontent.com
```

### GitHub Actions Security

**CRITICAL SECURITY REQUIREMENT:** All GitHub Actions MUST be pinned to commit SHAs, not tags.

**Why SHA Pinning?**
- ✅ Prevents supply chain attacks (tag poisoning)
- ✅ Ensures exact version is used (immutable)
- ✅ No silent updates that could introduce vulnerabilities
- ✅ Explicit control over when to update dependencies

**Tag Vulnerability:**
- ❌ Tags can be moved to point to different commits
- ❌ Attacker could compromise an action and republish under existing tag
- ❌ Automated updates could pull malicious code

**Correct Format:**
```yaml
# ✅ GOOD - Pinned to SHA with version comment
- uses: actions/checkout@08eba0b27e820071cde6df949e0beb9ba4906955 # v4.3.0

# ❌ BAD - Using tag (mutable)
- uses: actions/checkout@v4

# ❌ BAD - Using branch (very mutable)
- uses: actions/checkout@main
```

**Documentation:**
- All action SHAs are documented in `.github/ACTION_SHAS.md`
- When updating an action:
  1. Find the release on GitHub
  2. Get the 40-character commit SHA from the release tag
  3. Update all workflow files using the SHA
  4. Update `.github/ACTION_SHAS.md` with version, SHA, and date
  5. Always add a version comment after the SHA for readability

**Example Update Process:**
```bash
# Find the SHA for a release
git ls-remote https://github.com/actions/checkout refs/tags/v4.3.0

# Update workflow files
sed -i 's|uses: actions/checkout@v4|uses: actions/checkout@08eba0b27e820071cde6df949e0beb9ba4906955 # v4.3.0|g' .github/workflows/*.yml

# Update documentation
# Edit .github/ACTION_SHAS.md with new version and SHA
```

## Important Constraints

### What NOT to Do

❌ **Don't use Makefiles**
- Use standard `go build`, `go test`, etc.
- Everything should work with native Go tooling

❌ **Don't use Cobra**
- This project uses `urfave/cli/v2`
- All CLI code should use this framework

❌ **Don't use Python pre-commit**
- Use lefthook (Go-native) instead
- Configuration in `lefthook.yml`

❌ **Don't use external services for coverage**
- No Codecov, Coveralls, etc.
- Use GitHub Pages (self-hosted)

❌ **Don't skip CI checks locally**
- Run tests before committing
- Run linter before pushing
- Use lefthook to automate this

❌ **Don't use v1 golangci-lint**
- Must use v2 (configuration incompatible)
- Install: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`

### Commit Message Format

Use **Conventional Commits**:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `style` - Formatting, missing semicolons, etc.
- `refactor` - Code change that neither fixes a bug nor adds a feature
- `test` - Adding or updating tests
- `chore` - Maintenance tasks
- `perf` - Performance improvement
- `ci` - CI/CD changes

**Examples:**
```
feat(backup): add SSH key authentication support

fix(ssh): handle connection timeout correctly

docs: update installation instructions

ci: add benchmark comparison workflow
```

The commit message format is enforced by lefthook's `commit-msg` hook.

## Quick Reference

### Essential Commands

```bash
# Build
go build ./cmd/mikrotik-backup

# Test
go test -v -race ./...

# Lint
golangci-lint run

# Format
gofumpt -l -w .
goimports -w -local github.com/mountain-reverie/mikrotik-configuation-backup .

# Clean
go clean
go mod tidy
```

### File Locations

- **Main CLI:** `cmd/mikrotik-backup/main.go`
- **Core Logic:** `internal/backup/backup.go`
- **Tests:** `internal/backup/backup_test.go`
- **Integration Tests:** `internal/backup/backup_integration_test.go`
- **Workflows:** `.github/workflows/*.yml`
- **Linter Config:** `.golangci.yml`
- **Release Config:** `.goreleaser.yml`
- **Git Hooks:** `lefthook.yml`

### Important Links

- **Repository:** https://github.com/mountain-reverie/mikrotik-configuation-backup
- **Documentation:** https://mountain-reverie.github.io/mikrotik-configuation-backup/
- **Issues:** https://github.com/mountain-reverie/mikrotik-configuation-backup/issues
- **Discussions:** https://github.com/mountain-reverie/mikrotik-configuation-backup/discussions

### Getting Help

- Read [README.md](README.md) for user documentation
- Read [CONTRIBUTING.md](CONTRIBUTING.md) for contributor guide
- Read [.github/workflows/README.md](.github/workflows/README.md) for CI/CD details
- Check existing issues and discussions
- Review code review comments in closed PRs

## AI Agent Best Practices

### When Contributing Code

1. **Read existing code first** - Understand patterns and conventions
2. **Follow established patterns** - Don't introduce new paradigms
3. **Write tests** - Every new feature needs tests
4. **Run checks locally** - Don't rely solely on CI
5. **Use conventional commits** - Enables automatic changelog generation
6. **Update documentation** - Keep README and CONTRIBUTING in sync

### When Analyzing Issues

1. **Check file references** - Use exact file paths with line numbers
2. **Understand the context** - Read related code, not just the error
3. **Consider side effects** - Changes may affect tests, docs, CI
4. **Verify assumptions** - Check actual behavior, not just documentation

### When Suggesting Changes

1. **Be specific** - Provide exact file paths and line numbers
2. **Show examples** - Include code snippets demonstrating the change
3. **Explain rationale** - Why this change is needed
4. **Consider alternatives** - Acknowledge trade-offs
5. **Test locally** - Verify changes work before suggesting

---

**Last Updated:** 2025-01-05
**Version:** 1.0.0

For questions or improvements to this documentation, please open an issue or discussion on GitHub.
