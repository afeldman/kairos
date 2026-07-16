package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// InitTestRepo creates a temporary bare-minimum git repository and returns
// its root directory. The caller is responsible for cleanup (t.TempDir).
func InitTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmds := []*exec.Cmd{
		exec.Command("git", "init"),
		exec.Command("git", "config", "user.email", "test@test"),
		exec.Command("git", "config", "user.name", "Test"),
	}
	for _, cmd := range cmds {
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git init setup: %v\n%s", err, out)
		}
	}
	return dir
}

// StageAndCommit creates a file, adds it to the index, and commits.
func StageAndCommit(t *testing.T, dir, filename, content, msg string) {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", filename, err)
	}
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", msg)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
}

// Stage creates a file and adds it to the index without committing.
func Stage(t *testing.T, dir, filename, content string) {
	t.Helper()
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", filename, err)
	}
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
}

func TestExecRepo_Basic(t *testing.T) {
	dir := InitTestRepo(t)
	StageAndCommit(t, dir, "README.md", "# Hello", "initial commit")
	Stage(t, dir, "main.go", "package main\nfunc main() {}")

	repo := NewRepository(dir)

	ctx := context.Background()

	t.Run("Branch", func(t *testing.T) {
		b, err := repo.Branch(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if b != "master" && b != "main" {
			t.Fatalf("Branch() = %q, want master or main", b)
		}
	})

	t.Run("Log", func(t *testing.T) {
		log, err := repo.Log(ctx, 5)
		if err != nil {
			t.Fatal(err)
		}
		if len(log) < 1 {
			t.Fatal("Log() returned empty")
		}
		if !strings.Contains(log[0], "initial commit") {
			t.Fatalf("Log()[0] = %q, want 'initial commit'", log[0])
		}
	})

	t.Run("DiffStaged", func(t *testing.T) {
		diff, err := repo.DiffStaged(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(diff, "func main()") {
			t.Fatalf("DiffStaged() = %q, want func main()", diff)
		}
	})

	t.Run("ChangedFiles", func(t *testing.T) {
		files, err := repo.ChangedFiles(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 1 || files[0] != "main.go" {
			t.Fatalf("ChangedFiles() = %v, want [main.go]", files)
		}
	})

	t.Run("LastTag", func(t *testing.T) {
		tag, err := repo.LastTag(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if tag != "" {
			t.Fatalf("LastTag() = %q, want empty (no tags)", tag)
		}
	})

	t.Run("Describe", func(t *testing.T) {
		desc, err := repo.Describe(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if desc != "" {
			t.Fatalf("Describe() = %q, want empty", desc)
		}
	})
}

func TestEnsureStaged_Empty(t *testing.T) {
	dir := InitTestRepo(t)
	repo := NewRepository(dir)
	err := EnsureStaged(context.Background(), repo)
	if err == nil {
		t.Fatal("expected error for empty staged changes")
	}
}

func TestEnsureStaged_HasChanges(t *testing.T) {
	dir := InitTestRepo(t)
	Stage(t, dir, "foo.go", "package foo")
	repo := NewRepository(dir)
	if err := EnsureStaged(context.Background(), repo); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFirstLine(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"hello\nworld", "hello"},
		{"  hello  ", "hello"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := FirstLine(tt.input); got != tt.want {
			t.Errorf("FirstLine(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
