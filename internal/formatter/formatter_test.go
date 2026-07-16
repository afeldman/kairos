package formatter

import (
	"testing"
)

func TestParse_JSON(t *testing.T) {
	raw := `{"type":"feat","scope":"core","subject":"add ollama support","body":"This adds the ollama provider client.","breaking":""}`
	msg := Parse(raw)

	if msg.Type != "feat" {
		t.Fatalf("Type = %q, want feat", msg.Type)
	}
	if msg.Scope != "core" {
		t.Fatalf("Scope = %q, want core", msg.Scope)
	}
	if msg.Subject != "add ollama support" {
		t.Fatalf("Subject = %q, want 'add ollama support'", msg.Subject)
	}
	if msg.Body != "This adds the ollama provider client." {
		t.Fatalf("Body = %q, want 'This adds the ollama provider client.'", msg.Body)
	}
	if msg.Breaking != "" {
		t.Fatalf("Breaking = %q, want empty", msg.Breaking)
	}
}

func TestParse_JSONWithBreaking(t *testing.T) {
	raw := `{"type":"feat!","scope":"api","subject":"redesign auth","body":"Replaced OAuth 1.0 with OAuth 2.0","breaking":"OAuth 1.0 tokens no longer accepted"}`
	msg := Parse(raw)

	if msg.Type != "feat!" {
		t.Fatalf("Type = %q", msg.Type)
	}
	if msg.Breaking != "OAuth 1.0 tokens no longer accepted" {
		t.Fatalf("Breaking = %q", msg.Breaking)
	}
}

func TestParse_Fallback(t *testing.T) {
	raw := "fix typo in README\n\nCorrected the spelling of 'kairos'."
	msg := Parse(raw)

	if msg.Type != "" {
		t.Fatalf("Type = %q, want empty on fallback", msg.Type)
	}
	if msg.Subject != "fix typo in README" {
		t.Fatalf("Subject = %q", msg.Subject)
	}
	if msg.Body != "Corrected the spelling of 'kairos'." {
		t.Fatalf("Body = %q", msg.Body)
	}
}

func TestParse_FallbackSingleLine(t *testing.T) {
	raw := "fix typo in README"
	msg := Parse(raw)

	if msg.Subject != "fix typo in README" {
		t.Fatalf("Subject = %q", msg.Subject)
	}
	if msg.Body != "" {
		t.Fatalf("Body = %q, want empty", msg.Body)
	}
}

func TestParse_InvalidJSONButHasFields(t *testing.T) {
	// Missing "type" field but still valid JSON — falls back to text mode.
	raw := `{"subject":"only a subject"}`
	msg := Parse(raw)
	// In text-fallback mode, the first line of the raw JSON becomes the subject.
	if msg.Subject != `{"subject":"only a subject"}` {
		t.Fatalf("Subject = %q, want raw JSON as subject", msg.Subject)
	}
	if msg.Type != "" {
		t.Fatalf("Type = %q, want empty", msg.Type)
	}
}

func TestRender_Conventional(t *testing.T) {
	msg := CommitMessage{
		Type:    "feat",
		Scope:   "core",
		Subject: "add ollama support",
		Body:    "This adds the ollama provider.",
	}

	got, err := Render("conventional", msg)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want := "feat(core): add ollama support\n\nThis adds the ollama provider."
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRender_ConventionalWithBreaking(t *testing.T) {
	msg := CommitMessage{
		Type:     "feat!",
		Scope:    "api",
		Subject:  "redesign auth",
		Body:     "Replaced OAuth 1.0 with OAuth 2.0",
		Breaking: "OAuth 1.0 tokens no longer accepted",
	}

	got, err := Render("conventional", msg)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want := "feat!(api): redesign auth\n\nReplaced OAuth 1.0 with OAuth 2.0\n\nBREAKING CHANGE: OAuth 1.0 tokens no longer accepted"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRender_ConventionalNoScope(t *testing.T) {
	msg := CommitMessage{
		Type:    "fix",
		Subject: "correct typo in README",
	}

	got, err := Render("conventional", msg)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want := "fix: correct typo in README"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRender_UnknownStyle(t *testing.T) {
	_, err := Render("nonexistent", CommitMessage{})
	if err == nil {
		t.Fatal("expected error for unknown style")
	}
}

func TestRoundTrip(t *testing.T) {
	// Full pipeline: parse JSON → render conventional.
	raw := `{"type":"feat","scope":"","subject":"add ci pipeline","body":"Add GitHub Actions workflow.","breaking":""}`
	msg := Parse(raw)
	got, err := Render("conventional", msg)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want := "feat: add ci pipeline\n\nAdd GitHub Actions workflow."
	if got != want {
		t.Fatalf("RoundTrip = %q, want %q", got, want)
	}
}
