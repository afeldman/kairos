// Package context builds a rich ProjectContext from a git repository.
// This is the core value of Kairos — understanding project history before
// asking an LLM for output.
package context

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/afeldman/kairos/internal/config"
	"github.com/afeldman/kairos/internal/detect"
	"github.com/afeldman/kairos/internal/git"
)

// ProjectContext holds everything Kairos knows about the current state of a
// repository before calling an LLM.
type ProjectContext struct {
	Diff             string
	Branch           string
	RecentCommits    []string
	LastTag          string
	Describe         string
	ChangedFiles     []string
	ProjectType      string
	ReadmeExcerpt    string
	ChangelogExcerpt string
}

// Builder collects and assembles a ProjectContext from a Repository.
type Builder struct{}

// NewBuilder returns a ready-to-use Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// Build gathers all context from the repository and returns a populated
// ProjectContext.
func (b *Builder) Build(ctx context.Context, repo git.Repository, cfg config.Config) (ProjectContext, error) {
	var pc ProjectContext
	var errs []error

	diff, err := repo.DiffStaged(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("diff: %w", err))
	} else {
		pc.Diff = diff
	}

	branch, err := repo.Branch(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("branch: %w", err))
	} else {
		pc.Branch = branch
	}

	log, err := repo.Log(ctx, cfg.History)
	if err != nil {
		errs = append(errs, fmt.Errorf("log: %w", err))
	} else {
		pc.RecentCommits = log
	}

	tag, err := repo.LastTag(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("last-tag: %w", err))
	} else {
		pc.LastTag = tag
	}

	desc, err := repo.Describe(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("describe: %w", err))
	} else {
		pc.Describe = desc
	}

	files, err := repo.ChangedFiles(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("changed-files: %w", err))
	} else {
		pc.ChangedFiles = files
	}

	// Detect project type from repository root.
	pc.ProjectType = string(detect.Detect("."))

	// Read project files from CWD (best-effort).
	if cwd, err := os.Getwd(); err == nil {
		pc.ReadmeExcerpt = readFirstLines(cwd, "README.md", 10)
		pc.ChangelogExcerpt = readFirstLines(cwd, "CHANGELOG.md", 10)
	}

	if len(errs) > 0 {
		return pc, fmt.Errorf("context build: %w", errs[0])
	}
	return pc, nil
}

// readFirstLines reads up to n lines from a file inside the given root
// directory. It uses os.DirFS to scope file access and prevent directory
// traversal attacks.
func readFirstLines(root, name string, n int) string {
	fsys := os.DirFS(root)
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return ""
	}
	lines := strings.SplitN(string(data), "\n", n+1)
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
