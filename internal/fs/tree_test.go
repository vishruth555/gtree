package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanBuildsRecursiveSizesAndSortsBySize(t *testing.T) {
	root := t.TempDir()

	mustWriteSizedFile(t, filepath.Join(root, "small.txt"), 10)
	mustWriteSizedFile(t, filepath.Join(root, "large.txt"), 30)

	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	mustWriteSizedFile(t, filepath.Join(nested, "inside.bin"), 20)

	result, err := Scan(root)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if got, want := result.Root.Size, int64(60); got != want {
		t.Fatalf("root size = %d, want %d", got, want)
	}

	if len(result.Root.Children) != 3 {
		t.Fatalf("children count = %d, want 3", len(result.Root.Children))
	}

	gotOrder := []string{
		result.Root.Children[0].Name,
		result.Root.Children[1].Name,
		result.Root.Children[2].Name,
	}
	wantOrder := []string{"large.txt", "nested", "small.txt"}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Fatalf("child %d = %q, want %q", i, gotOrder[i], wantOrder[i])
		}
	}
}

func TestScanSkipsSymlinkTargets(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.txt")
	mustWriteSizedFile(t, target, 25)

	link := filepath.Join(root, "target-link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	result, err := Scan(root)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Fatal("expected warning for skipped symlink")
	}

	if got, want := result.Root.Size, int64(25); got != want {
		t.Fatalf("root size = %d, want %d", got, want)
	}
}

func mustWriteSizedFile(t *testing.T, path string, size int) {
	t.Helper()

	data := make([]byte, size)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}
