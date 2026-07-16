package git

import (
	"context"
	"testing"
)

// FakeRepo implements Repository with in-memory data for testing.
type FakeRepo struct {
	DiffStagedFunc  func(context.Context) (string, error)
	StatusFunc      func(context.Context) (string, error)
	BranchFunc      func(context.Context) (string, error)
	LogFunc         func(context.Context, int) ([]string, error)
	LastTagFunc     func(context.Context) (string, error)
	DescribeFunc    func(context.Context) (string, error)
	ChangedFilesFunc func(context.Context) ([]string, error)
}

func (f *FakeRepo) DiffStaged(ctx context.Context) (string, error) {
	return f.DiffStagedFunc(ctx)
}
func (f *FakeRepo) Status(ctx context.Context) (string, error) {
	return f.StatusFunc(ctx)
}
func (f *FakeRepo) Branch(ctx context.Context) (string, error) {
	return f.BranchFunc(ctx)
}
func (f *FakeRepo) Log(ctx context.Context, n int) ([]string, error) {
	return f.LogFunc(ctx, n)
}
func (f *FakeRepo) LastTag(ctx context.Context) (string, error) {
	return f.LastTagFunc(ctx)
}
func (f *FakeRepo) Describe(ctx context.Context) (string, error) {
	return f.DescribeFunc(ctx)
}
func (f *FakeRepo) ChangedFiles(ctx context.Context) ([]string, error) {
	return f.ChangedFilesFunc(ctx)
}

// NewFakeRepo returns a FakeRepo pre-populated with sensible no-op defaults.
func NewFakeRepo(t *testing.T) *FakeRepo {
	t.Helper()
	return &FakeRepo{
		DiffStagedFunc:  func(context.Context) (string, error) { return "+func main() {}", nil },
		StatusFunc:      func(context.Context) (string, error) { return "M  main.go", nil },
		BranchFunc:      func(context.Context) (string, error) { return "main", nil },
		LogFunc:         func(context.Context, int) ([]string, error) { return []string{"abc1234 init"}, nil },
		LastTagFunc:     func(context.Context) (string, error) { return "v0.1.0", nil },
		DescribeFunc:    func(context.Context) (string, error) { return "v0.1.0-5-gabc1234", nil },
		ChangedFilesFunc: func(context.Context) ([]string, error) { return []string{"main.go"}, nil },
	}
}
