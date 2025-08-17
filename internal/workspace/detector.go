package workspace

import (
	"os"
	"path/filepath"
)

var packageFiles = []string{
	"package.json",     // Node.js
	"go.mod",           // Go
	"Cargo.toml",       // Rust
	"pom.xml",          // Maven
	"build.gradle",     // Gradle
	"build.gradle.kts", // Gradle Kotlin
	"pyproject.toml",   // Python
	"requirements.txt", // Python
	"Gemfile",          // Ruby
	"composer.json",    // PHP
	"pubspec.yaml",     // Dart/Flutter
	"mix.exs",          // Elixir
	"Makefile",         // Generic
	"CMakeLists.txt",   // CMake
	".clang-format",    // C/C++
	"stack.yaml",       // Haskell
}

func IsWorkspace(dir string) bool {
	// Check for .git directory
	gitPath := filepath.Join(dir, ".git")
	if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
		return true
	}

	// Check for package/project files
	for _, file := range packageFiles {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func HasGitRepo(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

func HasPackageFile(dir string) bool {
	for _, file := range packageFiles {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

func GetPackageType(dir string) string {
	packageTypes := map[string]string{
		"package.json":     "node",
		"go.mod":           "go",
		"Cargo.toml":       "rust",
		"pom.xml":          "maven",
		"build.gradle":     "gradle",
		"build.gradle.kts": "gradle",
		"pyproject.toml":   "python",
		"requirements.txt": "python",
		"Gemfile":          "ruby",
		"composer.json":    "php",
		"pubspec.yaml":     "dart",
		"mix.exs":          "elixir",
		"stack.yaml":       "haskell",
	}

	for file, packageType := range packageTypes {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return packageType
		}
	}

	if _, err := os.Stat(filepath.Join(dir, "Makefile")); err == nil {
		return "make"
	}

	if _, err := os.Stat(filepath.Join(dir, "CMakeLists.txt")); err == nil {
		return "cmake"
	}

	return ""
}
