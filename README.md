# Kairos

**Kairos** is a **Git Context Engine** that understands a project's history to generate high-quality commit messages, tag messages, release notes, and changelog entries using LLMs.

> The best commit message is written by understanding the project's history.

## Philosophy

Kairos is **not** just an AI commit message generator. It builds rich context from your Git repository — history, branches, tags, changed files, project type — before asking any LLM for output.

## Features (Phase 1)

- 🔍 **Git Context Engine** — collects staged diff, branch, recent commits, last tag, changed files, project type, README/CHANGELOG excerpts
- 🤖 **LLM Providers** — Ollama (default), OpenAI-compatible (OpenAI, LM Studio, [GoModel](https://github.com/ENTERPILOT/GoModel) gateway), with pluggable interface for Anthropic, Gemini, OpenRouter
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

### LM Studio (macOS)

LM Studio only exposes an OpenAI-compatible API, no native Ollama API. Start the
local server in LM Studio (default `http://localhost:1234/v1`), then:

```bash
kairos --provider lmstudio --model <loaded-model-name>
```

Or in `config.yaml`:

```yaml
provider: lmstudio
model: qwen2.5-coder-14b-instruct
# base_url: http://localhost:1234/v1   # only needed if not the default
```

### GoModel gateway (optional)

To route through a [GoModel](https://github.com/ENTERPILOT/GoModel) gateway (multi-provider
routing, cost tracking, observability) instead of talking to LM Studio directly, point Kairos
at the gateway's endpoint (default `http://localhost:8080/v1`):

```bash
kairos --provider gomodel --model gpt-5-chat-latest --api-key "$GOMODEL_KEY"
```

GoModel is a separate service (Docker/binary), not a Go dependency of Kairos — run it
alongside and configure its own `OPENAI_LMSTUDIO_BASE_URL` etc. to reach LM Studio.
`--base-url` overrides the endpoint for any of `openai`, `lmstudio`, or `gomodel` if not
running on the default port.

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
  provider/openai/        OpenAI-compatible HTTP client (openai, lmstudio, gomodel)
pkg/                      Reserved for future public API
```

## Roadmap

- [x] Phase 1 — Core CLI, Git context, Ollama provider, Conventional Commits
- [x] Phase 2a — OpenAI-compatible provider (OpenAI, LM Studio, GoModel gateway)
- [ ] Phase 2b — Anthropic, Gemini, OpenRouter
- [ ] Phase 3 — `kairos tag`, `kairos release`, `kairos changelog`
- [ ] Phase 4 — `kairos explain`, `kairos doctor`
- [ ] Phase 5 — GitHub/GitLab integration, Jira, Linear, MCP support

## License

Apache 2.0
