package cmd

import (
	"github.com/afeldman/kairos/internal/config"
	"github.com/spf13/cobra"
)

func newCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Generate a commit message from staged changes",
		Long:  `Analyzes the staged changes and generates a Conventional Commits message.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommit(cmd, configFromFlags(cmd))
		},
	}
}

// configFromFlags reads the current Config from a command's persistent flags.
// The root command defines all flags; subcommands inherit them.
func configFromFlags(cmd *cobra.Command) config.Config {
	cfg := config.Config{}
	if v, err := cmd.Flags().GetString("provider"); err == nil {
		cfg.Provider = v
	}
	if v, err := cmd.Flags().GetString("model"); err == nil {
		cfg.Model = v
	}
	if v, err := cmd.Flags().GetFloat64("temperature"); err == nil {
		cfg.Temperature = v
	}
	if v, err := cmd.Flags().GetInt("history"); err == nil {
		cfg.History = v
	}
	if v, err := cmd.Flags().GetString("style"); err == nil {
		cfg.Style = v
	}
	if v, err := cmd.Flags().GetString("language"); err == nil {
		cfg.Language = v
	}
	if v, err := cmd.Flags().GetBool("update-changelog"); err == nil {
		cfg.UpdateChangelog = v
	}
	if v, err := cmd.Flags().GetString("base-url"); err == nil {
		cfg.BaseURL = v
	}
	if v, err := cmd.Flags().GetString("api-key"); err == nil {
		cfg.APIKey = v
	}
	return cfg
}
