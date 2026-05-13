package gollama_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/khchehab/gollama"
)

// roundTripFunc is an http.RoundTripper backed by a function.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// respond builds a synthetic *http.Response with the given status and JSON body.
func respond(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// newClientWithTransport creates a gollama.Client wired to a custom RoundTripper.
func newClientWithTransport(t *testing.T, rt http.RoundTripper) *gollama.Client {
	t.Helper()
	c, err := gollama.NewClient(gollama.WithHTTPClient(&http.Client{Transport: rt}))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

// TestNewClient_Defaults verifies that a client can be created with no options.
func TestNewClient_Defaults(t *testing.T) {
	_, err := gollama.NewClient()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestNewClient_WithHostAndPort verifies WithHostAndPort produces a valid client.
func TestNewClient_WithHostAndPort(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithHostAndPort("localhost", 9999))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestNewClient_WithURL verifies that a custom URL is accepted.
func TestNewClient_WithURL(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithURL("http://myhost:1234"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestNewClient_InvalidURL verifies that a malformed URL returns InvalidBaseURLError.
func TestNewClient_InvalidURL(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithURL("not-a-url"))
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
	var urlErr *gollama.InvalidBaseURLError
	if !errors.As(err, &urlErr) {
		t.Fatalf("expected *InvalidBaseURLError, got %T: %v", err, err)
	}
}

// TestNewClient_CloudURLWithoutAPIKey verifies that using the Ollama cloud URL without an API key fails.
func TestNewClient_CloudURLWithoutAPIKey(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithURL("https://ollama.com"))
	if err == nil {
		t.Fatal("expected ErrCloudMissingAPIKey, got nil")
	}
	if !errors.Is(err, gollama.ErrCloudMissingAPIKey) {
		t.Fatalf("expected ErrCloudMissingAPIKey, got %v", err)
	}
}

// TestNewClient_CloudURLWithAPIKey verifies that the cloud URL is accepted when an API key is provided.
func TestNewClient_CloudURLWithAPIKey(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithURL("https://ollama.com"), gollama.WithAPIKey("sk-test"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestNewClient_APIKeyWithLineBreak verifies that an API key containing newlines is rejected.
func TestNewClient_APIKeyWithLineBreak(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithAPIKey("key\ninjection"))
	if err == nil {
		t.Fatal("expected ErrAPIKeyHasLineBreaks, got nil")
	}
	if !errors.Is(err, gollama.ErrAPIKeyHasLineBreaks) {
		t.Fatalf("expected ErrAPIKeyHasLineBreaks, got %v", err)
	}
}

// TestNewClient_APIKeyWithCarriageReturn verifies that an API key containing \r is rejected.
func TestNewClient_APIKeyWithCarriageReturn(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithAPIKey("key\rvalue"))
	if err == nil {
		t.Fatal("expected ErrAPIKeyHasLineBreaks, got nil")
	}
	if !errors.Is(err, gollama.ErrAPIKeyHasLineBreaks) {
		t.Fatalf("expected ErrAPIKeyHasLineBreaks, got %v", err)
	}
}

// TestNewClient_LocalClientWithAPIKey verifies that a local client with an API key is allowed.
func TestNewClient_LocalClientWithAPIKey(t *testing.T) {
	_, err := gollama.NewClient(gollama.WithAPIKey("local-key"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
