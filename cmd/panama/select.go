package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yuya-takeyama/panama/internal/config"
	"github.com/yuya-takeyama/panama/internal/output"
	"github.com/yuya-takeyama/panama/internal/pipeline"
	"github.com/yuya-takeyama/panama/internal/ui/fuzzyfinder"
	"golang.org/x/term"
)

type selectOptions struct {
	query    string
	first    bool
	format   string
	maxDepth int
	noCache  bool
	silent   bool
	config   string
}

func newSelectCommand() *cobra.Command {
	opts := &selectOptions{}

	cmd := &cobra.Command{
		Use:   "select [path]",
		Short: "Select a workspace interactively",
		Long: `Select a workspace using the built-in fuzzy finder.
If no path is provided, it searches from the current directory.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSelect(args, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.query, "query", "q", "", "Initial search query")
	flags.BoolVar(&opts.first, "first", false, "Select first result without interaction")
	flags.StringVarP(&opts.format, "format", "f", "path", "Output format (path|cd|json)")
	flags.IntVar(&opts.maxDepth, "max-depth", 0, "Maximum search depth (0 uses config default)")
	flags.BoolVar(&opts.noCache, "no-cache", false, "Disable caching")
	flags.BoolVar(&opts.silent, "silent", false, "Suppress non-essential output")
	flags.StringVar(&opts.config, "config", "", "Path to configuration file")

	return cmd
}

func runSelect(args []string, opts *selectOptions) error {
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
		Query:    opts.query,
		First:    opts.first,
		MaxDepth: opts.maxDepth,
		NoCache:  opts.noCache,
	}

	workspaces, err := pipeline.CollectWorkspaces(absRoot, cfg, pipelineOpts)
	if err != nil {
		return fmt.Errorf("failed to collect workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		return fmt.Errorf("no workspaces found")
	}

	// Check if we should use interactive mode
	isInteractive := term.IsTerminal(int(os.Stdin.Fd())) && 
		term.IsTerminal(int(os.Stdout.Fd())) && 
		!opts.first &&
		cfg.UI != "stdio"

	var selectedPath string

	if isInteractive {
		// Convert to fuzzyfinder items
		items := make([]fuzzyfinder.Item, len(workspaces))
		for i, ws := range workspaces {
			items[i] = fuzzyfinder.Item{
				Label:       ws.Label(),
				Description: ws.Description,
				Path:        ws.Path,
				Score:       ws.Score,
			}
		}

		// Show fuzzy finder
		idx, err := fuzzyfinder.Select(items, opts.query)
		if err != nil {
			return err
		}

		selectedPath = workspaces[idx].Path
	} else if cfg.UI == "stdio" && !opts.first {
		// Use stdio selection
		items := make([]fuzzyfinder.Item, len(workspaces))
		for i, ws := range workspaces {
			items[i] = fuzzyfinder.Item{
				Label:       ws.Label(),
				Description: ws.Description,
				Path:        ws.Path,
				Score:       ws.Score,
			}
		}

		idx, err := fuzzyfinder.StdioSelect(items)
		if err != nil {
			return err
		}

		selectedPath = workspaces[idx].Path
	} else {
		// Non-interactive mode or --first flag
		selectedPath = workspaces[0].Path
	}

	// Output the selected path
	return output.Print(selectedPath, format)
}