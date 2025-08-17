package workspace

import (
	"path/filepath"
	"strings"
)

type Workspace struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Depth       int    `json:"depth"`
}

func (w *Workspace) Label() string {
	// Use path as-is for now, will be converted to relative in cmd
	return w.Path
}

func (w *Workspace) LabelWithBase(base string) string {
	// Show relative path from base directory
	return w.RelativePath(base)
}

func (w *Workspace) RelativePath(base string) string {
	rel, err := filepath.Rel(base, w.Path)
	if err != nil {
		return w.Path
	}
	return rel
}

func CalculateDepth(basePath, targetPath string) int {
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return 0
	}
	if rel == "." {
		return 0
	}
	// If the relative path starts with "..", it means targetPath is outside basePath
	if strings.HasPrefix(rel, "..") {
		return 0
	}
	return strings.Count(rel, string(filepath.Separator)) + 1
}
