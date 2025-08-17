package fuzzyfinder

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// StdioSelect provides a fallback selection method using standard I/O
func StdioSelect(items []Item) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("no items to select from")
	}

	// Display items
	fmt.Fprintln(os.Stderr, "Available workspaces:")
	for i, item := range items {
		fmt.Fprintf(os.Stderr, "  %d) %s\n", i+1, item.Label)
		if item.Description != "" {
			fmt.Fprintf(os.Stderr, "     %s\n", item.Description)
		}
	}

	// Prompt for selection
	fmt.Fprint(os.Stderr, "Select workspace (1-", len(items), "): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("invalid input: %s", input)
	}

	if choice < 1 || choice > len(items) {
		return -1, fmt.Errorf("choice %d is out of range (1-%d)", choice, len(items))
	}

	return choice - 1, nil
}
