package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

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

	// Create default configuration
	defaultConfig := map[string]interface{}{
		"max_depth": 3,
		"format":    "path",
		"ignore_dirs": []string{
			"node_modules",
			".git",
			"vendor",
			"target",
			"dist",
			"build",
			".next",
			".nuxt",
			".cache",
			"__pycache__",
		},
	}

	// Write configuration file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}
	defer file.Close()

	// Always write YAML format
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(defaultConfig); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("Configuration file created: %s\n", absPath)
	return nil
}
