package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Load resolves Config from, in increasing priority: built-in defaults,
// the user's ~/.config/kairos/config.yaml, KAIROS_* environment
// variables, and finally flags (if a non-nil FlagSet is given).
func Load(flags *pflag.FlagSet) (Config, error) {
	v := viper.New()
	d := defaults()
	v.SetDefault("provider", d.Provider)
	v.SetDefault("model", d.Model)
	v.SetDefault("temperature", d.Temperature)
	v.SetDefault("history", d.History)
	v.SetDefault("style", d.Style)
	v.SetDefault("language", d.Language)
	v.SetDefault("update_changelog", d.UpdateChangelog)
	v.SetDefault("base_url", d.BaseURL)
	v.SetDefault("api_key", d.APIKey)

	v.SetEnvPrefix("KAIROS")
	v.AutomaticEnv()

	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "kairos"))
	}
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("config: read config file: %w", err)
		}
	}

	if flags != nil {
		if err := v.BindPFlags(flags); err != nil {
			return Config{}, fmt.Errorf("config: bind flags: %w", err)
		}
	}

	return Config{
		Provider:        v.GetString("provider"),
		Model:           v.GetString("model"),
		Temperature:     v.GetFloat64("temperature"),
		History:         v.GetInt("history"),
		Style:           v.GetString("style"),
		Language:        v.GetString("language"),
		UpdateChangelog: v.GetBool("update_changelog"),
		BaseURL:         v.GetString("base_url"),
		APIKey:          v.GetString("api_key"),
	}, nil
}
