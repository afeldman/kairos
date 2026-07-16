package git

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// newTestRepo creates a throwaway git repo with one commit and one
// staged change, and returns its path.
func newTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(cmd.Env,
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
			"HOME="+dir,
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "test")

	readme := filepath.Join(dir, "README.md")
	if err := writeFile(readme, "# test repo\n"); err != nil {
		t.Fatalf("write README: %v", err)
	}
	run("add", "README.md")
	run("commit", "-m", "chore: initial commit")

	feature := filepath.Join(dir, "feature.go")
	if err := writeFile(feature, "package main\n"); err != nil {
		t.Fatalf("write feature.go: %v", err)
	}
	run("add", "feature.go")

	return dir
}

func TestExecRepository_DiffStaged(t *testing.T) {
	dir := newTestRepo(t)
	repo := NewExecRepository(dir)

	diff, err := repo.DiffStaged()
	if err != nil {
		t.Fatalf("DiffStaged() error = %v", err)
	}
	if !strings.Contains(diff, "feature.go") {
		t.Fatalf("DiffStaged() = %q, want it to mention feature.go", diff)
	}
}

func TestExecRepository_Branch(t *testing.T) {
	dir := newTestRepo(t)
	repo := NewExecRepository(dir)

	branch, err := repo.Branch()
	if err != nil {
		t.Fatalf("Branch() error = %v", err)
	}
	if branch == "" {
		t.Fatal("Branch() = \"\", want a branch name")
	}
}

func TestExecRepository_Log(t *testing.T) {
	dir := newTestRepo(t)
	repo := NewExecRepository(dir)

	log, err := repo.Log(5)
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
	if len(log) != 1 || log[0] != "chore: initial commit" {
		t.Fatalf("Log() = %v, want [\"chore: initial commit\"]", log)
	}
}

func TestExecRepository_LastTag_NoTags(t *testing.T) {
	dir := newTestRepo(t)
	repo := NewExecRepository(dir)

	tag, err := repo.LastTag()
	if err != nil {
		t.Fatalf("LastTag() error = %v, want nil (no tags is not an error)", err)
	}
	if tag != "" {
		t.Fatalf("LastTag() = %q, want empty string", tag)
	}
}

func TestExecRepository_ChangedFiles(t *testing.T) {
	dir := newTestRepo(t)
	repo := NewExecRepository(dir)

	files, err := repo.ChangedFiles()
	if err != nil {
		t.Fatalf("ChangedFiles() error = %v", err)
	}
	if len(files) != 1 || files[0] != "feature.go" {
		t.Fatalf("ChangedFiles() = %v, want [\"feature.go\"]", files)
	}
}
