package prompt

import (
	"testing"

	"github.com/afeldman/kairos/internal/config"
	"github.com/afeldman/kairos/internal/context"
)

func TestBuild_ContainsStagedDiff(t *testing.T) {
	pc := context.ProjectContext{
		Diff:          "+func main() {}",
		Branch:        "main",
		RecentCommits: []string{"abc1234 initial"},
		ChangedFiles:  []string{"main.go"},
		LastTag:       "v0.1.0",
		Describe:      "v0.1.0-5-gabc1234",
		ProjectType:   "go",
	}
	cfg := config.Config{
		Language: "english",
		Style:    "conventional",
	}

	msgs := Build(pc, cfg)
	if len(msgs) != 2 {
		t.Fatalf("Build() returned %d messages, want 2", len(msgs))
	}

	sys := msgs[0]
	user := msgs[1]

	if sys.Role != "system" {
		t.Fatalf("first message role = %q, want system", sys.Role)
	}
	if user.Role != "user" {
		t.Fatalf("second message role = %q, want user", user.Role)
	}

	if !contains(user.Content, "<staged-diff>") {
		t.Fatal("user prompt missing <staged-diff> section")
	}
	if !contains(user.Content, "+func main() {}") {
		t.Fatal("user prompt missing diff content")
	}
	if !contains(user.Content, "<branch>main</branch>") {
		t.Fatal("user prompt missing branch")
	}
	if !contains(user.Content, "<project-type>go</project-type>") {
		t.Fatal("user prompt missing project type")
	}
	if !contains(user.Content, "<recent-commits>") {
		t.Fatal("user prompt missing recent commits")
	}
	if !contains(user.Content, "<changed-files>") {
		t.Fatal("user prompt missing changed files")
	}
	if !contains(sys.Content, "Conventional Commits") {
		t.Fatal("system prompt missing Conventional Commits mention")
	}
}

func TestBuild_NoOptionalFields(t *testing.T) {
	pc := context.ProjectContext{
		Diff:          "some diff",
		RecentCommits: nil,
		ChangedFiles:  nil,
	}
	cfg := config.Config{Language: "english"}

	msgs := Build(pc, cfg)
	user := msgs[1]

	if contains(user.Content, "<recent-commits>") {
		t.Fatal("should not include recent-commits section when empty")
	}
	if contains(user.Content, "unknown") {
		t.Log("unknown project type may appear, but is harmless")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
