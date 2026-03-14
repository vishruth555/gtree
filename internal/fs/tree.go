package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Node models a filesystem entry and the total space it occupies.
type Node struct {
	Name     string
	Path     string
	Size     int64
	IsDir    bool
	Children []*Node
}

// ScanResult contains the scanned tree plus non-fatal issues.
type ScanResult struct {
	Root     *Node
	Warnings []string
}

// Scan recursively crawls the provided path and returns a size tree.
func Scan(root string) (*ScanResult, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve root: %w", err)
	}

	warnings := make([]string, 0)
	node, err := scan(absRoot, &warnings)
	if err != nil {
		return nil, err
	}

	return &ScanResult{
		Root:     node,
		Warnings: warnings,
	}, nil
}

func scan(path string, warnings *[]string) (*Node, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}

	node := &Node{
		Name:  info.Name(),
		Path:  path,
		Size:  info.Size(),
		IsDir: info.IsDir(),
	}

	if info.Mode()&os.ModeSymlink != 0 {
		node.Size = 0
		*warnings = append(*warnings, fmt.Sprintf("skipped symlink: %s", path))
		return node, nil
	}

	if !info.IsDir() {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		*warnings = append(*warnings, fmt.Sprintf("could not read %s: %v", path, err))
		return node, nil
	}

	node.Size = 0
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		child, err := scan(childPath, warnings)
		if err != nil {
			*warnings = append(*warnings, err.Error())
			continue
		}

		node.Children = append(node.Children, child)
		node.Size += child.Size
	}

	sort.Slice(node.Children, func(i, j int) bool {
		left := node.Children[i]
		right := node.Children[j]

		if left.Size != right.Size {
			return left.Size > right.Size
		}
		if left.IsDir != right.IsDir {
			return left.IsDir
		}
		return left.Name < right.Name
	})

	return node, nil
}
