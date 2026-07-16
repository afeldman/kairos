# Kairos — Phase 1: Foundation — Design

## Context

Kairos is a Git Context Engine: it builds rich project/git context before asking an LLM to
generate commit messages (and later tags, releases, changelogs). This spec covers **Phase 1
only** — the foundation slice that makes `kairos` (default command) produce a Conventional
Commits message from staged changes, using Ollama as the sole provider.

Full project scope (multi-provider, tag/release, changelog engine, explain/doctor, custom
templates) is decomposed into later phases (2–5), each with its own spec/plan cycle. This
document does not cover those.

Module path: `github.com/afeldman/kairos`.

## Goals (Phase 1)

- `kairos` (no subcommand) reads **staged** changes and prints a Conventional Commits message
  to stdout, usable as `git commit -am "$(kairos)"` (after `git add`).
- Git Context Engine collects: staged diff, branch, last N commit subjects, last tag, changed
  files, detected project type, README/CHANGELOG excerpts (read-only).
- Single LLM provider: Ollama, via hand-rolled HTTP client (no extra deps).
- Prompt Builder asks the LLM for structured JSON; Formatter renders it via a Conventional
  Commits template.
- Config via Viper: `~/.config/kairos/config.yaml`, overridable by env vars and CLI flags.
- Fully unit tested without network access or a running Ollama instance.

Out of scope for Phase 1: other providers, `tag`/`release`/`changelog`/`explain`/`doctor`/
`config`/`providers` subcommands, custom templates, changelog writing.

## Package Layout

```
cmd/kairos/main.go        Cobra root command wiring

internal/
  git/                     Repository interface + exec-based implementation
  detect/                  ProjectDetector (go.mod, package.json, Cargo.toml, pyproject.toml)
  context/                 ProjectContext struct + Builder
  provider/                Provider interface + registry
  provider/ollama/         Ollama HTTP client implementing Provider
  prompt/                  Builder: ProjectContext -> LLM messages (JSON-output instructed)
  formatter/               CommitMessage struct, JSON parser, style templates
  config/                  Viper-backed Config loader

pkg/                       Empty for now — reserved for future public API surface
```

Each package has a single responsibility and communicates through small interfaces
(`git.Repository`, `provider.Provider`), so each is unit-testable in isolation via
fakes/mocks.

## Data Flow

1. `kairos` invoked → `config.Load()` resolves `Config{Provider, Model, Temperature, History,
   Style, Language}` from flags > `KAIROS_*` env vars > config file > defaults.
2. `git.NewRepository(".")` (exec-based). If `git diff --cached` is empty, return a clear
   error: `"nothing staged; run git add first"`.
3. `context.Builder.Build(repo, detector, cfg)` assembles:
   ```go
   type ProjectContext struct {
       Diff            string
       Branch          string
       RecentCommits   []string // last cfg.History subjects
       LastTag         string
       ChangedFiles    []string
       ProjectType     string   // "go", "node", "rust", "python", "unknown"
       ReadmeExcerpt   string   // first N lines if README.md exists
       ChangelogExcerpt string  // first N lines if CHANGELOG.md exists
   }
   ```
4. `prompt.Build(ctx, task="commit", cfg)` produces system + user messages. The system
   message instructs the LLM to respond with **only** JSON:
   `{"type": "...", "scope": "...", "subject": "...", "body": "...", "breaking": "..."}`.
5. `provider.Registry.Get(cfg.Provider)` returns the `ollama.Client`;
   `client.Generate(ctx, messages, cfg.Model, cfg.Temperature)` returns the raw LLM text.
6. `formatter.Parse(raw)` → `CommitMessage`. If the response isn't valid JSON, fall back:
   first line becomes `Subject`, remaining lines become `Body`, `Type`/`Scope` left empty.
   No retry call to the LLM (kept simple/cheap for Phase 1; can be revisited later).
7. `formatter.Render("conventional", msg)` → final string, printed to stdout (no trailing
   newline issues — exactly what `git commit -m` expects).

## Key Interfaces

```go
// internal/git
type Repository interface {
    DiffStaged() (string, error)
    Status() (string, error)
    Branch() (string, error)
    Log(n int) ([]string, error)
    LastTag() (string, error)
    ChangedFiles() ([]string, error)
}
```

```go
// internal/provider
type Message struct {
    Role    string // "system" | "user"
    Content string
}

type Request struct {
    Model       string
    Temperature float64
    Messages    []Message
}

type Provider interface {
    Name() string
    Generate(ctx context.Context, req Request) (string, error)
}
```

Formatter styles are implemented as Go `text/template` templates keyed by style name; Phase 1
ships only `"conventional"`. The `Style` lookup is a `map[string]*template.Template` so adding
`github`/`semantic-release`/`custom` later is additive, not a rewrite.

## Config

`~/.config/kairos/config.yaml`:
```yaml
provider: ollama
model: qwen3:30b
temperature: 0.2
history: 20
style: conventional
language: english
update_changelog: true   # parsed now, unused until Phase 4
```
Viper precedence: CLI flags > `KAIROS_*` env vars > config file > built-in defaults. Config
package exposes `config.Load() (Config, error)` — no global state; `Config` is passed
explicitly through the call chain.

## Error Handling

- No staged changes → explicit error, no LLM call made.
- Ollama unreachable/non-200 → wrapped error naming the provider and underlying cause (e.g.
  `"ollama: connection refused (is ollama running? try 'ollama serve')"`). `kairos doctor` is
  Phase 2; Phase 1's error message alone must be enough for a user to self-diagnose.
- Malformed JSON from LLM → text-fallback parse (above), not a hard failure.
- Git command failures (not a repo, git not installed) → wrapped error, no panic.

## Testing Strategy

- `git`: `Repository` interface; exec-based impl tested against a real temp git repo fixture
  created in test setup (local `git` binary only, no network); context/prompt tests use a fake
  in-memory `Repository`.
- `provider`: fake `Provider` returning canned JSON for prompt/formatter integration tests;
  `ollama.Client` tested against `httptest.Server` mocking the Ollama REST API.
- `formatter`: table-driven tests covering valid JSON, malformed JSON (fallback path), and
  template rendering per field (including empty `scope`/`breaking`).
- `config`: precedence tests using `t.TempDir()` config files + `t.Setenv()` for env override
  cases.
- `detect`: table-driven tests with fixture directories containing each manifest file
  (and one with none → `"unknown"`).

No test requires network access or a running Ollama instance.
