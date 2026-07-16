// Package cmd contains all Cobra command definitions.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/afeldman/kairos/internal/build"
	"github.com/afeldman/kairos/internal/config"
	projctx "github.com/afeldman/kairos/internal/context"
	"github.com/afeldman/kairos/internal/detect"
	"github.com/afeldman/kairos/internal/formatter"
	"github.com/afeldman/kairos/internal/git"
	"github.com/afeldman/kairos/internal/i18n"
	"github.com/afeldman/kairos/internal/prompt"
	"github.com/afeldman/kairos/internal/provider"
	_ "github.com/afeldman/kairos/internal/provider/ollama"
	_ "github.com/afeldman/kairos/internal/provider/openai"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root kairos command.
func NewRootCmd() *cobra.Command {
	var (
		cfgFile         string
		providerName    string
		model           string
		temperature     float64
		history         int
		style           string
		language        string
		updateChangelog bool
		baseURL         string
		apiKey          string
	)

	root := &cobra.Command{
		Use:     "kairos",
		Short:   "Kairos — Git Context Engine",
		Version: fmt.Sprintf("%s (commit=%s date=%s)", build.Version, build.Commit, build.Date),
		Long: `Kairos understands the history of a Git repository to generate
high-quality commit messages, tags, releases, and changelog entries.

By default (no subcommand), Kairos generates a commit message from staged changes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommit(cmd, config.Config{
				Provider:        providerName,
				Model:           model,
				Temperature:     temperature,
				History:         history,
				Style:           style,
				Language:        language,
				UpdateChangelog: updateChangelog,
				BaseURL:         baseURL,
				APIKey:          apiKey,
			})
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Flags
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: platform config dir, e.g. ~/Library/Application Support/kairos/config.yaml on macOS)")
	root.PersistentFlags().StringVarP(&providerName, "provider", "p", "", "LLM provider (ollama, openai, lmstudio, gomodel, anthropic, gemini, openrouter)")
	root.PersistentFlags().StringVarP(&model, "model", "m", "", "LLM model name")
	root.PersistentFlags().Float64Var(&temperature, "temperature", 0, "LLM temperature (0.0-1.0)")
	root.PersistentFlags().IntVar(&history, "history", 0, "number of recent commits to include")
	root.PersistentFlags().StringVarP(&style, "style", "s", "", "output style (conventional, github)")
	root.PersistentFlags().StringVarP(&language, "language", "l", "", "output language")
	root.PersistentFlags().BoolVar(&updateChangelog, "update-changelog", false, "update CHANGELOG.md")
	root.PersistentFlags().StringVar(&baseURL, "base-url", "", "override the provider's API endpoint (e.g. http://localhost:1234/v1 for LM Studio)")
	root.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for the selected provider")

	// Subcommands
	root.AddCommand(newCommitCmd())
	root.AddCommand(newTagCmd())
	root.AddCommand(newReleaseCmd())
	root.AddCommand(newChangelogCmd())
	root.AddCommand(newExplainCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newConfigCmd())
	root.AddCommand(newProvidersCmd())

	return root
}

// Execute is the entry point for the CLI.
func Execute() {
	root := NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "kairos: error: %s\n", err)
		os.Exit(1)
	}
}

// runCommit is the core logic shared by the root command and `kairos commit`.
func runCommit(cmd *cobra.Command, cfg config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// 1. Load config (merged with CLI flags).
	loadedCfg, err := config.Load(cmd.Flags())
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	// Merge — CLI flags passed directly take precedence.
	if cfg.Provider != "" {
		loadedCfg.Provider = cfg.Provider
	}
	if cfg.Model != "" {
		loadedCfg.Model = cfg.Model
	}
	if cfg.Temperature != 0 {
		loadedCfg.Temperature = cfg.Temperature
	}
	if cfg.History != 0 {
		loadedCfg.History = cfg.History
	}
	if cfg.Style != "" {
		loadedCfg.Style = cfg.Style
	}
	if cfg.Language != "" {
		loadedCfg.Language = cfg.Language
	}
	if cfg.BaseURL != "" {
		loadedCfg.BaseURL = cfg.BaseURL
	}
	if cfg.APIKey != "" {
		loadedCfg.APIKey = cfg.APIKey
	}

	// 2. Detect project type (before git context, independent).
	_ = detect.Detect(".")

	locale := i18n.Detect(loadedCfg.Language)

	// 3. Open repository.
	if !git.IsRepo(ctx, ".") {
		return fmt.Errorf("%s", i18n.T(locale, i18n.NotAGitRepo))
	}
	repo := git.NewRepository(".")

	// 4. Ensure there are staged changes.
	if err := git.EnsureStaged(ctx, repo); err != nil {
		if errors.Is(err, git.ErrNothingStaged) {
			return fmt.Errorf("%s", i18n.T(locale, i18n.NothingStaged))
		}
		return err
	}

	// 5. Build project context.
	builder := projctx.NewBuilder()
	pc, err := builder.Build(ctx, repo, loadedCfg)
	if err != nil {
		return fmt.Errorf("build context: %w", err)
	}

	// 6. Build prompt.
	messages := prompt.Build(pc, loadedCfg)

	// 7. Get provider.
	prov, err := provider.Get(loadedCfg.Provider, loadedCfg)
	if err != nil {
		return err
	}

	// 8. Generate.
	req := provider.Request{
		Model:       loadedCfg.Model,
		Temperature: loadedCfg.Temperature,
		Messages:    messages,
	}
	raw, err := prov.Generate(ctx, req)
	if err != nil {
		return fmt.Errorf("%s: %w", prov.Name(), err)
	}

	// 9. Parse and render.
	msg := formatter.Parse(raw)
	result, err := formatter.Render(loadedCfg.Style, msg)
	if err != nil {
		return err
	}

	// 10. Print result (no trailing newline — works with git commit -m).
	fmt.Print(result)
	return nil
}
