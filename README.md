# gtree

`gtree` is a terminal UI (TUI) disk usage visualizer written in Go. It scans a directory recursively, calculates file and folder sizes, and lets you navigate the directory tree interactively from your terminal.


## Features

- Recursive directory scanning
- Files and folders ranked by size
- Interactive navigation through directories
- Fast and lightweight implementation in Go
- Runs entirely inside the terminal

## Demo Workflow

1. Start `gtree`.
2. Let it scan the target directory recursively.
3. View files and folders sorted by size.
4. Navigate with the keyboard controls.
5. Drill into large directories.
6. Go back to parent folders when needed.

## Installation

### Option 1: Install with Go (recommended)

```bash
go install github.com/vishruth555/gtree/cmd/gtree@latest
```

Then run:

```bash
gtree
```
this will work anywhere on your terminal!

### Option 2: Build locally

From this repository:

```bash
go build ./cmd/gtree
```

This creates a local `gtree` binary in the project directory.

### Option 3: Run from source

```bash
go run ./cmd/gtree
```

Scan a specific directory:

```bash
go run ./cmd/gtree ./internal
```

## Usage

Scan the current directory:

```bash
gtree
```

Scan a specific directory:

```bash
gtree ./internal
gtree ~/Downloads
```

If no argument is provided, `gtree` scans the current working directory. If one argument is provided, it resolves that path from the current directory and scans it instead.

## Navigation Controls

| Key | Action |
| --- | --- |
| `j` or `Down` | Move selection down |
| `k` or `Up` | Move selection up |
| `Enter` | Open or focus the selected folder |
| `Backspace`, `h`, or `Left` | Go back to the parent directory |
| `g` | Jump to the top |
| `G` | Jump to the bottom |
| `q` | Quit the application |

## How It Works

`gtree` is split into a few focused parts so the code stays maintainable:

- `cmd/gtree`: the CLI entrypoint. It resolves the target directory and starts the Bubble Tea program.
- `internal/fs`: the recursive scanner. It walks the filesystem, computes sizes, sorts children by size, and records non-fatal warnings such as skipped symlinks.
- `internal/ui`: the TUI layer. It tracks the current folder, cursor position, navigation stack, and screen rendering.

## Example Use Cases

- Find large folders consuming disk space
- Explore build artifacts or log directories
- Analyze project directory sizes
- Quickly inspect disk usage during development

## Requirements

- Go 1.25 or newer for local development in this repository
- A terminal that supports interactive TUI applications
