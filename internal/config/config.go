package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	UI          string     `yaml:"ui" toml:"ui" json:"ui"`
	MaxDepth    int        `yaml:"max_depth" toml:"max_depth" json:"max_depth"`
	Format      string     `yaml:"format" toml:"format" json:"format"`
	Silent      bool       `yaml:"silent" toml:"silent" json:"silent"`
	NoCache     bool       `yaml:"no_cache" toml:"no_cache" json:"no_cache"`
	Workspaces  []string   `yaml:"workspaces" toml:"workspaces" json:"workspaces"`
	IgnoreDirs  []string   `yaml:"ignore_dirs" toml:"ignore_dirs" json:"ignore_dirs"`
	ScoreConfig ScoreConfig `yaml:"score" toml:"score" json:"score"`
}

type ScoreConfig struct {
	RecentAccessWeight float64 `yaml:"recent_access_weight" toml:"recent_access_weight" json:"recent_access_weight"`
	FrequencyWeight    float64 `yaml:"frequency_weight" toml:"frequency_weight" json:"frequency_weight"`
	DepthPenalty       float64 `yaml:"depth_penalty" toml:"depth_penalty" json:"depth_penalty"`
}

func DefaultConfig() *Config {
	return &Config{
		UI:       "fuzzyfinder",
		MaxDepth: 3,
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
		},
		ScoreConfig: ScoreConfig{
			RecentAccessWeight: 0.5,
			FrequencyWeight:    0.3,
			DepthPenalty:       0.1,
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
		return cfg
	}

	// Search for config file upward from rootDir
	dir := rootDir
	for {
		for _, name := range []string{".panama.yaml", ".panama.yml", ".panama.toml", ".panama.json"} {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				if err := loadFromFile(path, cfg); err != nil {
					log.Printf("Warning: failed to load config from %s: %v", path, err)
				}
				return cfg
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

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
	case ".toml":
		return toml.Unmarshal(data, cfg)
	case ".json":
		return json.Unmarshal(data, cfg)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

func (c *Config) Validate() error {
	if c.MaxDepth < 1 {
		return fmt.Errorf("max_depth must be at least 1")
	}

	if c.UI != "fuzzyfinder" && c.UI != "stdio" {
		return fmt.Errorf("ui must be either 'fuzzyfinder' or 'stdio'")
	}

	if c.Format != "path" && c.Format != "cd" && c.Format != "json" {
		return fmt.Errorf("format must be one of: path, cd, json")
	}

	return nil
}