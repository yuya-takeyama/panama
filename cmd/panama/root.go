package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yuya-takeyama/panama/internal/output"
)

type rootOptions struct {
	format string
	config string
}

func newRootCommand() *cobra.Command {
	opts := &rootOptions{}

	cmd := &cobra.Command{
		Use:   "root",
		Short: "Print the root directory containing panama config or .git",
		Long: `Print the path to the first parent directory containing a panama configuration file or .git directory.
This is useful for navigating to the monorepo or project root directory.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "path", "Output format (path|cd|json)")
	flags.StringVar(&opts.config, "config", "", "Path to configuration file")

	return cmd
}

func runRoot(opts *rootOptions) error {
	// Parse output format
	format, err := output.ParseFormat(opts.format)
	if err != nil {
		return err
	}

	// If config path is provided, use its directory
	if opts.config != "" {
		configDir := filepath.Dir(opts.config)
		absDir, err := filepath.Abs(configDir)
		if err != nil {
			return fmt.Errorf("failed to resolve config directory: %w", err)
		}
		return output.Print(absDir, format)
	}

	// Search for config file or .git directory upward from current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	dir := currentDir
	for {
		// First, check for panama config files
		for _, name := range []string{".panama.yaml", ".panama.yml"} {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				// Found config file, return this directory
				return output.Print(dir, format)
			}
		}

		// Also check for .git directory as fallback
		gitPath := filepath.Join(dir, ".git")
		if stat, err := os.Stat(gitPath); err == nil && stat.IsDir() {
			// Found .git directory, return this directory
			return output.Print(dir, format)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// No config or .git found
	return fmt.Errorf("no root workspace found in any parent directory")
}
