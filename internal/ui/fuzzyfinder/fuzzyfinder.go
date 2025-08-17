package fuzzyfinder

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/term"
)

type Item struct {
	Label       string
	Description string
	Path        string
	Score       float64
}

func Select(items []Item, query string) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("no items to select from")
	}

	// TTY check
	if !isTerminal() {
		// Return first item in non-interactive mode
		return 0, nil
	}

	opts := []fuzzyfinder.Option{
		fuzzyfinder.WithPromptString("workspaces > "),
	}

	if query != "" {
		opts = append(opts, fuzzyfinder.WithQuery(query))
	}

	if len(items) > 0 && items[0].Description != "" {
		opts = append(opts, fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i < 0 || i >= len(items) {
				return ""
			}
			preview := fmt.Sprintf("Path: %s\n\n", items[i].Path)
			if items[i].Description != "" {
				preview += fmt.Sprintf("Description:\n%s", items[i].Description)
			}
			return wrapText(preview, w)
		}))
	}

	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			return items[i].Label
		},
		opts...,
	)

	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return -1, fmt.Errorf("selection cancelled")
		}
		return -1, err
	}

	return idx, nil
}

func SelectMulti(items []Item, query string) ([]int, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no items to select from")
	}

	if !isTerminal() {
		// Return first item in non-interactive mode
		return []int{0}, nil
	}

	opts := []fuzzyfinder.Option{
		fuzzyfinder.WithPromptString("workspaces > "),
	}

	if query != "" {
		opts = append(opts, fuzzyfinder.WithQuery(query))
	}

	if len(items) > 0 && items[0].Description != "" {
		opts = append(opts, fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i < 0 || i >= len(items) {
				return ""
			}
			preview := fmt.Sprintf("Path: %s\n\n", items[i].Path)
			if items[i].Description != "" {
				preview += fmt.Sprintf("Description:\n%s", items[i].Description)
			}
			return wrapText(preview, w)
		}))
	}

	idxs, err := fuzzyfinder.FindMulti(
		items,
		func(i int) string {
			return items[i].Label
		},
		opts...,
	)

	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return nil, fmt.Errorf("selection cancelled")
		}
		return nil, err
	}

	return idxs, nil
}

func isTerminal() bool {
	// Check if stdin is a terminal (fuzzyfinder uses /dev/tty directly)
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result []string
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(line) <= width {
			result = append(result, line)
			continue
		}

		// Simple word wrapping
		for len(line) > width {
			result = append(result, line[:width])
			line = line[width:]
		}
		if len(line) > 0 {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}