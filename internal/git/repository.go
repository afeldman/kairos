// Package git provides a Repository interface and an exec-based implementation
// for collecting Git context. The interface makes it easy to swap in fakes
// during testing.
package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Repository abstracts a Git repository so callers can collect context
// without depending on a real git binary.
type Repository interface {
	DiffStaged(ctx context.Context) (string, error)
	Status(ctx context.Context) (string, error)
	Branch(ctx context.Context) (string, error)
	Log(ctx context.Context, n int) ([]string, error)
	LastTag(ctx context.Context) (string, error)
	Describe(ctx context.Context) (string, error)
	ChangedFiles(ctx context.Context) ([]string, error)
}

// ExecRepo is a Repository implementation that shells out to the git binary.
type ExecRepo struct {
	// Dir is the working directory for git commands.
	Dir string
}

// NewRepository returns an ExecRepo rooted at dir.
func NewRepository(dir string) *ExecRepo {
	return &ExecRepo{Dir: dir}
}

func (r *ExecRepo) git(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = r.Dir
	out, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = ee.Stderr
		}
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, string(stderr))
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func (r *ExecRepo) DiffStaged(ctx context.Context) (string, error) {
	return r.git(ctx, "diff", "--cached")
}

func (r *ExecRepo) Status(ctx context.Context) (string, error) {
	return r.git(ctx, "status", "--short")
}

func (r *ExecRepo) Branch(ctx context.Context) (string, error) {
	s, err := r.git(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return s, nil
}

func (r *ExecRepo) Log(ctx context.Context, n int) ([]string, error) {
	out, err := r.git(ctx, "log", "--oneline", fmt.Sprintf("-%d", n))
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func (r *ExecRepo) LastTag(ctx context.Context) (string, error) {
	s, err := r.git(ctx, "describe", "--tags", "--abbrev=0")
	if err != nil {
		// If there are no tags, git describe exits non-zero.
		return "", nil
	}
	return s, nil
}

func (r *ExecRepo) Describe(ctx context.Context) (string, error) {
	s, err := r.git(ctx, "describe", "--tags", "--long")
	if err != nil {
		return "", nil
	}
	return s, nil
}

func (r *ExecRepo) ChangedFiles(ctx context.Context) ([]string, error) {
	out, err := r.git(ctx, "diff", "--cached", "--name-only")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// EnsureStaged returns an error if there is nothing staged for commit.
func EnsureStaged(ctx context.Context, repo Repository) error {
	diff, err := repo.DiffStaged(ctx)
	if err != nil {
		return fmt.Errorf("check staged changes: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("nothing staged; run 'git add' first")
	}
	return nil
}

// FirstLine returns the first line of s, or the full string if there is
// only one line. Used for commit subjects.
func FirstLine(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.IndexByte(s, '\n'); idx != -1 {
		return s[:idx]
	}
	return s
}

// Buffer is a reusable buffer for building git command output.
type Buffer struct{ bytes.Buffer }

// ErrGitNotFound is a sentinel error for tests indicating git is not available.
var ErrGitNotFound = fmt.Errorf("git not found")
