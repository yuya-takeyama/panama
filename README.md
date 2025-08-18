# Panama

Go ahead and jump into directories in a monorepo

## Features

- 🚀 **Fast workspace detection** - Automatically finds Git repositories and project directories
- 🔍 **Built-in fuzzy finder** - No external dependencies like `fzf` required
- 🌍 **Cross-platform** - Works on macOS, Linux, and Windows
- 📁 **Smart detection** - Recognizes projects by package files (package.json, go.mod, pyproject.toml, etc.)
- ⚙️ **Customizable patterns** - Configure your own workspace detection patterns

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
```

## Configuration

Panama looks for configuration files in the following order:
- `.panama.yaml` / `.panama.yml`

The configuration file is searched upward from the current directory. When a configuration file is found, Panama uses that directory as the search root.

### Example configuration

```yaml
# Maximum search depth
max_depth: 6

# Default output format: path, cd, or json
format: path

# Workspace detection patterns
# No defaults - configure based on your needs
patterns:
  - package.json
  - go.mod
  - pyproject.toml
  # Add more patterns as needed:
  # - Cargo.toml        # Rust
  # - "*.xcodeproj"     # Xcode project

# Directories to ignore during search
# No defaults - configure based on your needs
ignored_dirs:
  - node_modules
  - .git
  - vendor
  - target
  - dist
  - build
  - .next
  - .nuxt
  - .cache
  - __pycache__
```

## Workspace Detection

Panama detects workspaces by looking for:

### Version Control
- `.git` directories (always detected)

### Package Files
Without configuration, Panama only detects `.git` directories as workspaces.
To detect package files, configure patterns in your `.panama.yaml`:

```yaml
patterns:
  - package.json    # Node.js
  - go.mod          # Go
  - pyproject.toml  # Python
  - Cargo.toml      # Rust
  - "*.xcodeproj"   # Xcode projects (glob pattern)
```

Run `panama init` to create a configuration file with commonly used patterns.

## Shell Integration

### Bash/Zsh

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
# Quick workspace switcher
# You can name this function whatever you like (e.g., j, pw, workspace, etc.)
jump() {
  local dir
  dir=$(panama select "$@")
  if [[ -n "$dir" ]]; then
    cd "$dir"
  fi
}
```

### Fish

Add to your `~/.config/fish/functions/jump.fish`:

```fish
# You can name this function whatever you like (e.g., j, pw, workspace, etc.)
function jump
  set -l dir (panama select $argv)
  if test -n "$dir"
    cd $dir
  end
end
```

## Environment Variables

- `PANAMA_CONFIG` - Path to configuration file

## Keyboard Shortcuts (Interactive Mode)

- `↑`/`↓` or `Ctrl+P`/`Ctrl+N` - Navigate through workspaces
- `Enter` - Select current workspace
- `Ctrl+C` or `Esc` - Cancel selection
- Type to filter workspaces in real-time

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

## Requirements

- Go 1.25 or later (for building from source)
- No external runtime dependencies

## Name Origin

Panama stands for **P**roject-**A**ware **N**avigator **A**cross **M**onorepo **A**pps.

### See also

- https://www.youtube.com/watch?v=fuKDBPw8wQA
- https://www.youtube.com/watch?v=SwYN7mTi6HM

## License

MIT
