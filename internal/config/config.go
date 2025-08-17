package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	MaxDepth   int      `yaml:"max_depth"`
	Format     string   `yaml:"format"`
	Silent     bool     `yaml:"silent"`
	NoCache    bool     `yaml:"no_cache"`
	Workspaces []string `yaml:"workspaces"`
	IgnoreDirs []string `yaml:"ignored_dirs"`
	Patterns   []string `yaml:"patterns"` // Custom workspace detection patterns
	ConfigDir  string   `yaml:"-"`        // Directory where config was found
}

func DefaultConfig() *Config {
	return &Config{
		MaxDepth: 6,
		Format:   "path",
		Silent:   false,
		NoCache:  false,
		IgnoreDirs: []string{
			"node_modules",
			".git",
			"vendor",
			"target",
			"dist",
			"build",
			".next",
			".nuxt",
			".cache",
			"__pycache__",
			".terraform",
		},
	}
}

func Load(configPath, rootDir string) *Config {
	cfg := DefaultConfig()

	// If config path is provided, use it directly
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			log.Printf("Warning: failed to load config from %s: %v", configPath, err)
		}
		cfg.ConfigDir = filepath.Dir(configPath)
		return cfg
	}

	// Search for config file upward from rootDir
	dir := rootDir
	for {
		for _, name := range []string{".panama.yaml", ".panama.yml"} {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				if err := loadFromFile(path, cfg); err != nil {
					log.Printf("Warning: failed to load config from %s: %v", path, err)
				}
				cfg.ConfigDir = dir // Store the directory where config was found
				return cfg
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// No config found, use rootDir as default
	cfg.ConfigDir = rootDir
	return cfg
}

func loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, cfg)
	default:
		return fmt.Errorf("unsupported config file format: %s (only .yaml and .yml are supported)", ext)
	}
}

func (c *Config) Validate() error {
	if c.MaxDepth < 1 {
		return fmt.Errorf("max_depth must be at least 1")
	}

	if c.Format != "path" && c.Format != "cd" && c.Format != "json" {
		return fmt.Errorf("format must be one of: path, cd, json")
	}

	return nil
}
