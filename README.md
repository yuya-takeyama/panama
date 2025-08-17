# Panama

Fast workspace finder and switcher with built-in fuzzy finder.

## Features

- 🚀 **Fast workspace detection** - Automatically finds Git repositories and project directories
- 🔍 **Built-in fuzzy finder** - No external dependencies like `fzf` or `peco` required
- 🌍 **Cross-platform** - Works on macOS, Linux, and Windows
- 📁 **Smart detection** - Recognizes projects by package files (package.json, go.mod, Cargo.toml, etc.)
- 🎯 **Intelligent scoring** - Prioritizes workspaces based on depth, type, and query matching
- 🤖 **CI/CD friendly** - Fallback to non-interactive mode for automation

## Installation

### From source

```bash
go install github.com/yuya-takeyama/panama/cmd/panama@latest
```

### Pre-built binaries

Download from [releases page](https://github.com/yuya-takeyama/panama/releases).

## Usage

### Interactive selection

```bash
# Select from current directory
panama select

# Select from specific directory
panama select ~/projects

# Start with initial query
panama select -q api

# Output as cd command
panama select -f cd
```

### Non-interactive mode

```bash
# Select first result
panama select --first

# Use in scripts
cd "$(panama select --first)"
```

### List workspaces

```bash
# List all workspaces
panama list

# Output as JSON
panama list -f json

# Limit search depth
panama list --max-depth 2
```

### Initialize configuration

```bash
# Create .panama.yaml
panama init

# Create with different format
panama init --format toml
```

## Configuration

Panama looks for configuration files in the following order:
- `.panama.yaml` / `.panama.yml`

The configuration file is searched upward from the current directory.

### Example configuration

```yaml
# UI mode: fuzzyfinder or stdio
ui: fuzzyfinder

# Maximum search depth
max_depth: 3

# Default output format: path, cd, or json
format: path

# Directories to ignore
ignore_dirs:
  - node_modules
  - .git
  - vendor
  - target
  - dist
  - build

# Workspace scoring configuration
score:
  recent_access_weight: 0.5
  frequency_weight: 0.3
  depth_penalty: 0.1
```

## Workspace Detection

Panama detects workspaces by looking for:

### Version Control
- `.git` directories

### Package/Project Files
- `package.json` (Node.js)
- `go.mod` (Go)
- `Cargo.toml` (Rust)
- `pom.xml` (Maven)
- `build.gradle` / `build.gradle.kts` (Gradle)
- `pyproject.toml` / `requirements.txt` (Python)
- `Gemfile` (Ruby)
- `composer.json` (PHP)
- `pubspec.yaml` (Dart/Flutter)
- `mix.exs` (Elixir)
- `Makefile`
- `CMakeLists.txt` (CMake)
- `stack.yaml` (Haskell)

## Shell Integration

### Bash/Zsh

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
# Quick workspace switcher
pw() {
  local dir
  dir=$(panama select "$@")
  if [[ -n "$dir" ]]; then
    cd "$dir"
  fi
}
```

### Fish

Add to your `~/.config/fish/functions/pw.fish`:

```fish
function pw
  set -l dir (panama select $argv)
  if test -n "$dir"
    cd $dir
  end
end
```

## Environment Variables

- `PANAMA_CONFIG` - Path to configuration file
- `PANAMA_UI` - Override UI mode (fuzzyfinder/stdio)

## Development

```bash
# Clone repository
git clone https://github.com/yuya-takeyama/panama.git
cd panama

# Build
go build -o panama cmd/panama/*.go

# Run tests
go test ./...

# Install locally
go install ./cmd/panama
```

## License

MIT
