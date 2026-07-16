// Package provider defines the Provider interface that all LLM providers must
// satisfy, along with a simple registry for looking up providers by name.
package provider

import (
	"context"
	"fmt"
)

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request holds the parameters for an LLM generation call.
type Request struct {
	Model       string
	Temperature float64
	Messages    []Message
}

// Provider is the interface that all LLM backends (Ollama, OpenAI, Anthropic,
// etc.) must implement.
type Provider interface {
	// Name returns a short identifier like "ollama" or "openai".
	Name() string
	// Generate sends a request and returns the raw text response.
	Generate(ctx context.Context, req Request) (string, error)
}

// Registry maps provider names to constructors. Providers register themselves
// in their init() function or are registered explicitly during CLI wiring.
var Registry = map[string]func(cfg interface{}) (Provider, error){}

// Register adds a provider constructor to the global registry.
func Register(name string, fn func(cfg interface{}) (Provider, error)) {
	Registry[name] = fn
}

// Get returns a provider by name, or an error if it is not registered.
func Get(name string, cfg interface{}) (Provider, error) {
	fn, ok := Registry[name]
	if !ok {
		known := make([]string, 0, len(Registry))
		for k := range Registry {
			known = append(known, k)
		}
		return nil, fmt.Errorf("unknown provider %q (known: %v)", name, known)
	}
	return fn(cfg)
}
