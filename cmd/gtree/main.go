package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vishruth555/gtree/internal/ui"
)

func main() {
	root, err := resolveRoot(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	program := tea.NewProgram(
		ui.New(root),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func resolveRoot(args []string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine working directory: %w", err)
	}

	switch len(args) {
	case 0:
		return cwd, nil
	case 1:
		target := args[0]
		if !filepath.IsAbs(target) {
			target = filepath.Join(cwd, target)
		}

		resolved, err := filepath.Abs(target)
		if err != nil {
			return "", fmt.Errorf("could not resolve path %q: %w", args[0], err)
		}

		return resolved, nil
	default:
		return "", fmt.Errorf("usage: gtree [relative-path]")
	}
}
