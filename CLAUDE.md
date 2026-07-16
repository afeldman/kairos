# CLAUDE.md

# Kairos Development Guide

## Project Vision

Kairos is a Git Context Engine.

It understands the history of a Git repository to generate meaningful commit messages, release notes, changelog updates, and other Git-related artifacts.

Never think of Kairos as "just another AI commit tool."

The context engine is the core product.

---

## Philosophy

Context beats prompts.

A commit message should never be generated from only a git diff.

Always prefer understanding over summarization.

---

## Design Principles

- Simple architecture
- Idiomatic Go
- Modular packages
- Dependency injection
- Interface-driven design
- Testability first
- Small packages
- Clean APIs
- No unnecessary abstractions

---

## CLI

Frameworks

- Cobra
- Viper

The CLI must remain intuitive.

The default command should always generate a commit message.

Examples

git commit -am "$(kairos)"

git tag -a v1.5.0 -m "$(kairos tag)"

---

## Provider Design

Providers must implement a common interface.

Business logic must never depend on a specific provider.

Default providers

- Ollama
- LM Studio
- OpenAI
- Anthropic
- Gemini
- OpenRouter

---

## Git Context

Prefer collecting

git status

git diff

git diff --cached

git log

git describe

git branch

changed files

last tag

Only include relevant history.

Avoid overwhelming the LLM with unnecessary data.

---

## Prompt Design

Prompts should be deterministic.

Keep temperature low.

Use structured prompts.

Avoid conversational prompts.

Prefer XML or Markdown sections.

Always specify the expected output format.

---

## Code Style

Use

context.Context

errors.Join()

slog

table-driven tests

Go 1.25+

Never panic inside business logic.

Return errors.

---

## Testing

Every package should have tests.

Mock external dependencies.

Avoid network access.

Tests should be deterministic.

---

## Documentation

Every exported type should be documented.

Keep README examples updated.

Document architectural decisions.

---

## Performance

Git operations should be cached where appropriate.

Avoid unnecessary provider calls.

Large repositories should remain responsive.

---

## Security

Never execute generated shell commands.

Never automatically commit.

Never automatically push.

Never expose secrets.

Respect .gitignore.

---

## Future Features

Design with extensibility for

- Release automation
- Semantic Versioning
- AI repository memory
- MCP servers
- GitHub Actions
- GitLab CI
- Jira
- Linear

Avoid implementing speculative features until requested.

---

## Guiding Principle

Every change should make Kairos a better Git Context Engine, not merely a better prompt wrapper.
