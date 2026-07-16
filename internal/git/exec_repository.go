package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ExecRepository implements Repository by shelling out to the git binary.
type ExecRepository struct {
	Dir string
}

// NewExecRepository returns a Repository rooted at dir.
func NewExecRepository(dir string) *ExecRepository {
	return &ExecRepository{Dir: dir}
}

func (r *ExecRepository) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Dir
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out.String(), nil
}

// DiffStaged implements Repository.
func (r *ExecRepository) DiffStaged() (string, error) {
	return r.run("diff", "--cached")
}

// Status implements Repository.
func (r *ExecRepository) Status() (string, error) {
	return r.run("status", "--porcelain")
}

// Branch implements Repository.
func (r *ExecRepository) Branch() (string, error) {
	out, err := r.run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Log implements Repository.
func (r *ExecRepository) Log(n int) ([]string, error) {
	out, err := r.run("log", "-"+strconv.Itoa(n), "--pretty=format:%s")
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

// LastTag implements Repository. A repo with no tags is not an error.
func (r *ExecRepository) LastTag() (string, error) {
	out, err := r.run("describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

// ChangedFiles implements Repository.
func (r *ExecRepository) ChangedFiles() ([]string, error) {
	out, err := r.run("diff", "--cached", "--name-only")
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return nil, nil
	}
	return strings.Split(trimmed, "\n"), nil
}
