package gollama_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/khchehab/gollama"
)

// TestGenerateResponse_Success verifies a happy-path non-streaming generate response.
func TestGenerateResponse_Success(t *testing.T) {
	body := `{"model":"llama3","created_at":"2024-01-01T00:00:00Z","response":"Hello!","done":true,"done_reason":"stop"}`
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/generate" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("could not decode request: %v", err)
		}
		if req["stream"] != false {
			t.Errorf("expected stream=false, got %v", req["stream"])
		}
		return respond(http.StatusOK, body), nil
	}))

	resp, err := c.GenerateResponse(context.Background(), gollama.GenerateResponseRequest{
		Model:  "llama3",
		Prompt: "Hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Model != "llama3" {
		t.Errorf("expected model %q, got %q", "llama3", resp.Model)
	}
	if resp.Response != "Hello!" {
		t.Errorf("expected response %q, got %q", "Hello!", resp.Response)
	}
	if !resp.Done {
		t.Error("expected done=true")
	}
}

// TestGenerateResponse_APIError verifies that a non-200 API response propagates as *ErrorResponse.
func TestGenerateResponse_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusNotFound, `{"error":"model not found"}`), nil
	}))

	_, err := c.GenerateResponse(context.Background(), gollama.GenerateResponseRequest{Model: "missing"})
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

// TestGenerateResponse_TransportError verifies that a transport-level error is propagated.
func TestGenerateResponse_TransportError(t *testing.T) {
	transportErr := errors.New("dial tcp: connection refused")
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, transportErr
	}))

	_, err := c.GenerateResponse(context.Background(), gollama.GenerateResponseRequest{Model: "llama3"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestGenerateResponse_AuthorizationHeader verifies the Authorization header is set when an API key is provided.
func TestGenerateResponse_AuthorizationHeader(t *testing.T) {
	var gotAuthHeader string
	c, err := gollama.NewClient(
		gollama.WithAPIKey("sk-secret"),
		gollama.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			gotAuthHeader = r.Header.Get("Authorization")
			return respond(http.StatusOK, `{"model":"llama3","response":"hi","done":true}`), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, _ = c.GenerateResponse(context.Background(), gollama.GenerateResponseRequest{Model: "llama3"})
	if gotAuthHeader != "Bearer sk-secret" {
		t.Errorf("expected Authorization header %q, got %q", "Bearer sk-secret", gotAuthHeader)
	}
}

// TestGenerateResponseStream_Success verifies streaming yields all chunks in order.
func TestGenerateResponseStream_Success(t *testing.T) {
	ndjson := `{"model":"llama3","response":"He","done":false}` + "\n" +
		`{"model":"llama3","response":"llo","done":false}` + "\n" +
		`{"model":"llama3","response":"","done":true,"done_reason":"stop"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("could not decode request: %v", err)
		}
		if req["stream"] != true {
			t.Errorf("expected stream=true, got %v", req["stream"])
		}
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.GenerateResponseStream(context.Background(), gollama.GenerateResponseRequest{
		Model:  "llama3",
		Prompt: "Hello",
	})

	var chunks []*gollama.GenerateResponseChunk
	for chunk, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if chunks[0].Response != "He" {
		t.Errorf("chunk[0]: expected %q, got %q", "He", chunks[0].Response)
	}
	if chunks[1].Response != "llo" {
		t.Errorf("chunk[1]: expected %q, got %q", "llo", chunks[1].Response)
	}
	if !chunks[2].Done {
		t.Error("last chunk: expected done=true")
	}
}

// TestGenerateResponseStream_APIError verifies that a non-200 response propagates as an error from the iterator.
func TestGenerateResponseStream_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusInternalServerError, `{"error":"server error"}`), nil
	}))

	seq := c.GenerateResponseStream(context.Background(), gollama.GenerateResponseRequest{Model: "llama3"})
	var gotErr error
	for _, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
	}
	if gotErr == nil {
		t.Fatal("expected error from stream, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(gotErr, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", gotErr, gotErr)
	}
}

// TestGenerateResponseStream_InlineError verifies that an error embedded in the NDJSON stream is propagated.
func TestGenerateResponseStream_InlineError(t *testing.T) {
	ndjson := `{"model":"llama3","response":"partial","done":false}` + "\n" +
		`{"error":"context length exceeded"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.GenerateResponseStream(context.Background(), gollama.GenerateResponseRequest{Model: "llama3"})
	var chunks []*gollama.GenerateResponseChunk
	var gotErr error
	for chunk, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
		chunks = append(chunks, chunk)
	}
	if gotErr == nil {
		t.Fatal("expected inline error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(gotErr, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T", gotErr)
	}
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk before error, got %d", len(chunks))
	}
}

// TestGenerateResponseStream_EarlyCancellation verifies returning false stops iteration cleanly.
func TestGenerateResponseStream_EarlyCancellation(t *testing.T) {
	ndjson := `{"model":"llama3","response":"a","done":false}` + "\n" +
		`{"model":"llama3","response":"b","done":false}` + "\n" +
		`{"model":"llama3","response":"c","done":true}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.GenerateResponseStream(context.Background(), gollama.GenerateResponseRequest{Model: "llama3"})
	var count int
	for _, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		count++
		break
	}
	if count != 1 {
		t.Errorf("expected 1 chunk before cancellation, got %d", count)
	}
}
