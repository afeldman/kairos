package context

import (
	"context"
	"testing"

	"github.com/afeldman/kairos/internal/config"
	"github.com/afeldman/kairos/internal/git"
)

func TestBuild(t *testing.T) {
	ctx := context.Background()
	repo := git.NewFakeRepo(t)
	cfg := config.Config{History: 5}

	b := NewBuilder()
	pc, err := b.Build(ctx, repo, cfg)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if pc.Branch != "main" {
		t.Fatalf("Branch = %q, want main", pc.Branch)
	}
	if len(pc.RecentCommits) != 1 || pc.RecentCommits[0] != "abc1234 init" {
		t.Fatalf("RecentCommits = %v", pc.RecentCommits)
	}
	if pc.LastTag != "v0.1.0" {
		t.Fatalf("LastTag = %q, want v0.1.0", pc.LastTag)
	}
	if len(pc.ChangedFiles) != 1 || pc.ChangedFiles[0] != "main.go" {
		t.Fatalf("ChangedFiles = %v", pc.ChangedFiles)
	}
	if pc.Diff == "" {
		t.Fatal("Diff should not be empty")
	}
}

func TestBuild_ErrorPropagation(t *testing.T) {
	ctx := context.Background()
	repo := &git.FakeRepo{
		DiffStagedFunc: func(context.Context) (string, error) {
			return "", git.ErrGitNotFound
		},
		StatusFunc:       func(context.Context) (string, error) { return "", nil },
		BranchFunc:       func(context.Context) (string, error) { return "main", nil },
		LogFunc:          func(context.Context, int) ([]string, error) { return nil, nil },
		LastTagFunc:      func(context.Context) (string, error) { return "", nil },
		DescribeFunc:     func(context.Context) (string, error) { return "", nil },
		ChangedFilesFunc: func(context.Context) ([]string, error) { return nil, nil },
	}

	b := NewBuilder()
	_, err := b.Build(ctx, repo, config.Config{})
	if err == nil {
		t.Fatal("expected error from Build()")
	}
	t.Logf("got expected error: %v", err)
}
