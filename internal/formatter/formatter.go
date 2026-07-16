// Package formatter parses LLM responses into structured commit messages and
// renders them according to a style template (e.g. Conventional Commits).
package formatter

import (
	"encoding/json"
	"strings"
	"text/template"
)

// CommitMessage holds the parsed fields of a commit message.
type CommitMessage struct {
	Type     string `json:"type"`
	Scope    string `json:"scope"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Breaking string `json:"breaking"`
}

// Parse attempts to unmarshal raw JSON into a CommitMessage. If that fails,
// it falls back to treating the first line as the subject and the rest as
// the body.
func Parse(raw string) CommitMessage {
	raw = strings.TrimSpace(raw)

	var msg CommitMessage
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		// Basic validation: at minimum we expect a type and subject.
		if msg.Type != "" && msg.Subject != "" {
			return msg
		}
	}

	// Fallback: split on newlines.
	lines := strings.SplitN(raw, "\n", 2)
	msg = CommitMessage{
		Subject: strings.TrimSpace(lines[0]),
	}
	if len(lines) > 1 {
		msg.Body = strings.TrimSpace(lines[1])
	}
	return msg
}

// Render produces the final commit message string using the named style.
// Currently only "conventional" is supported.
func Render(style string, msg CommitMessage) (string, error) {
	tmpl, ok := templates[style]
	if !ok {
		return "", errUnknownStyle(style)
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, msg); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// templates holds the named templates. Additive — adding a new style is a
// single map entry.
var templates = map[string]*template.Template{
	"conventional": template.Must(template.New("conventional").Parse(conventionalTemplate)),
	"github":       template.Must(template.New("github").Parse(githubTemplate)),
}

const conventionalTemplate = `{{.Type -}}{{if .Scope}}({{.Scope}}){{end -}}: {{.Subject -}}
{{if .Body}}

{{.Body -}}
{{end -}}
{{if .Breaking}}

BREAKING CHANGE: {{.Breaking -}}
{{end}}`

const githubTemplate = `{{.Type -}}{{if .Scope}}({{.Scope}}){{end -}}: {{.Subject -}}
{{if .Body}}

{{.Body -}}
{{end -}}
{{if .Breaking}}

BREAKING CHANGE: {{.Breaking -}}
{{end}}

Co-authored-by: kairos`
