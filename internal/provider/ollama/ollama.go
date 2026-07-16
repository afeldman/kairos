// Package ollama provides a Provider that talks to a local Ollama instance
// via its REST API. It requires no additional dependencies beyond the standard
// library's net/http.
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afeldman/kairos/internal/provider"
)

const defaultBaseURL = "http://localhost:11434"

// Client implements provider.Provider for Ollama.
type Client struct {
	baseURL string
	http    *http.Client
}

// Config holds optional settings for the Ollama client.
type Config struct {
	BaseURL string
}

// NewClient creates a new Ollama provider client.
func NewClient(cfg Config) *Client {
	base := cfg.BaseURL
	if base == "" {
		base = defaultBaseURL
	}
	return &Client{
		baseURL: base,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) Name() string { return "ollama" }

// chatRequest is the Ollama /api/chat request body.
type chatRequest struct {
	Model    string             `json:"model"`
	Messages []provider.Message `json:"messages"`
	Options  map[string]float64 `json:"options,omitempty"`
	Stream   bool               `json:"stream"`
}

// chatResponse is the Ollama /api/chat response body (non-streaming).
type chatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func (c *Client) Generate(ctx context.Context, req provider.Request) (string, error) {
	body := chatRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   false,
	}
	if req.Temperature > 0 {
		body.Options = map[string]float64{"temperature": req.Temperature}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/chat", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("ollama: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ollama: %w (is ollama running? try 'ollama serve')", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ollama: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama: %s: %s", resp.Status, string(raw))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return "", fmt.Errorf("ollama: parse response: %w", err)
	}

	return chatResp.Message.Content, nil
}

func init() {
	provider.Register("ollama", func(cfg interface{}) (provider.Provider, error) {
		// No special config for now; could read from cfg if typed.
		return NewClient(Config{}), nil
	})
}
