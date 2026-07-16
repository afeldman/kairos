package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newTagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tag",
		Short: "Generate an annotated tag message",
		Long:  `Generates a git tag message based on the latest commits.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos tag: not yet implemented")
		},
	}
}

func newReleaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "release",
		Short: "Generate release notes",
		Long:  `Generates release notes from commits since the last tag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos release: not yet implemented")
		},
	}
}

func newChangelogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "changelog",
		Short: "Update or generate CHANGELOG.md",
		Long:  `Updates CHANGELOG.md with new entries while preserving manual edits.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos changelog: not yet implemented")
		},
	}
}

func newExplainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "explain",
		Short: "Explain staged or committed changes",
		Long:  `Provides a plain-language explanation of what the changes do.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos explain: not yet implemented")
		},
	}
}

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose Kairos configuration and environment",
		Long: `Checks that Kairos can find Git, reach the configured provider,
and that the configuration is valid.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos doctor: not yet implemented")
		},
	}
}

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Display or edit Kairos configuration",
		Long:  `Shows the current configuration or opens the config file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos config: not yet implemented")
		},
	}
}

func newProvidersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "providers",
		Short: "List available LLM providers",
		Long:  `Lists all registered LLM providers and their status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("kairos providers: not yet implemented")
		},
	}
}
