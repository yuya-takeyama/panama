package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "panama",
		Short: "Fast workspace finder and switcher",
		Long: `panama is a CLI tool that helps you quickly find and navigate to project workspaces.
It automatically detects Git repositories and project directories with package files.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	rootCmd.AddCommand(
		newSelectCommand(),
		newListCommand(),
		newInitCommand(),
		newVersionCommand(),
	)

	return rootCmd.Execute()
}