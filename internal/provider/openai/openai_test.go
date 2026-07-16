package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/afeldman/kairos/internal/provider"
)

func TestClient_Generate(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantErr    bool
		wantResult string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
					t.Errorf("Authorization header = %q, want %q", got, "Bearer test-key")
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"choices": []map[string]any{
						{"message": map[string]string{"role": "assistant", "content": "feat: add thing"}},
					},
				})
			},
			wantResult: "feat: add thing",
		},
		{
			name: "api error payload",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]string{"message": "invalid api key"},
				})
			},
			wantErr: true,
		},
		{
			name: "empty choices",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{}})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			c := NewClient(Config{Name: "lmstudio", BaseURL: srv.URL, APIKey: "test-key"})
			got, err := c.Generate(context.Background(), provider.Request{
				Model:    "test-model",
				Messages: []provider.Message{{Role: "user", Content: "diff"}},
			})

			if tt.wantErr {
				if err == nil {
					t.Fatal("Generate() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Generate() unexpected error: %v", err)
			}
			if got != tt.wantResult {
				t.Errorf("Generate() = %q, want %q", got, tt.wantResult)
			}
		})
	}
}

func TestClient_Name(t *testing.T) {
	if got := NewClient(Config{}).Name(); got != "openai" {
		t.Errorf("Name() = %q, want %q (default)", got, "openai")
	}
	if got := NewClient(Config{Name: "gomodel"}).Name(); got != "gomodel" {
		t.Errorf("Name() = %q, want %q", got, "gomodel")
	}
}

func TestRegistry_DefaultBaseURLs(t *testing.T) {
	for name, want := range defaultBaseURLs {
		p, err := provider.Get(name, nil)
		if err != nil {
			t.Fatalf("provider.Get(%q) error: %v", name, err)
		}
		c, ok := p.(*Client)
		if !ok {
			t.Fatalf("provider.Get(%q) = %T, want *Client", name, p)
		}
		if c.baseURL != want {
			t.Errorf("%s baseURL = %q, want default %q", name, c.baseURL, want)
		}
	}
}
