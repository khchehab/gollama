package gollama_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/khchehab/gollama"
)

// TestGenerateEmbeddings_Success verifies a happy-path embedding response.
func TestGenerateEmbeddings_Success(t *testing.T) {
	body := `{"model":"nomic-embed","embeddings":[[0.1,0.2,0.3],[0.4,0.5,0.6]],"total_duration":1000,"prompt_eval_count":2}`

	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/embed" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, body), nil
	}))

	resp, err := c.GenerateEmbeddings(context.Background(), gollama.GenerateEmbeddingRequest{
		Model: "nomic-embed",
		Input: []string{"hello", "world"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Model != "nomic-embed" {
		t.Errorf("expected model %q, got %q", "nomic-embed", resp.Model)
	}
	if len(resp.Embeddings) != 2 {
		t.Fatalf("expected 2 embeddings, got %d", len(resp.Embeddings))
	}
	if len(resp.Embeddings[0]) != 3 {
		t.Errorf("expected 3 dimensions in first embedding, got %d", len(resp.Embeddings[0]))
	}
	if resp.Embeddings[0][0] != 0.1 {
		t.Errorf("embeddings[0][0]: expected 0.1, got %f", resp.Embeddings[0][0])
	}
}

// TestGenerateEmbeddings_APIError verifies that a non-200 response propagates as *ErrorResponse.
func TestGenerateEmbeddings_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusUnprocessableEntity, `{"error":"model not found"}`), nil
	}))

	_, err := c.GenerateEmbeddings(context.Background(), gollama.GenerateEmbeddingRequest{
		Model: "missing",
		Input: []string{"text"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
	if errResp.Error() != "model not found" {
		t.Errorf("unexpected message %q", errResp.Error())
	}
}

// TestGenerateEmbeddings_TransportError verifies transport-level errors are propagated.
func TestGenerateEmbeddings_TransportError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}))

	_, err := c.GenerateEmbeddings(context.Background(), gollama.GenerateEmbeddingRequest{
		Model: "nomic-embed",
		Input: []string{"hello"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestGenerateEmbeddings_EmptyEmbeddings verifies that an empty embeddings array is handled without error.
func TestGenerateEmbeddings_EmptyEmbeddings(t *testing.T) {
	body := `{"model":"nomic-embed","embeddings":[]}`

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, body), nil
	}))

	resp, err := c.GenerateEmbeddings(context.Background(), gollama.GenerateEmbeddingRequest{
		Model: "nomic-embed",
		Input: []string{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Embeddings) != 0 {
		t.Errorf("expected 0 embeddings, got %d", len(resp.Embeddings))
	}
}
