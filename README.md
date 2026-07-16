# Kairos

**Kairos** is a **Git Context Engine** that understands a project's history to generate high-quality commit messages, tag messages, release notes, and changelog entries using LLMs.

> The best commit message is written by understanding the project's history.

## Philosophy

Kairos is **not** just an AI commit message generator. It builds rich context from your Git repository — history, branches, tags, changed files, project type — before asking any LLM for output.

## Features (Phase 1)

- 🔍 **Git Context Engine** — collects staged diff, branch, recent commits, last tag, changed files, project type, README/CHANGELOG excerpts
- 🤖 **LLM Providers** — Ollama (default), with pluggable interface for OpenAI, Anthropic, Gemini, OpenRouter, LM Studio
- 📝 **Conventional Commits** — structured JSON output parsed into type, scope, subject, body, and breaking changes
- ⚙️ **Configurable** — `~/.config/kairos/config.yaml`, `KAIROS_*` environment variables, CLI flags
- 🧪 **Fully Tested** — unit tests with mocked Git and provider dependencies, no network required

## Quick Start

### Prerequisites

- Go 1.26+
- Git
- [Ollama](https://ollama.com/) (default provider) with a model like `qwen3:30b`

### Install

```bash
go install github.com/afeldman/kairos/cmd/kairos@latest
```

### Usage

**Generate a commit message:**

```bash
git add .
git commit -am "$(kairos)"
```

**Or using `kairos commit`:**

```bash
git add .
kairos commit
```

**With a different model:**

```bash
kairos -m llama3.2:3b
```

### Configuration

Create `~/.config/kairos/config.yaml`:

```yaml
provider: ollama
model: qwen3:30b
temperature: 0.2
history: 20
style: conventional
language: english
update_changelog: true
```

## Commands

| Command | Description |
|---------|-------------|
| `kairos` | Generate commit message (default) |
| `kairos commit` | Generate commit message |
| `kairos tag` | Generate annotated tag message |
| `kairos release` | Generate release notes |
| `kairos changelog` | Update or generate CHANGELOG.md |
| `kairos explain` | Explain staged or committed changes |
| `kairos doctor` | Diagnose configuration and environment |
| `kairos config` | Display or edit configuration |
| `kairos providers` | List available LLM providers |

## Architecture

```
cmd/kairos/main.go        Entry point
internal/
  cmd/                    Cobra command definitions
  config/                 Viper-backed configuration
  context/                Project context builder (core engine)
  detect/                 Project type detection
  formatter/              Commit message parser and renderer
  git/                    Git repository abstraction
  prompt/                 LLM prompt builder
  provider/               Provider interface + registry
  provider/ollama/        Ollama HTTP client
pkg/                      Reserved for future public API
```

## Roadmap

- [x] Phase 1 — Core CLI, Git context, Ollama provider, Conventional Commits
- [ ] Phase 2 — Multi-provider (OpenAI, Anthropic, Gemini, LM Studio, OpenRouter)
- [ ] Phase 3 — `kairos tag`, `kairos release`, `kairos changelog`
- [ ] Phase 4 — `kairos explain`, `kairos doctor`
- [ ] Phase 5 — GitHub/GitLab integration, Jira, Linear, MCP support

## License

Apache 2.0
