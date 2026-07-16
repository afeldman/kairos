// Package i18n resolves the CLI's UI locale and translates the small set
// of user-facing error and easter-egg messages Kairos prints itself. It
// does not touch LLM prompt/output language (see internal/prompt), only
// what Kairos says on its own behalf.
package i18n

import (
	"os"
	"strings"
)

// Locale identifies a supported UI language.
type Locale string

const (
	En Locale = "en"
	De Locale = "de"
)

// Detect resolves the UI locale: an explicit cfgLanguage value (from
// config.yaml, KAIROS_LANGUAGE, or --language) takes priority, falling
// back to the system locale (LC_ALL, LC_MESSAGES, LANG), and finally
// English if nothing matches.
func Detect(cfgLanguage string) Locale {
	if l, ok := fromName(cfgLanguage); ok {
		return l
	}
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if l, ok := fromName(os.Getenv(env)); ok {
			return l
		}
	}
	return En
}

func fromName(s string) (Locale, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case s == "":
		return "", false
	case strings.HasPrefix(s, "de"), strings.Contains(s, "german"), strings.Contains(s, "deutsch"):
		return De, true
	case strings.HasPrefix(s, "en"), strings.Contains(s, "english"):
		return En, true
	}
	return "", false
}
