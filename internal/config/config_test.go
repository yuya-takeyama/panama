package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.UI != "fuzzyfinder" {
		t.Errorf("expected UI to be 'fuzzyfinder', got '%s'", cfg.UI)
	}

	if cfg.MaxDepth != 3 {
		t.Errorf("expected MaxDepth to be 3, got %d", cfg.MaxDepth)
	}

	if cfg.Format != "path" {
		t.Errorf("expected Format to be 'path', got '%s'", cfg.Format)
	}

	if len(cfg.IgnoreDirs) == 0 {
		t.Error("expected IgnoreDirs to have default values")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				UI:       "fuzzyfinder",
				MaxDepth: 3,
				Format:   "path",
			},
			wantErr: false,
		},
		{
			name: "invalid UI",
			config: Config{
				UI:       "invalid",
				MaxDepth: 3,
				Format:   "path",
			},
			wantErr: true,
		},
		{
			name: "invalid max depth",
			config: Config{
				UI:       "fuzzyfinder",
				MaxDepth: 0,
				Format:   "path",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: Config{
				UI:       "fuzzyfinder",
				MaxDepth: 3,
				Format:   "invalid",
			},
			wantErr: true,
		},
		{
			name: "stdio UI",
			config: Config{
				UI:       "stdio",
				MaxDepth: 3,
				Format:   "cd",
			},
			wantErr: false,
		},
		{
			name: "json format",
			config: Config{
				UI:       "fuzzyfinder",
				MaxDepth: 3,
				Format:   "json",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "panama-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test YAML config
	yamlPath := filepath.Join(tmpDir, ".panama.yaml")
	yamlContent := `
ui: stdio
max_depth: 5
format: json
`
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := DefaultConfig()
	if err := loadFromFile(yamlPath, cfg); err != nil {
		t.Errorf("failed to load YAML config: %v", err)
	}

	if cfg.UI != "stdio" {
		t.Errorf("expected UI to be 'stdio', got '%s'", cfg.UI)
	}

	if cfg.MaxDepth != 5 {
		t.Errorf("expected MaxDepth to be 5, got %d", cfg.MaxDepth)
	}

	if cfg.Format != "json" {
		t.Errorf("expected Format to be 'json', got '%s'", cfg.Format)
	}
}
