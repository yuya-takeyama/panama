package pipeline

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

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

	visited := make(map[string]bool)

	// If specific workspace directories are configured, search those
	if len(cfg.Workspaces) > 0 {
		for _, wsPath := range cfg.Workspaces {
			absPath := wsPath
			if !filepath.IsAbs(wsPath) {
				absPath = filepath.Join(rootDir, wsPath)
			}

			if err := collectFromPath(absPath, rootDir, maxDepth, ignorePatterns, visited, &workspaces); err != nil {
				// Log but continue with other paths
				continue
			}
		}
	} else {
		// Search from root directory
		if err := collectFromPath(rootDir, rootDir, maxDepth, ignorePatterns, visited, &workspaces); err != nil {
			return nil, err
		}
	}

	// Score and sort workspaces
	scoreWorkspaces(workspaces, cfg.ScoreConfig, opts.Query)
	sort.Sort(byScore(workspaces))

	return workspaces, nil
}

func collectFromPath(searchPath, basePath string, maxDepth int, ignorePatterns []string, visited map[string]bool, workspaces *[]*workspace.Workspace) error {
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
		if workspace.IsWorkspace(path) {
			ws := &workspace.Workspace{
				Path:       path,
				Name:       filepath.Base(path),
				Depth:      depth,
				HasGit:     workspace.HasGitRepo(path),
				HasPackage: workspace.HasPackageFile(path),
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

func scoreWorkspaces(workspaces []*workspace.Workspace, cfg config.ScoreConfig, query string) {
	for _, ws := range workspaces {
		score := 1.0

		// Depth penalty
		score -= float64(ws.Depth) * cfg.DepthPenalty

		// Bonus for git repos
		if ws.HasGit {
			score += 0.2
		}

		// Bonus for package files
		if ws.HasPackage {
			score += 0.1
		}

		// Query matching bonus
		if query != "" {
			lowerQuery := strings.ToLower(query)
			lowerName := strings.ToLower(ws.Name)
			lowerPath := strings.ToLower(ws.Path)

			if strings.Contains(lowerName, lowerQuery) {
				score += 0.5
			} else if strings.Contains(lowerPath, lowerQuery) {
				score += 0.3
			}
		}

		ws.Score = score
	}
}

type byScore []*workspace.Workspace

func (s byScore) Len() int           { return len(s) }
func (s byScore) Less(i, j int) bool { return s[i].Score > s[j].Score }
func (s byScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
