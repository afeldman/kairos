// Package git provides a Repository interface and an exec-based implementation
// for collecting Git context. The interface makes it easy to swap in fakes
// during testing.
package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// allowedGitArgs is the set of git subcommands and flags that the
// Repository methods are permitted to pass to the git binary.
// This prevents argument injection through the exec.CommandContext call.
var allowedGitArgs = []string{
	"diff", "--cached", "--name-only", "--no-color",
	"status", "--short",
	"rev-parse", "--abbrev-ref", "--is-inside-work-tree", "HEAD",
	"log", "--oneline",
	"describe", "--tags", "--abbrev=0", "--long",
}

// isAllowedGitArg reports whether arg is in the allowlist.
func isAllowedGitArg(arg string) bool {
	// In Go ≥1.26, slices.Contains handles exact string matches.
	// Flag-like args (e.g. "-d", "--name-only") are matched literally.
	return slices.Contains(allowedGitArgs, arg)
}

// ErrNothingStaged indicates there are no staged changes to build a commit
// message from. Callers can wrap it with a localized message.
var ErrNothingStaged = errors.New("nothing staged")

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
	// Validate every argument against an allowlist to prevent command injection.
	for _, arg := range args {
		// Allow numeric-only args like "-5" used for git log -N.
		if isNumericArg(arg) {
			continue
		}
		if !isAllowedGitArg(arg) {
			return "", fmt.Errorf("git: disallowed argument %q", arg)
		}
	}
	cmd := exec.CommandContext(ctx, "git", args...) // #nosec G204 — args validated above via allowlist
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

// isNumericArg reports whether arg is a "-N" style flag where N is a
// non-negative integer (e.g. "-5", "-10"). These are used for git log -N.
func isNumericArg(arg string) bool {
	if len(arg) < 2 || arg[0] != '-' {
		return false
	}
	for _, c := range arg[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
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

// IsRepo reports whether dir is inside a Git working tree.
func IsRepo(ctx context.Context, dir string) bool {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// EnsureStaged returns an error if there is nothing staged for commit.
func EnsureStaged(ctx context.Context, repo Repository) error {
	diff, err := repo.DiffStaged(ctx)
	if err != nil {
		return fmt.Errorf("check staged changes: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return ErrNothingStaged
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
