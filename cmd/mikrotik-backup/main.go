// Package main is the entry point for the mikrotik-backup CLI application.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	defaultSSHPort = 22
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	app := &cli.App{
		Name:  "mikrotik-backup",
		Usage: "MikroTik RouterOS configuration backup tool",
		Description: `A CLI tool to backup MikroTik RouterOS configurations.
This tool connects to MikroTik devices via SSH and exports their configurations
to local files for version control and disaster recovery.`,
		Version: getVersion(),
		Authors: []*cli.Author{
			{
				Name: "Mountain Reverie",
			},
		},
		Commands: []*cli.Command{
			backupCommand(),
			versionCommand(),
		},
		EnableBashCompletion: true,
	}

	err := app.RunContext(ctx, os.Args)
	cancel()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func backupCommand() *cli.Command {
	return &cli.Command{
		Name:  "backup",
		Usage: "Backup MikroTik configuration",
		Description: `Connect to a MikroTik device and backup its configuration.
Supports both password and SSH key-based authentication.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"H"},
				Usage:    "MikroTik device hostname or IP address",
				Required: true,
				EnvVars:  []string{"MIKROTIK_HOST"},
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "SSH port",
				Value:   defaultSSHPort,
				EnvVars: []string{"MIKROTIK_PORT"},
			},
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "SSH username",
				Value:   "admin",
				EnvVars: []string{"MIKROTIK_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"P"},
				Usage:   "SSH password (use with caution, prefer SSH key)",
				EnvVars: []string{"MIKROTIK_PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "Path to SSH private key file",
				EnvVars: []string{"MIKROTIK_KEY_FILE"},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path for the backup",
				Value:   "backup.rsc",
			},
		},
		Action: runBackup,
	}
}

func runBackup(c *cli.Context) error {
	// TODO: Implement backup logic
	_, _ = fmt.Fprintf(c.App.Writer, "Backing up configuration from %s:%d\n", c.String("host"), c.Int("port"))
	_, _ = fmt.Fprintf(c.App.Writer, "Username: %s\n", c.String("username"))
	_, _ = fmt.Fprintf(c.App.Writer, "Output: %s\n", c.String("output"))

	// Validate authentication method
	password := c.String("password")
	keyFile := c.String("key")
	if password == "" && keyFile == "" {
		return errors.New("either --password or --key must be provided")
	}

	return errors.New("not implemented yet")
}

func versionCommand() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print version information",
		Action:  printVersion,
	}
}

func printVersion(c *cli.Context) error {
	info := getBuildInfo()
	_, _ = fmt.Fprintf(c.App.Writer, "mikrotik-backup version %s\n", info.Version)
	_, _ = fmt.Fprintf(c.App.Writer, "  commit: %s\n", info.Commit)
	_, _ = fmt.Fprintf(c.App.Writer, "  built:  %s\n", info.Date)
	_, _ = fmt.Fprintf(c.App.Writer, "  go:     %s\n", info.GoVersion)
	return nil
}

// BuildInfo contains version information about the build.
type BuildInfo struct {
	Version   string
	Commit    string
	Date      string
	GoVersion string
}

func getVersion() string {
	return getBuildInfo().Version
}

func getBuildInfo() BuildInfo {
	info := BuildInfo{
		Version:   "dev",
		Commit:    "none",
		Date:      "unknown",
		GoVersion: "unknown",
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		info.GoVersion = buildInfo.GoVersion

		const shortCommitLength = 7
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				if len(setting.Value) > shortCommitLength {
					info.Commit = setting.Value[:shortCommitLength]
				} else {
					info.Commit = setting.Value
				}
			case "vcs.time":
				info.Date = setting.Value
			}
		}

		// If main module has a version, use it
		if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
			info.Version = buildInfo.Main.Version
		}
	}

	return info
}

// Ensure io.Writer is imported and used.
var _ io.Writer
