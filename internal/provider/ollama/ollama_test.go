package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afeldman/kairos/internal/provider"
)

func TestGenerate_Success(t *testing.T) {
	var called bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != "POST" {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/api/chat" {
			t.Fatalf("path = %s, want /api/chat", r.URL.Path)
		}

		resp := map[string]interface{}{
			"message": map[string]interface{}{
				"content": `{"type":"feat","scope":"core","subject":"add ollama support","body":"","breaking":""}`,
			},
			"done": true,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClient(Config{BaseURL: srv.URL})
	req := provider.Request{
		Model:       "qwen3:30b",
		Temperature: 0.2,
		Messages: []provider.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Generate a commit message."},
		},
	}

	got, err := client.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !called {
		t.Fatal("server was not called")
	}
	if got == "" {
		t.Fatal("Generate() returned empty string")
	}
}

func TestGenerate_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"something broke"}`))
	}))
	defer srv.Close()

	client := NewClient(Config{BaseURL: srv.URL})
	_, err := client.Generate(context.Background(), provider.Request{Model: "test"})
	if err == nil {
		t.Fatal("expected error for server error")
	}
	t.Logf("got expected error: %v", err)
}

func TestGenerate_ConnectionRefused(t *testing.T) {
	// Use a port that's likely not in use.
	client := NewClient(Config{BaseURL: "http://localhost:42999"})
	_, err := client.Generate(context.Background(), provider.Request{Model: "test"})
	if err == nil {
		t.Fatal("expected error for connection refused")
	}
	t.Logf("got expected error: %v", err)
}

func TestClientName(t *testing.T) {
	c := NewClient(Config{})
	if got := c.Name(); got != "ollama" {
		t.Fatalf("Name() = %q, want ollama", got)
	}
}
