# GitHub Actions CI/CD Documentation

This directory contains the GitHub Actions workflows for the mikrotik-configuration-backup project.

## Workflows Overview

### 1. CI Workflow (`ci.yml`)

**Trigger:** Reusable workflow, called by other workflows

**Purpose:** Core continuous integration checks

**Jobs:**
- **Lint**: Runs golangci-lint with comprehensive rules
- **Test**: Runs unit tests on latest stable and previous stable Go releases with:
  - `gotestsum` for better test output
  - `goteststats` for test analytics
  - Race detector enabled
  - Coverage reporting to Codecov
  - HTML coverage reports
  - Test result artifacts
- **Integration Test**: Runs integration tests with build tags
- **Build**: Builds binaries for multiple platforms (Linux, macOS, Windows) and architectures (amd64, arm64)
- **Security**: Runs gosec security scanner and uploads SARIF results
- **CodeQL**: Performs CodeQL security analysis
- **Dependency Review**: Reviews dependencies for vulnerabilities (PR only)

### 2. Main Branch Workflow (`main.yml`)

**Trigger:** Push to `main` branch

**Purpose:** Automation for main branch

**Jobs:**
- **CI**: Runs the complete CI workflow
- **Tag**: Automatically creates version tags
  - Calculates next semantic version (patch bump)
  - Generates changelog from commits since last tag
  - Creates and pushes new tag (which automatically triggers the release workflow)

### 3. Pull Request Workflow (`pr.yml`)

**Trigger:** Pull requests to `main` or `develop`

**Purpose:** PR-specific checks and automation

**Jobs:**
- **CI**: Runs the complete CI workflow
- **Coverage Diff**: Compares coverage between base and PR
  - Uses `golang-cover-diff` to show coverage changes
  - Posts comparison as PR comment
- **Dependabot Auto-merge**: Automatically merges safe Dependabot PRs
  - Only for Go modules and GitHub Actions updates
  - Requires CI to pass
- **Claude Auto-fix**: Attempts to fix failing Dependabot PRs
  - Uses Anthropic Claude Code Action
  - Runs build, lint, and tests
  - Commits fixes if found
  - Prevents infinite loops with attempt tracking

### 4. Release Workflow (`release.yml`)

**Trigger:** Push of version tags (`v*`)

**Purpose:** Create GitHub releases with GoReleaser

**Features:**
- Multi-platform builds (Linux, macOS, Windows)
- Multiple architectures (amd64, arm64, arm)
- Creates archives (tar.gz, zip)
- Generates checksums
- **Keyless binary signing with cosign**
  - Uses GitHub OIDC tokens (no private keys to manage)
  - Signs checksums file with ephemeral certificates
  - Records signatures in Rekor transparency log
  - Certificate embeds GitHub workflow and tag information
- **SBOM generation with Syft**
  - Creates Software Bill of Materials for all archives
  - Provides inventory of dependencies
  - Supports supply chain security
- Creates GitHub release with:
  - Changelog grouped by type (features, fixes, etc.)
  - Installation instructions
  - Binary downloads
  - Signature files (.sig) and certificates (.pem)
  - SBOM files (.sbom.json)
- Uploads release assets as artifacts

### 5. GitHub Pages Workflow (`pages.yml`)

**Trigger:** workflow_call (called by main.yml and pr.yml)

**Purpose:** Build and deploy documentation, coverage, and benchmarks

**Features:**
- Generates HTML coverage report
- Creates SVG coverage badge
- Runs performance benchmarks
- Compares benchmarks with previous version using `benchstat`
- Builds beautiful documentation site with:
  - Interactive dashboard showing test stats
  - Coverage reports with drill-down
  - Benchmark results and comparisons
  - Test results in JUnit XML format
- Deploys to GitHub Pages (main branch) or creates preview artifact (PRs)

## Binary Signing and Verification

### How It Works

This project uses **keyless signing with cosign** for maximum security without key management overhead:

1. **During Release** (automated in GitHub Actions):
   - GoReleaser builds binaries and creates checksums.txt
   - cosign signs checksums.txt using GitHub's OIDC token
   - An ephemeral certificate is created (valid for 30 minutes)
   - The signature and certificate are uploaded to the release
   - The signature is recorded in Rekor (public transparency log)

2. **Benefits**:
   - No private keys to store or manage
   - No risk of key compromise or loss
   - Signatures tied to specific GitHub workflow runs
   - Publicly auditable via Rekor transparency log
   - Certificate contains GitHub metadata (repo, workflow, tag)

### Verifying Releases

Users can verify downloaded binaries to ensure they came from this repository's official release workflow:

#### Install cosign
```bash
# macOS
brew install sigstore/tap/cosign

# Linux
wget https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64
chmod +x cosign-linux-amd64
sudo mv cosign-linux-amd64 /usr/local/bin/cosign

# Windows (via winget)
winget install Sigstore.cosign
```

#### Verify the checksums file
```bash
# Download the release files
wget https://github.com/mountain-reverie/mikrotik-configuation-backup/releases/download/v1.0.0/checksums.txt
wget https://github.com/mountain-reverie/mikrotik-configuation-backup/releases/download/v1.0.0/checksums.txt.pem
wget https://github.com/mountain-reverie/mikrotik-configuation-backup/releases/download/v1.0.0/checksums.txt.sig

# Verify the signature
cosign verify-blob checksums.txt \
  --certificate checksums.txt.pem \
  --signature checksums.txt.sig \
  --certificate-identity=https://github.com/mountain-reverie/mikrotik-configuation-backup/.github/workflows/release.yml@refs/tags/v1.0.0 \
  --certificate-oidc-issuer=https://token.actions.githubusercontent.com
```

