// Package config loads Kairos configuration from flags, environment
// variables, and the user's config file, in that order of precedence.
package config

// Config holds all user-tunable Kairos settings.
type Config struct {
	Provider        string
	Model           string
	Temperature     float64
	History         int
	Style           string
	Language        string
	UpdateChangelog bool
}

func defaults() Config {
	return Config{
		Provider:        "ollama",
		Model:           "qwen3:30b",
		Temperature:     0.2,
		History:         20,
		Style:           "conventional",
		Language:        "english",
		UpdateChangelog: true,
	}
}
