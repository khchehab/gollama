package gollama_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/khchehab/gollama"
)

// blockingTransport returns a RoundTripper that blocks until the request context is cancelled,
// then returns the context's error.
func blockingTransport() http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
}

// TestGenerateResponse_ContextCancelled verifies that a cancelled context aborts the request.
func TestGenerateResponse_ContextCancelled(t *testing.T) {
	c := newClientWithTransport(t, blockingTransport())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GenerateResponse(ctx, gollama.GenerateResponseRequest{Model: "llama3", Prompt: "hello"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// TestGenerateResponseStream_ContextCancelled verifies that a cancelled context aborts a stream before it yields.
func TestGenerateResponseStream_ContextCancelled(t *testing.T) {
	c := newClientWithTransport(t, blockingTransport())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	seq := c.GenerateResponseStream(ctx, gollama.GenerateResponseRequest{Model: "llama3", Prompt: "hello"})
	var gotErr error
	for _, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
	}
	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", gotErr)
	}
}

// TestGenerateChatMessage_ContextCancelled verifies that a cancelled context aborts the request.
func TestGenerateChatMessage_ContextCancelled(t *testing.T) {
	c := newClientWithTransport(t, blockingTransport())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GenerateChatMessage(ctx, gollama.GenerateChatMessageRequest{
		Model:    "llama3",
		Messages: []gollama.GenerateChatMessageRequestMessage{{Role: gollama.RoleUser, Content: "hello"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// TestGenerateChatMessageStream_ContextCancelled verifies that a cancelled context aborts a stream before it yields.
func TestGenerateChatMessageStream_ContextCancelled(t *testing.T) {
	c := newClientWithTransport(t, blockingTransport())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	seq := c.GenerateChatMessageStream(ctx, gollama.GenerateChatMessageRequest{
		Model:    "llama3",
		Messages: []gollama.GenerateChatMessageRequestMessage{{Role: gollama.RoleUser, Content: "hello"}},
	})
	var gotErr error
	for _, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
	}
	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", gotErr)
	}
}

// TestGenerateEmbeddings_ContextCancelled verifies that a cancelled context aborts the request.
func TestGenerateEmbeddings_ContextCancelled(t *testing.T) {
	c := newClientWithTransport(t, blockingTransport())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GenerateEmbeddings(ctx, gollama.GenerateEmbeddingRequest{Model: "nomic-embed", Input: []string{"hello"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

// TestPullModelStream_ContextCancelled verifies that a cancelled context aborts a stream before it yields.
func TestPullModelStream_ContextCancelled(t *testing.T) {
	c := newClientWithTransport(t, blockingTransport())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	seq := c.PullModelStream(ctx, gollama.PullModelRequest{Model: "llama3"})
	var gotErr error
	for _, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
	}
	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", gotErr)
	}
}
