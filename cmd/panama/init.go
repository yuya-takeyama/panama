package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed templates/default-config.yaml
var defaultConfigTemplate string

type initOptions struct {
	force bool
}

func newInitCommand() *cobra.Command {
	opts := &initOptions{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a panama configuration file",
		Long:  `Create a .panama.yaml configuration file in the current directory with default settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.force, "force", false, "Overwrite existing configuration file")

	return cmd
}

func runInit(opts *initOptions) error {
	// Always use .panama.yaml
	filename := ".panama.yaml"

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil && !opts.force {
		return fmt.Errorf("configuration file %s already exists (use --force to overwrite)", filename)
	}

	// Write configuration file from template
	if err := os.WriteFile(filename, []byte(defaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("Configuration file created: %s\n", absPath)
	return nil
}
