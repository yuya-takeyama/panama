package pipeline

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/yuya-takeyama/panama/internal/config"
	"github.com/yuya-takeyama/panama/internal/workspace"
)

type Options struct {
	Query    string
	MaxDepth int
	NoCache  bool
}

func CollectWorkspaces(rootDir string, cfg *config.Config, opts Options) ([]*workspace.Workspace, error) {
	workspaces := []*workspace.Workspace{}
	maxDepth := cfg.MaxDepth
	if opts.MaxDepth > 0 {
		maxDepth = opts.MaxDepth
	}

	// Create ignore patterns
	ignorePatterns := make([]string, len(cfg.IgnoreDirs))
	for i, dir := range cfg.IgnoreDirs {
		ignorePatterns[i] = "**/" + dir
	}

	// Create detector with custom patterns
	detector := workspace.NewDetector(cfg.Patterns)

	visited := make(map[string]bool)

	// If specific workspace directories are configured, search those
	if len(cfg.Workspaces) > 0 {
		for _, wsPath := range cfg.Workspaces {
			absPath := wsPath
			if !filepath.IsAbs(wsPath) {
				absPath = filepath.Join(rootDir, wsPath)
			}

			if err := collectFromPath(absPath, rootDir, maxDepth, ignorePatterns, visited, detector, &workspaces); err != nil {
				// Log but continue with other paths
				continue
			}
		}
	} else {
		// Search from root directory
		if err := collectFromPath(rootDir, rootDir, maxDepth, ignorePatterns, visited, detector, &workspaces); err != nil {
			return nil, err
		}
	}

	// Sort workspaces by path
	sort.Slice(workspaces, func(i, j int) bool {
		return workspaces[i].Path < workspaces[j].Path
	})

	return workspaces, nil
}

func collectFromPath(searchPath, basePath string, maxDepth int, ignorePatterns []string, visited map[string]bool, detector *workspace.Detector, workspaces *[]*workspace.Workspace) error {
	return filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip on error
		}

		if !info.IsDir() {
			return nil
		}

		// Check if already visited
		if visited[path] {
			return filepath.SkipDir
		}
		visited[path] = true

		// Skip if exceeds max depth
		depth := workspace.CalculateDepth(basePath, path)
		if depth > maxDepth {
			return filepath.SkipDir
		}

		// Skip ignored directories
		for _, pattern := range ignorePatterns {
			if matched, _ := doublestar.Match(pattern, path); matched {
				return filepath.SkipDir
			}
		}

		// Check if it's a workspace
		if detector.IsWorkspaceWithPatterns(path) {
			ws := &workspace.Workspace{
				Path:  path,
				Name:  filepath.Base(path),
				Depth: depth,
			}

			// Add package type as description
			if packageType := workspace.GetPackageType(path); packageType != "" {
				ws.Description = "Type: " + packageType
			}

			*workspaces = append(*workspaces, ws)

			// Don't recurse into detected workspaces
			if path != searchPath {
				return filepath.SkipDir
			}
		}

		return nil
	})
}
