// Package git provides a small abstraction over the git commands Kairos
// needs, plus an in-memory fake for tests that don't want to shell out.
package git

// Repository is the subset of git operations Kairos' context engine needs.
type Repository interface {
	// DiffStaged returns `git diff --cached` output.
	DiffStaged() (string, error)
	// Status returns `git status --porcelain` output.
	Status() (string, error)
	// Branch returns the current branch name.
	Branch() (string, error)
	// Log returns up to n recent commit subjects, newest first.
	Log(n int) ([]string, error)
	// LastTag returns the most recent tag reachable from HEAD, or "" if
	// there is none. A missing tag is not an error.
	LastTag() (string, error)
	// ChangedFiles returns the paths of staged files.
	ChangedFiles() ([]string, error)
}
