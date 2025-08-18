package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) string // Returns temp dir
		opts      *rootOptions
		wantErr   bool
		wantInOut bool // Check if output contains the temp dir path
	}{
		{
			name: "finds config in current directory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, ".panama.yaml")
				if err := os.WriteFile(configPath, []byte("max_depth: 5\n"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.Chdir(tmpDir); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			opts:      &rootOptions{format: "path"},
			wantErr:   false,
			wantInOut: true,
		},
		{
			name: "finds config in parent directory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, ".panama.yaml")
				if err := os.WriteFile(configPath, []byte("max_depth: 5\n"), 0644); err != nil {
					t.Fatal(err)
				}
				subDir := filepath.Join(tmpDir, "subdir")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.Chdir(subDir); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			opts:      &rootOptions{format: "path"},
			wantErr:   false,
			wantInOut: true,
		},
		{
			name: "finds .git directory when no config",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				gitDir := filepath.Join(tmpDir, ".git")
				if err := os.MkdirAll(gitDir, 0755); err != nil {
					t.Fatal(err)
				}
				subDir := filepath.Join(tmpDir, "subdir")
				if err := os.MkdirAll(subDir, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.Chdir(subDir); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			opts:      &rootOptions{format: "path"},
			wantErr:   false,
			wantInOut: true,
		},
		{
			name: "no config or git found",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.Chdir(tmpDir); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			opts:    &rootOptions{format: "path"},
			wantErr: true,
		},
		{
			name: "uses provided config path",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configDir := filepath.Join(tmpDir, "config")
				if err := os.MkdirAll(configDir, 0755); err != nil {
					t.Fatal(err)
				}
				configPath := filepath.Join(configDir, "custom.yaml")
				if err := os.WriteFile(configPath, []byte("max_depth: 5\n"), 0644); err != nil {
					t.Fatal(err)
				}
				return configDir
			},
			opts: &rootOptions{
				format: "path",
				config: filepath.Join(t.TempDir(), "config", "custom.yaml"),
			},
			wantErr:   false,
			wantInOut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current directory to restore later
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(originalDir)

			// Setup test environment
			tmpDir := tt.setup(t)

			// Update config path if it's relative to temp dir
			if tt.opts.config != "" && !filepath.IsAbs(tt.opts.config) {
				tt.opts.config = filepath.Join(tmpDir, "config", "custom.yaml")
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the command
			err = runRoot(tt.opts)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("runRoot() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check output contains expected path
			if tt.wantInOut && err == nil {
				buf := make([]byte, 1024)
				n, _ := r.Read(buf)
				output := string(buf[:n])
				if output == "" {
					t.Errorf("expected output to contain path, got empty string")
				}
			}
		})
	}
}
