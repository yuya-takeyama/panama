package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsWorkspace(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(dir string) error
		want      bool
	}{
		{
			name: "with .git directory",
			setupFunc: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".git"), 0755)
			},
			want: true,
		},
		{
			name: "with package.json",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
			},
			want: true,
		},
		{
			name: "with go.mod",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)
			},
			want: true,
		},
		{
			name: "with terraform.tf",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "terraform.tf"), []byte("terraform {}"), 0644)
			},
			want: true,
		},
		{
			name: "with main.tf",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "main.tf"), []byte("provider {}"), 0644)
			},
			want: true,
		},
		{
			name: "empty directory",
			setupFunc: func(dir string) error {
				return nil
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setupFunc(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := IsWorkspace(dir)
			if got != tt.want {
				t.Errorf("IsWorkspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPackageType(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(dir string) error
		want      string
	}{
		{
			name: "Node.js project",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
			},
			want: "node",
		},
		{
			name: "Go project",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644)
			},
			want: "go",
		},
		{
			name: "Terraform project with terraform.tf",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "terraform.tf"), []byte("terraform {}"), 0644)
			},
			want: "terraform",
		},
		{
			name: "Terraform project with main.tf",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "main.tf"), []byte("provider {}"), 0644)
			},
			want: "terraform",
		},
		{
			name: "Rust project",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[package]"), 0644)
			},
			want: "rust",
		},
		{
			name: "Python project with pyproject.toml",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[tool.poetry]"), 0644)
			},
			want: "python",
		},
		{
			name: "No package file",
			setupFunc: func(dir string) error {
				return nil
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setupFunc(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := GetPackageType(dir)
			if got != tt.want {
				t.Errorf("GetPackageType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasPackageFile(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(dir string) error
		want      bool
	}{
		{
			name: "has package.json",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0644)
			},
			want: true,
		},
		{
			name: "has terraform.tf",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "terraform.tf"), []byte("terraform {}"), 0644)
			},
			want: true,
		},
		{
			name: "has main.tf",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "main.tf"), []byte("provider {}"), 0644)
			},
			want: true,
		},
		{
			name: "no package files",
			setupFunc: func(dir string) error {
				return nil
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tt.setupFunc(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := HasPackageFile(dir)
			if got != tt.want {
				t.Errorf("HasPackageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
