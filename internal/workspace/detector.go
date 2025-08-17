package workspace

import (
	"os"
	"path/filepath"
	"strings"
)

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
	// Always check for .git directory
	gitPath := filepath.Join(dir, ".git")
	if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
		return true
	}

	// If no patterns configured, only .git directories are considered workspaces
	if d == nil || len(d.customPatterns) == 0 {
		return false
	}

	// Check custom patterns
	for _, pattern := range d.customPatterns {
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
// Checks common package files if they exist
func GetPackageType(dir string) string {
	// Common package file mappings
	packageTypes := map[string]string{
		"package.json":   "node",
		"go.mod":         "go",
		"pyproject.toml": "python",
		"Cargo.toml":     "rust",
		"pom.xml":        "maven",
		"build.gradle":   "gradle",
		"Gemfile":        "ruby",
		"composer.json":  "php",
	}

	for file, packageType := range packageTypes {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return packageType
		}
	}

	return ""
}
