package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func touch(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("{}"), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name string
		file string
		want ProjectType
	}{
		{"go", "go.mod", TypeGo},
		{"node", "package.json", TypeNode},
		{"rust", "Cargo.toml", TypeRust},
		{"python", "pyproject.toml", TypePython},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			touch(t, dir, tt.file)

			if got := Detect(dir); got != tt.want {
				t.Fatalf("Detect() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetect_Unknown(t *testing.T) {
	dir := t.TempDir()

	if got := Detect(dir); got != TypeUnknown {
		t.Fatalf("Detect() = %q, want %q", got, TypeUnknown)
	}
}
