package workspace

import (
	"path/filepath"
	"strings"
)

type Workspace struct {
	Path        string  `json:"path"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Score       float64 `json:"score,omitempty"`
	Depth       int     `json:"depth"`
	HasGit      bool    `json:"has_git"`
	HasPackage  bool    `json:"has_package"`
}

func (w *Workspace) Label() string {
	// Use path as-is for now, will be converted to relative in cmd
	label := w.Path
	if w.HasGit {
		label += " [git]"
	}
	if w.HasPackage {
		label += " [pkg]"
	}
	return label
}

func (w *Workspace) LabelWithBase(base string) string {
	// Show relative path from base directory
	label := w.RelativePath(base)
	if w.HasGit {
		label += " [git]"
	}
	if w.HasPackage {
		label += " [pkg]"
	}
	return label
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
