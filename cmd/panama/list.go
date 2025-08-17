package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yuya-takeyama/panama/internal/config"
	"github.com/yuya-takeyama/panama/internal/output"
	"github.com/yuya-takeyama/panama/internal/pipeline"
)

type listOptions struct {
	format   string
	maxDepth int
	noCache  bool
	config   string
}

func newListCommand() *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list [path]",
		Short: "List all available workspaces",
		Long: `List all workspaces found in the specified directory or current directory.
Output can be formatted as paths or JSON.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(args, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "path", "Output format (path|json)")
	flags.IntVar(&opts.maxDepth, "max-depth", 0, "Maximum search depth (0 uses config default)")
	flags.BoolVar(&opts.noCache, "no-cache", false, "Disable caching")
	flags.StringVar(&opts.config, "config", "", "Path to configuration file")

	return cmd
}

func runList(args []string, opts *listOptions) error {
	// Determine root directory
	rootDir := "."
	if len(args) > 0 {
		rootDir = args[0]
	}

	// Convert to absolute path
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Load configuration
	cfg := config.Load(opts.config, absRoot)
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Parse output format
	format, err := output.ParseFormat(opts.format)
	if err != nil {
		return err
	}

	// Collect workspaces
	pipelineOpts := pipeline.Options{
		MaxDepth: opts.maxDepth,
		NoCache:  opts.noCache,
	}

	// Use config directory as root if config was found
	searchRoot := absRoot
	if cfg.ConfigDir != "" {
		searchRoot = cfg.ConfigDir
	}

	workspaces, err := pipeline.CollectWorkspaces(searchRoot, cfg, pipelineOpts)
	if err != nil {
		return fmt.Errorf("failed to collect workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		return fmt.Errorf("no workspaces found")
	}

	// Output workspaces
	return output.PrintWorkspaces(workspaces, format)
}
