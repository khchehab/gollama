package gollama_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/khchehab/gollama"
)

// TestGenerateChatMessage_Success verifies a happy-path non-streaming chat response.
func TestGenerateChatMessage_Success(t *testing.T) {
	body := `{
		"model":"llama3",
		"created_at":"2024-01-01T00:00:00Z",
		"message":{"role":"assistant","content":"Hi there!"},
		"done":true,
		"done_reason":"stop"
	}`

	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/chat" {
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

	resp, err := c.GenerateChatMessage(context.Background(), gollama.GenerateChatMessageRequest{
		Model: "llama3",
		Messages: []gollama.GenerateChatMessageRequestMessage{
			{Role: gollama.RoleUser, Content: "Hello"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Model != "llama3" {
		t.Errorf("expected model %q, got %q", "llama3", resp.Model)
	}
	if resp.Message.Content != "Hi there!" {
		t.Errorf("expected content %q, got %q", "Hi there!", resp.Message.Content)
	}
	if !resp.Done {
		t.Error("expected done=true")
	}
}

// TestGenerateChatMessage_APIError verifies that a non-200 response propagates as *ErrorResponse.
func TestGenerateChatMessage_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`), nil
	}))

	_, err := c.GenerateChatMessage(context.Background(), gollama.GenerateChatMessageRequest{Model: "llama3"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
	if errResp.Error() != "invalid request" {
		t.Errorf("unexpected message %q", errResp.Error())
	}
}

// TestGenerateChatMessage_TransportError verifies transport-level errors are propagated.
func TestGenerateChatMessage_TransportError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}))

	_, err := c.GenerateChatMessage(context.Background(), gollama.GenerateChatMessageRequest{Model: "llama3"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestGenerateChatMessageStream_Success verifies the streaming chat yields all chunks.
func TestGenerateChatMessageStream_Success(t *testing.T) {
	ndjson := `{"model":"llama3","message":{"role":"assistant","content":"Hi"},"done":false}` + "\n" +
		`{"model":"llama3","message":{"role":"assistant","content":"!"},"done":false}` + "\n" +
		`{"model":"llama3","message":{"role":"assistant","content":""},"done":true}` + "\n"

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

	seq := c.GenerateChatMessageStream(context.Background(), gollama.GenerateChatMessageRequest{
		Model:    "llama3",
		Messages: []gollama.GenerateChatMessageRequestMessage{{Role: gollama.RoleUser, Content: "Hello"}},
	})

	var chunks []*gollama.GenerateChatMessageChunk
	for chunk, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		chunks = append(chunks, chunk)
	}

	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if chunks[0].Message.Content != "Hi" {
		t.Errorf("chunk[0] content: expected %q, got %q", "Hi", chunks[0].Message.Content)
	}
	if !chunks[2].Done {
		t.Error("last chunk: expected done=true")
	}
}

// TestGenerateChatMessageStream_APIError verifies a non-200 status yields an error from the iterator.
func TestGenerateChatMessageStream_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusInternalServerError, `{"error":"server error"}`), nil
	}))

	seq := c.GenerateChatMessageStream(context.Background(), gollama.GenerateChatMessageRequest{Model: "llama3"})
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

// TestGenerateChatMessageStream_InlineError verifies inline stream errors are propagated.
func TestGenerateChatMessageStream_InlineError(t *testing.T) {
	ndjson := `{"model":"llama3","message":{"role":"assistant","content":"partial"},"done":false}` + "\n" +
		`{"error":"context exceeded"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.GenerateChatMessageStream(context.Background(), gollama.GenerateChatMessageRequest{Model: "llama3"})
	var chunks []*gollama.GenerateChatMessageChunk
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

// TestGenerateChatMessageStream_EarlyCancellation verifies returning false stops iteration cleanly.
func TestGenerateChatMessageStream_EarlyCancellation(t *testing.T) {
	ndjson := `{"model":"llama3","message":{"role":"assistant","content":"a"},"done":false}` + "\n" +
		`{"model":"llama3","message":{"role":"assistant","content":"b"},"done":true}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.GenerateChatMessageStream(context.Background(), gollama.GenerateChatMessageRequest{Model: "llama3"})
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
