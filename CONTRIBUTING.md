# Contributing to MikroTik Configuration Backup

Thank you for your interest in contributing to this project! This document provides guidelines and instructions for contributing.

> **For AI Agents:** If you're an AI agent (like Claude Code, GitHub Copilot, ChatGPT, etc.), please see [AI_AGENTS.md](AI_AGENTS.md) for comprehensive technical documentation including project structure, coding patterns, testing guidelines, and CI/CD workflows tailored for AI-assisted development.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [CI/CD Workflows](#cicd-workflows)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a code of conduct based on respect and professionalism. By participating, you are expected to:

- Be respectful and inclusive
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go (latest stable)** - [Download](https://go.dev/doc/install)
  - This project always uses the latest stable Go release
  - Check `go.mod` for the current minimum version
- **Git** - [Download](https://git-scm.com/downloads)
- **golangci-lint v2** - [Installation guide](https://golangci-lint.run/docs/welcome/install/)
  ```bash
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
  ```

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/mikrotik-configuation-backup.git
   cd mikrotik-configuation-backup
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/mountain-reverie/mikrotik-configuation-backup.git
   ```

### Set Up Development Environment

1. Download dependencies:
   ```bash
   go mod download
   ```

2. Verify your setup:
   ```bash
   go test ./...
   golangci-lint run
   ```

## Development Workflow

### 1. Create a Branch

Always create a new branch for your work:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test improvements
- `chore/` - Maintenance tasks

### 2. Make Your Changes

- Write clean, idiomatic Go code
- Follow the [Coding Standards](#coding-standards)
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests and Linting

Before committing, ensure all checks pass:

```bash
go test -v -race ./...
golangci-lint run
```

This runs:
- Unit tests with race detector
- Code linting with 60+ linters

### 4. Commit Your Changes

Follow the [Commit Messages](#commit-messages) guidelines:

```bash
git add .
git commit -m "feat: add new backup scheduling feature"
```

### 5. Keep Your Branch Updated

Regularly sync with upstream:

```bash
git fetch upstream
git rebase upstream/main
```

### 6. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 7. Create a Pull Request

1. Go to the repository on GitHub
2. Click "New Pull Request"
3. Select your branch
4. Fill out the PR template
5. Submit the pull request

## Coding Standards

This project follows Go best practices and enforces them through automated linting.

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofumpt` formatting (stricter than `gofmt`)
- Use `golangci-lint` with the project's configuration

### Key Principles

1. **Interfaces** - Use interfaces for dependency injection and testability
2. **Error Handling** - Always wrap errors with context using `fmt.Errorf` with `%w`
3. **Context** - Propagate context.Context for cancellation and timeouts
4. **Naming** - Use clear, descriptive names
5. **Comments** - Add godoc comments for exported functions, types, and packages
6. **Simplicity** - Keep functions small and focused

### Example Code Style

```go
package backup

import (
    "context"
    "fmt"
)

// Service handles backup operations for MikroTik devices.
type Service struct {
    client SSHClient
}

// New creates a new backup service with the provided SSH client.
func New(client SSHClient) *Service {
    return &Service{
        client: client,
    }
}

// Execute performs a backup operation and writes the result to the output.
func (s *Service) Execute(ctx context.Context, config Config, output io.Writer) error {
    if err := s.client.Connect(ctx, config); err != nil {
        return fmt.Errorf("failed to connect to device: %w", err)
    }
    defer s.client.Close()

    // ... implementation
    return nil
}
```

### Code Organization

- Put public APIs in package root
- Put internal implementation in `internal/`
- Put shared utilities in `pkg/`
- Keep files focused and cohesive
- Group related functionality together

## Testing Guidelines

### Test Coverage

- Aim for at least 80% code coverage
- All new features must include tests
- All bug fixes must include regression tests

### Unit Tests

```go
func TestService_Execute_Success(t *testing.T) {
    t.Parallel() // Run tests in parallel when possible

    // Arrange
    mockClient := &mockSSHClient{
        executeCommandFunc: func(ctx context.Context, cmd string) (string, error) {
            return "config output", nil
        },
    }
    service := backup.New(mockClient)
    output := &bytes.Buffer{}

    // Act
    err := service.Execute(context.Background(), config, output)

    // Assert
    if err != nil {
        t.Fatalf("Execute() error = %v, want nil", err)
    }
}
```

### Integration Tests

- Use build tags: `//go:build integration`
- Place in `test/integration/`
- Document required environment setup
- Skip if environment is not configured

```go
//go:build integration

package integration_test

func TestBackupIntegration(t *testing.T) {
    host := os.Getenv("MIKROTIK_HOST")
    if host == "" {
        t.Skip("Skipping: MIKROTIK_HOST not set")
    }
    // ... test implementation
}
```

### Running Tests

```bash
# Unit tests
go test ./...

# Unit tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Integration tests
go test -tags=integration ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

## Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks
- `perf` - Performance improvements
- `ci` - CI/CD changes

### Examples

```
feat(backup): add support for scheduled backups

Implement cron-based scheduling for automatic backups.
Adds new --schedule flag to backup command.

Closes #123
```

```
fix(ssh): handle connection timeout correctly

Previously, connection timeouts would panic. Now they
return a proper error with context.

Fixes #456
```

### Best Practices

- Use imperative mood ("add feature" not "added feature")
- Keep subject line under 50 characters
- Capitalize subject line
- Don't end subject with a period
- Separate subject from body with blank line
- Wrap body at 72 characters
- Use body to explain what and why, not how

## Pull Request Process

### PR Checklist

Before submitting a PR, ensure:

- [ ] Code follows project style guidelines
- [ ] All tests pass (`go test ./...`)
- [ ] Linting passes (`golangci-lint run`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] PR description is clear and complete
- [ ] Changes are covered by tests

### PR Description Template

```markdown
## Description
Brief description of the changes

## Motivation and Context
Why is this change required? What problem does it solve?

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## How Has This Been Tested?
Describe the tests you ran and how to reproduce them.

## Screenshots (if appropriate)

## Checklist
- [ ] My code follows the code style of this project
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] All new and existing tests pass
- [ ] I have updated the documentation accordingly
```

### Review Process

1. Automated CI checks must pass
2. At least one maintainer approval required
3. Address review comments
4. Keep the PR focused and atomic
5. Rebase on main if needed

### After Approval

- Maintainers will merge using "Squash and Merge"
- Delete your branch after merge
- Update your local repository:
  ```bash
  git checkout main
  git pull upstream main
  ```

## CI/CD Workflows

This project uses comprehensive GitHub Actions workflows for automation. All workflows are fully documented in [.github/workflows/README.md](.github/workflows/README.md).

### Overview

The project includes 6 main workflows:

1. **Go Version Update** (`go-version-update.yml`) - Monthly check for new Go releases, creates PR to update go.mod
2. **CI Workflow** (`ci.yml`) - Reusable workflow for linting, testing, building, and security scanning
3. **Main Branch** (`main.yml`) - Auto-tags releases and deploys GitHub Pages on main branch pushes
4. **Pull Requests** (`pr.yml`) - Runs CI, coverage diff, Dependabot auto-merge, and Claude auto-fix
5. **Release** (`release.yml`) - Creates signed releases with GoReleaser when tags are pushed
6. **GitHub Pages** (`pages.yml`) - Builds and deploys documentation, coverage, and benchmarks

### Key Features

- **Automated Testing**: Unit tests, integration tests, and benchmarks on latest stable and previous stable Go releases
- **Code Quality**: golangci-lint with 60+ linters, gosec security scanning, CodeQL analysis
- **Binary Signing**: Keyless signing with cosign using GitHub OIDC tokens
- **SBOM Generation**: Software Bill of Materials for supply chain security
- **Coverage Reports**: Self-hosted on GitHub Pages with interactive dashboards
- **Auto-tagging**: Automatic semantic versioning on main branch
- **Dependabot Integration**: Auto-merge safe updates, Claude auto-fix for failures
- **Go Version Updates**: Monthly automated checks for new Go releases
- **GitHub Actions Security**: All actions pinned to commit SHAs (not tags) for supply chain security

### For Contributors

When you create a pull request:
- CI automatically runs all checks
- Coverage diff is posted as a comment
- Test results and artifacts are available in the Actions tab
- Documentation preview is generated as an artifact

For detailed information about workflows, secrets, configuration, and troubleshooting, see the [GitHub Actions documentation](.github/workflows/README.md).

### GitHub Actions Security Policy

**CRITICAL: All GitHub Actions MUST be pinned to commit SHAs, not tags.**

This is a security requirement to prevent supply chain attacks through tag poisoning. Tags are mutable and can be changed to point to malicious code, while commit SHAs are immutable.

**Correct format:**
```yaml
# ✅ GOOD - Pinned to SHA with version comment
- uses: actions/checkout@08eba0b27e820071cde6df949e0beb9ba4906955 # v4.3.0

# ❌ BAD - Using tag (mutable, vulnerable to attacks)
- uses: actions/checkout@v4
```

**When updating GitHub Actions:**

1. Find the release on the action's GitHub repository
2. Get the full 40-character commit SHA from the release tag:
   ```bash
   git ls-remote https://github.com/actions/checkout refs/tags/v4.3.0
   ```
3. Update all workflow files with the SHA
4. Update `.github/ACTION_SHAS.md` with the new version, SHA, and date
5. Always add a version comment after the SHA for readability

**Documentation:**
- All action SHAs are documented in `.github/ACTION_SHAS.md`
- Includes version, commit SHA, and last updated date
- Provides instructions for finding and updating SHAs

**Why this matters:**
- Prevents attackers from moving tags to malicious commits
- Ensures reproducible builds with exact versions
- Provides explicit control over dependency updates
- Aligns with security best practices (SLSA, OpenSSF Scorecards)

## Release Process

Releases are automated using GoReleaser and GitHub Actions.

### Creating a Release

**Option 1: Automatic (Recommended)**
- Push to `main` branch
- GitHub Actions automatically creates a tag (patch version bump)
- Tag push triggers the release workflow

**Option 2: Manual (for major/minor versions)**
- Only maintainers can create manual releases
- Create and push a version tag:
  ```bash
  git tag -a v1.2.3 -m "Release v1.2.3"
  git push origin v1.2.3
  ```

### What Happens During Release

The release workflow automatically:
- Builds binaries for all platforms (Linux, macOS, Windows)
- Runs all tests and security scans
- Signs binaries with cosign (keyless, using GitHub OIDC)
- Generates SBOM for all artifacts
- Creates GitHub release with changelog
- Uploads signed binaries, signatures, certificates, and SBOMs

### Version Numbering

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** - Incompatible API changes
- **MINOR** - Backwards-compatible new functionality
- **PATCH** - Backwards-compatible bug fixes (auto-incremented on main branch)

## Getting Help

If you need help:

- 📖 Read the [README.md](README.md)
- 💬 Start a [Discussion](https://github.com/mountain-reverie/mikrotik-configuation-backup/discussions)
- 🐛 Check existing [Issues](https://github.com/mountain-reverie/mikrotik-configuation-backup/issues)

## Additional Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Conventional Commits](https://www.conventionalcommits.org/)

## Recognition

Contributors will be recognized in:
- Release notes
- GitHub contributors page
- Project documentation

Thank you for contributing! 🎉
