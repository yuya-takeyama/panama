package workspace

import (
	"os"
	"path/filepath"
	"strings"
)

// Default minimal patterns for workspace detection
var defaultPatterns = []string{
	"package.json",   // Node.js
	"go.mod",         // Go
	"pyproject.toml", // Python
}

// Detector holds configuration for workspace detection
type Detector struct {
	customPatterns []string
}

// NewDetector creates a new Detector with custom patterns
func NewDetector(patterns []string) *Detector {
	return &Detector{
		customPatterns: patterns,
	}
}

// IsWorkspaceWithPatterns checks if a directory is a workspace with custom patterns
func (d *Detector) IsWorkspaceWithPatterns(dir string) bool {
	// Check for .git directory
	gitPath := filepath.Join(dir, ".git")
	if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
		return true
	}

	// Use custom patterns if provided, otherwise use defaults
	patterns := d.customPatterns
	if len(patterns) == 0 {
		patterns = defaultPatterns
	}

	// Check patterns
	for _, pattern := range patterns {
		// Check if pattern contains glob characters
		if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
			// Use glob matching
			matches, err := filepath.Glob(filepath.Join(dir, pattern))
			if err == nil && len(matches) > 0 {
				return true
			}
		} else {
			// Simple file existence check
			path := filepath.Join(dir, pattern)
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
	}

	return false
}

// IsWorkspace checks if a directory is a workspace using default patterns
func IsWorkspace(dir string) bool {
	detector := NewDetector(nil)
	return detector.IsWorkspaceWithPatterns(dir)
}

// GetPackageType returns the package type for a directory
// Only checks for the default patterns
func GetPackageType(dir string) string {
	packageTypes := map[string]string{
		"package.json":   "node",
		"go.mod":         "go",
		"pyproject.toml": "python",
	}

	for file, packageType := range packageTypes {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return packageType
		}
	}

	return ""
}
