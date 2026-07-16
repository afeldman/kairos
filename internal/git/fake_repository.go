package git

// FakeRepository is an in-memory Repository for tests that don't want to
// shell out to git.
type FakeRepository struct {
	Diff         string
	StatusOutput string
	BranchName   string
	Commits      []string
	Tag          string
	Files        []string
}

// DiffStaged implements Repository.
func (f *FakeRepository) DiffStaged() (string, error) { return f.Diff, nil }

// Status implements Repository.
func (f *FakeRepository) Status() (string, error) { return f.StatusOutput, nil }

// Branch implements Repository.
func (f *FakeRepository) Branch() (string, error) { return f.BranchName, nil }

// Log implements Repository. It returns up to n commits from Commits.
func (f *FakeRepository) Log(n int) ([]string, error) {
	if n >= len(f.Commits) {
		return f.Commits, nil
	}
	return f.Commits[:n], nil
}

// LastTag implements Repository.
func (f *FakeRepository) LastTag() (string, error) { return f.Tag, nil }

// ChangedFiles implements Repository.
func (f *FakeRepository) ChangedFiles() ([]string, error) { return f.Files, nil }