#### Verify your downloaded binary
Once the checksums file is verified, check your binary matches:
```bash
# Download your platform's binary
wget https://github.com/mountain-reverie/mikrotik-configuation-backup/releases/download/v1.0.0/mikrotik-backup_1.0.0_Linux_x86_64.tar.gz

# Verify checksum
sha256sum --ignore-missing -c checksums.txt
```

**Note**: Replace `v1.0.0` with the actual version you're verifying.

### SBOM (Software Bill of Materials)

Each release includes SBOM files generated by Syft:
- `*.sbom.json` files for each archive
- Contains complete inventory of all dependencies
- Useful for security audits and compliance

## Required Secrets

Configure these secrets in your GitHub repository settings:

### Required
- `GITHUB_TOKEN` - Automatically provided by GitHub Actions
  - Must have `id-token: write` permission for cosign OIDC signing

### Optional (for enhanced features)
- `ANTHROPIC_API_KEY` - For Claude Code auto-fix feature
  - Get from: https://console.anthropic.com/
  - Used in: pr.yml claude-fix job
  - Can be disabled by removing the claude-fix job

## GitHub Pages Setup

### Enable GitHub Pages

1. Go to your repository Settings → Pages
2. Under "Build and deployment":
   - Source: **GitHub Actions**
3. Save the settings

The site will be available at: `https://<username>.github.io/<repository-name>/`

### What's Included

The GitHub Pages site includes:
- **Interactive Dashboard**: Real-time stats for tests, coverage, and builds
- **Coverage Report**: Detailed HTML coverage with line-by-line drill-down
- **Coverage Badge**: SVG badge showing current coverage percentage
- **Benchmarks**: Performance benchmark results
- **Benchmark Comparisons**: Statistical comparison with previous results using `benchstat`
- **Test Results**: JUnit XML for CI/CD integration

### PR Previews

Each PR automatically generates a preview of the documentation site as an artifact. Download from the Actions tab to view locally.

## Feature Flags

### Disable Claude Auto-fix
To disable Claude auto-fix for Dependabot PRs, remove or comment out the `claude-fix` job in `pr.yml`.

### Disable Auto-tagging
To disable automatic version tagging on main branch pushes, remove or comment out the `tag` job in `main.yml`.

### Disable Dependabot Auto-merge
To disable automatic merging of Dependabot PRs, remove or comment out the `dependabot-automerge` job in `pr.yml`.

### Disable GitHub Pages
To disable GitHub Pages deployment, remove or comment out the `deploy-pages` job in `main.yml` and `pages-preview` job in `pr.yml`.

## Local Testing

### Test Linting
```bash
golangci-lint run
```

### Test with gotestsum
```bash
go install gotest.tools/gotestsum@latest
gotestsum --format testname -- -v -race -cover ./...
```

### Test Integration
```bash
go test -v -tags=integration ./...
```

### Test Build
```bash
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" ./cmd/mikrotik-backup
```

### Test GoReleaser
```bash
goreleaser release --snapshot --clean --skip=publish
```

### Git Hooks with Lefthook

The project uses **Lefthook**, a Go-native git hooks manager, for local quality checks:

```bash
# Install lefthook
go install github.com/evilmartians/lefthook@latest

# Install hooks to local repository
lefthook install

# Run pre-commit checks manually
lefthook run pre-commit

# Run pre-push checks manually
lefthook run pre-push
```

Lefthook automatically runs before commits and pushes:
- **Pre-commit**: Formatting (gofumpt, goimports), linting (golangci-lint), security (gosec), tests
- **Pre-push**: Build verification and full test suite
- **Commit-msg**: Validates conventional commit message format

Configuration is in `lefthook.yml` at the repository root.

## Workflow Triggers

### Manual Triggers
You can manually trigger workflows from the Actions tab:
- CI workflow: Can be manually dispatched
- PR workflow: Can be manually dispatched

### Automatic Triggers
- **CI**: Automatically runs on every PR and main branch push
- **Main**: Runs on every push to main branch
- **PR**: Runs on every pull request
- **Release**: Runs when version tags are pushed

## Artifacts

### Test Results
- **Location**: `test-results-{go-version}`
- **Contains**: JUnit XML, test stats JSON
- **Retention**: 30 days

### Coverage Reports
- **Location**: `coverage-report`
- **Contains**: coverage.out, coverage.html
- **Retention**: 30 days

### Build Binaries
- **Location**: `binary-{goos}-{goarch}`
- **Contains**: Compiled binaries for each platform
- **Retention**: 7 days

### Release Assets
- **Location**: `release-assets`
- **Contains**: All release files from GoReleaser
- **Retention**: 90 days

## Best Practices

1. **Conventional Commits**: Use conventional commit format for automatic changelog generation
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation
   - `ci:` for CI changes
   - `test:` for test changes

2. **Version Tags**: Follow semantic versioning (`v1.2.3`)
   - Main workflow auto-increments patch version
   - Manual tags for major/minor bumps

3. **PR Size**: Keep PRs focused and small for easier review and faster CI

4. **Test Coverage**: Maintain or improve coverage with each PR

5. **Security**: Review security scan results in the Security tab

## Troubleshooting

### CI Failing
1. Check the specific job that failed
2. Review the logs for error messages
3. Run the equivalent command locally
4. Check golangci-lint version compatibility

### GoReleaser Failing
1. Ensure tag format is `v*.*.*`
2. Check goreleaser.yml syntax
3. Verify all files referenced exist
4. Test locally with `--snapshot` flag

### Dependabot Auto-merge Not Working
1. Verify PR is from Dependabot
2. Check if CI passed
3. Ensure PR title matches safe patterns
4. Verify GitHub token has necessary permissions

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GoReleaser Documentation](https://goreleaser.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [Conventional Commits](https://www.conventionalcommits.org/)
