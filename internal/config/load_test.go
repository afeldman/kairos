package config

import (
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func writeConfigFile(t *testing.T, home, contents string) {
	t.Helper()
	dir := filepath.Join(home, ".config", "kairos")
	if err := osMkdirAll(dir); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	path := filepath.Join(dir, "config.yaml")
	if err := osWriteFile(path, contents); err != nil {
		t.Fatalf("write config file: %v", err)
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Provider != "ollama" || cfg.Model != "qwen3:30b" || cfg.Style != "conventional" {
		t.Fatalf("Load() = %+v, want defaults", cfg)
	}
	if cfg.Temperature != 0.2 || cfg.History != 20 {
		t.Fatalf("Load() numeric defaults = %+v", cfg)
	}
}

func TestLoad_ConfigFileOverridesDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeConfigFile(t, home, "provider: openai\nmodel: gpt-4o\n")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Provider != "openai" || cfg.Model != "gpt-4o" {
		t.Fatalf("Load() = %+v, want provider=openai model=gpt-4o", cfg)
	}
	if cfg.Style != "conventional" {
		t.Fatalf("Load() Style = %q, want default to survive", cfg.Style)
	}
}

func TestLoad_EnvOverridesConfigFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeConfigFile(t, home, "provider: openai\n")
	t.Setenv("KAIROS_PROVIDER", "anthropic")

	cfg, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Provider != "anthropic" {
		t.Fatalf("Load().Provider = %q, want env override anthropic", cfg.Provider)
	}
}

func TestLoad_FlagsOverrideEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("KAIROS_PROVIDER", "anthropic")

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("provider", "", "")
	if err := flags.Set("provider", "gemini"); err != nil {
		t.Fatalf("flags.Set: %v", err)
	}

	cfg, err := Load(flags)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Provider != "gemini" {
		t.Fatalf("Load().Provider = %q, want flag override gemini", cfg.Provider)
	}
}
