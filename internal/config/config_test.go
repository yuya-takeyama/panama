package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.MaxDepth != 3 {
		t.Errorf("expected MaxDepth to be 3, got %d", cfg.MaxDepth)
	}

	if cfg.Format != "path" {
		t.Errorf("expected Format to be 'path', got '%s'", cfg.Format)
	}

	if len(cfg.IgnoreDirs) == 0 {
		t.Error("expected IgnoreDirs to have default values")
	}

	// Check if .terraform is in the ignore list
	hasTerraform := false
	for _, dir := range cfg.IgnoreDirs {
		if dir == ".terraform" {
			hasTerraform = true
			break
		}
	}
	if !hasTerraform {
		t.Error("expected .terraform to be in IgnoreDirs")
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
				MaxDepth: 3,
				Format:   "path",
			},
			wantErr: false,
		},
		{
			name: "invalid max depth",
			config: Config{
				MaxDepth: 0,
				Format:   "path",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: Config{
				MaxDepth: 3,
				Format:   "invalid",
			},
			wantErr: true,
		},
		{
			name: "json format",
			config: Config{
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

	if cfg.MaxDepth != 5 {
		t.Errorf("expected MaxDepth to be 5, got %d", cfg.MaxDepth)
	}

	if cfg.Format != "json" {
		t.Errorf("expected Format to be 'json', got '%s'", cfg.Format)
	}
}
