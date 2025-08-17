package workspace

import (
	"testing"
)

func TestCalculateDepth(t *testing.T) {
	tests := []struct {
		name       string
		basePath   string
		targetPath string
		want       int
	}{
		{
			name:       "same directory",
			basePath:   "/home/user/projects",
			targetPath: "/home/user/projects",
			want:       0,
		},
		{
			name:       "one level deep",
			basePath:   "/home/user/projects",
			targetPath: "/home/user/projects/app",
			want:       1,
		},
		{
			name:       "two levels deep",
			basePath:   "/home/user/projects",
			targetPath: "/home/user/projects/frontend/app",
			want:       2,
		},
		{
			name:       "different path",
			basePath:   "/home/user",
			targetPath: "/var/log",
			want:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDepth(tt.basePath, tt.targetPath)
			if got != tt.want {
				t.Errorf("CalculateDepth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkspace_Label(t *testing.T) {
	tests := []struct {
		name string
		ws   Workspace
		want string
	}{
		{
			name: "plain workspace",
			ws: Workspace{
				Path: "/home/user/myproject",
				Name: "myproject",
			},
			want: "/home/user/myproject",
		},
		{
			name: "with git",
			ws: Workspace{
				Path: "/home/user/myproject",
				Name: "myproject",
			},
			want: "/home/user/myproject",
		},
		{
			name: "with package",
			ws: Workspace{
				Path: "/home/user/myproject",
				Name: "myproject",
			},
			want: "/home/user/myproject",
		},
		{
			name: "with both git and package",
			ws: Workspace{
				Path: "/home/user/myproject",
				Name: "myproject",
			},
			want: "/home/user/myproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ws.Label()
			if got != tt.want {
				t.Errorf("Workspace.Label() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkspace_LabelWithBase(t *testing.T) {
	tests := []struct {
		name string
		ws   Workspace
		base string
		want string
	}{
		{
			name: "relative path",
			ws: Workspace{
				Path: "/home/user/projects/myapp",
				Name: "myapp",
			},
			base: "/home/user/projects",
			want: "myapp",
		},
		{
			name: "nested relative path",
			ws: Workspace{
				Path: "/home/user/projects/frontend/app",
				Name: "app",
			},
			base: "/home/user/projects",
			want: "frontend/app",
		},
		{
			name: "with git",
			ws: Workspace{
				Path: "/home/user/projects/myapp",
				Name: "myapp",
			},
			base: "/home/user/projects",
			want: "myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ws.LabelWithBase(tt.base)
			if got != tt.want {
				t.Errorf("Workspace.LabelWithBase() = %v, want %v", got, tt.want)
			}
		})
	}
}
