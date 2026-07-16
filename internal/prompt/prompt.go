// Package prompt builds structured prompts for LLM interactions.
// Prompts are deterministic, structured, and include full project context.
package prompt

import (
	"fmt"
	"strings"

	"github.com/afeldman/kairos/internal/config"
	"github.com/afeldman/kairos/internal/context"
	"github.com/afeldman/kairos/internal/provider"
)

// Build creates a set of chat messages (system + user) for the "commit" task.
func Build(pc context.ProjectContext, cfg config.Config) []provider.Message {
	return []provider.Message{
		{Role: "system", Content: systemPrompt(cfg)},
		{Role: "user", Content: userPrompt(pc, cfg)},
	}
}

// outputLanguage returns cfg's configured LLM output language, defaulting
// to English when unset.
func outputLanguage(lang string) string {
	if lang == "" {
		return "english"
	}
	return lang
}

func systemPrompt(cfg config.Config) string {
	var b strings.Builder
	b.WriteString("You are a Git Context Engine that generates high-quality commit messages.\n\n")
	b.WriteString("Rules:\n")
	b.WriteString("- Respond with ONLY a JSON object, no other text.\n")
	b.WriteString("- Use Conventional Commits format.\n")
	fmt.Fprintf(&b, "- Language: %s\n", outputLanguage(cfg.Language))
	b.WriteString("- Keep the subject line under 72 characters.\n")
	b.WriteString("- The body should explain what and why, not how.\n")
	b.WriteString("\n")
	b.WriteString("JSON format:\n")
	b.WriteString(`{"type":"<type>","scope":"<scope or empty>","subject":"<short description>","body":"<detailed body or empty>","breaking":"<breaking change description or empty>"}`)
	b.WriteString("\n\n")
	b.WriteString("Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert\n")
	b.WriteString("If there is a breaking change, set type to feat! or fix! and describe it in the breaking field.\n")
	return b.String()
}

func userPrompt(pc context.ProjectContext, cfg config.Config) string {
	var b strings.Builder
	b.WriteString("<task>Generate a commit message for the staged changes below.</task>\n\n")

	if pc.Branch != "" {
		fmt.Fprintf(&b, "<branch>%s</branch>\n\n", pc.Branch)
	}

	if pc.ProjectType != "" && pc.ProjectType != "unknown" {
		fmt.Fprintf(&b, "<project-type>%s</project-type>\n\n", pc.ProjectType)
	}

	if pc.LastTag != "" {
		fmt.Fprintf(&b, "<last-tag>%s</last-tag>\n", pc.LastTag)
	}
	if pc.Describe != "" {
		fmt.Fprintf(&b, "<describe>%s</describe>\n", pc.Describe)
	}

	if len(pc.RecentCommits) > 0 {
		b.WriteString("\n<recent-commits>\n")
		for _, c := range pc.RecentCommits {
			fmt.Fprintf(&b, "  %s\n", c)
		}
		b.WriteString("</recent-commits>\n\n")
	}

	if len(pc.ChangedFiles) > 0 {
		b.WriteString("<changed-files>\n")
		for _, f := range pc.ChangedFiles {
			fmt.Fprintf(&b, "  %s\n", f)
		}
		b.WriteString("</changed-files>\n\n")
	}

	if pc.ReadmeExcerpt != "" {
		b.WriteString("<project-readme-excerpt>\n")
		b.WriteString(pc.ReadmeExcerpt)
		b.WriteString("\n</project-readme-excerpt>\n\n")
	}

	if pc.ChangelogExcerpt != "" {
		b.WriteString("<changelog-excerpt>\n")
		b.WriteString(pc.ChangelogExcerpt)
		b.WriteString("\n</changelog-excerpt>\n\n")
	}

	b.WriteString("<staged-diff>\n")
	b.WriteString(pc.Diff)
	b.WriteString("\n</staged-diff>\n")

	return b.String()
}
