# MikroTik Configuration Backup

[![Go Reference](https://pkg.go.dev/badge/github.com/mountain-reverie/mikrotik-configuration-backup.svg)](https://pkg.go.dev/github.com/mountain-reverie/mikrotik-configuration-backup)
[![Go Report Card](https://goreportcard.com/badge/github.com/mountain-reverie/mikrotik-configuration-backup)](https://goreportcard.com/report/github.com/mountain-reverie/mikrotik-configuration-backup)
[![Coverage](https://mountain-reverie.github.io/mikrotik-configuration-backup/coverage-badge.svg)](https://mountain-reverie.github.io/mikrotik-configuration-backup/coverage.html)
[![CI](https://github.com/mountain-reverie/mikrotik-configuration-backup/actions/workflows/ci.yml/badge.svg)](https://github.com/mountain-reverie/mikrotik-configuration-backup/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/mountain-reverie/mikrotik-configuration-backup)](LICENSE)

A CLI tool to backup MikroTik RouterOS configurations via SSH.

## Installation

```bash
go install github.com/mountain-reverie/mikrotik-configuration-backup/cmd/mikrotik-backup@latest
```

Or download pre-built binaries from the [releases page](https://github.com/mountain-reverie/mikrotik-configuration-backup/releases). All binaries are signed with cosign - see [verification instructions](.github/workflows/README.md#binary-signing-and-verification).

## Usage

```bash
# Password authentication
mikrotik-backup backup --host 192.168.88.1 --username admin --password mypassword --output backup.rsc

# SSH key authentication
mikrotik-backup backup --host 192.168.88.1 --username admin --key ~/.ssh/mikrotik_rsa --output backup.rsc

# Using environment variables
export MIKROTIK_HOST=192.168.88.1
export MIKROTIK_USERNAME=admin
export MIKROTIK_PASSWORD=mypassword
mikrotik-backup backup --output backup.rsc
```

Run `mikrotik-backup backup --help` for all options.

## Development

```bash
# Run tests
go test -v -race ./...

# Run integration tests
go test -v -tags=integration ./...

# Run linting
golangci-lint run
```

See [AI_AGENTS.md](AI_AGENTS.md) for detailed development documentation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
