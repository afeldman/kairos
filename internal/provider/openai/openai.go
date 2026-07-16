// Package openai provides a Provider that talks to any OpenAI-compatible
// chat completions API. The same client backs three registered providers:
//
//   - "openai": the real OpenAI API
//   - "lmstudio": a local LM Studio server (no native Ollama API, so it
//     must be reached through its OpenAI-compatible endpoint)
//   - "gomodel": a GoModel gateway (https://github.com/ENTERPILOT/GoModel)
//     fronting one or more upstream providers
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afeldman/kairos/internal/config"
	"github.com/afeldman/kairos/internal/provider"
)

// Client implements provider.Provider for any server exposing an
// OpenAI-compatible POST /chat/completions endpoint.
type Client struct {
	name    string
	baseURL string
	apiKey  string
	http    *http.Client
}

// Config holds settings for an OpenAI-compatible client.
type Config struct {
	// Name is the identifier reported by Name(), e.g. "lmstudio".
	Name string
	// BaseURL is the API root, e.g. "http://localhost:1234/v1".
	BaseURL string
	// APIKey is sent as a Bearer token. Servers that don't check it
	// (LM Studio) accept any non-empty value.
	APIKey string
}

// NewClient creates a new OpenAI-compatible provider client.
func NewClient(cfg Config) *Client {
	name := cfg.Name
	if name == "" {
		name = "openai"
	}
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = "none"
	}
	return &Client{
		name:    name,
		baseURL: cfg.BaseURL,
		apiKey:  apiKey,
		http: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) Name() string { return c.name }

// chatRequest is the OpenAI /chat/completions request body.
type chatRequest struct {
	Model       string             `json:"model"`
	Messages    []provider.Message `json:"messages"`
	Temperature float64            `json:"temperature,omitempty"`
	Stream      bool               `json:"stream"`
}

// chatResponse is the OpenAI /chat/completions response body (non-streaming).
type chatResponse struct {
	Choices []struct {
		Message provider.Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) Generate(ctx context.Context, req provider.Request) (string, error) {
	body := chatRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		Stream:      false,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("%s: marshal request: %w", c.name, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("%s: create request: %w", c.name, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("%s: %w (is the server running at %s?)", c.name, err, c.baseURL)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: read response: %w", c.name, err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return "", fmt.Errorf("%s: %s: parse response: %w", c.name, resp.Status, err)
	}
	if chatResp.Error != nil {
		return "", fmt.Errorf("%s: %s", c.name, chatResp.Error.Message)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s: %s", c.name, resp.Status, string(raw))
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("%s: empty response", c.name)
	}

	return chatResp.Choices[0].Message.Content, nil
}

// defaultBaseURLs maps each registered provider name to its default endpoint.
var defaultBaseURLs = map[string]string{
	"openai":   "https://api.openai.com/v1",
	"lmstudio": "http://localhost:1234/v1",
	"gomodel":  "http://localhost:8080/v1",
}

func init() {
	for name, defaultURL := range defaultBaseURLs {
		provider.Register(name, func(cfg any) (provider.Provider, error) {
			c, _ := cfg.(config.Config)
			baseURL := c.BaseURL
			if baseURL == "" {
				baseURL = defaultURL
			}
			return NewClient(Config{Name: name, BaseURL: baseURL, APIKey: c.APIKey}), nil
		})
	}
}
