package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type initOptions struct {
	format string
	force  bool
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
	flags.StringVar(&opts.format, "format", "yaml", "Configuration format (yaml|toml|json)")
	flags.BoolVar(&opts.force, "force", false, "Overwrite existing configuration file")

	return cmd
}

func runInit(opts *initOptions) error {
	// Determine filename based on format
	var filename string
	switch opts.format {
	case "yaml", "yml":
		filename = ".panama.yaml"
	case "toml":
		filename = ".panama.toml"
	case "json":
		filename = ".panama.json"
	default:
		return fmt.Errorf("unsupported format: %s", opts.format)
	}

	// Check if file already exists
	if _, err := os.Stat(filename); err == nil && !opts.force {
		return fmt.Errorf("configuration file %s already exists (use --force to overwrite)", filename)
	}

	// Create default configuration
	defaultConfig := map[string]interface{}{
		"ui":        "fuzzyfinder",
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
		"score": map[string]interface{}{
			"recent_access_weight": 0.5,
			"frequency_weight":     0.3,
			"depth_penalty":        0.1,
		},
	}

	// Write configuration file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}
	defer file.Close()

	switch opts.format {
	case "yaml", "yml":
		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		if err := encoder.Encode(defaultConfig); err != nil {
			return fmt.Errorf("failed to write configuration: %w", err)
		}
	case "toml":
		// For simplicity, we'll output a static TOML template
		tomlContent := `# Panama configuration file

ui = "fuzzyfinder"
max_depth = 3
format = "path"

ignore_dirs = [
  "node_modules",
  ".git",
  "vendor",
  "target",
  "dist",
  "build",
  ".next",
  ".nuxt",
  ".cache",
  "__pycache__"
]

[score]
recent_access_weight = 0.5
frequency_weight = 0.3
depth_penalty = 0.1
`
		if _, err := file.WriteString(tomlContent); err != nil {
			return fmt.Errorf("failed to write configuration: %w", err)
		}
	case "json":
		jsonContent := `{
  "ui": "fuzzyfinder",
  "max_depth": 3,
  "format": "path",
  "ignore_dirs": [
    "node_modules",
    ".git",
    "vendor",
    "target",
    "dist",
    "build",
    ".next",
    ".nuxt",
    ".cache",
    "__pycache__"
  ],
  "score": {
    "recent_access_weight": 0.5,
    "frequency_weight": 0.3,
    "depth_penalty": 0.1
  }
}
`
		if _, err := file.WriteString(jsonContent); err != nil {
			return fmt.Errorf("failed to write configuration: %w", err)
		}
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("Configuration file created: %s\n", absPath)
	return nil
}
