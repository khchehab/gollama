package gollama

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// options represent the options for the Ollama client.
type options struct {
	url     string
	apiKey  string
	timeout time.Duration
	logger  *slog.Logger
	client  *http.Client
}

// OptionFunc is a function that modifies the options.
type OptionFunc func(*options)

// applyOptions applies the given options to a new options struct.
func applyOptions(opts ...OptionFunc) *options {
	o := &options{
		url:     "http://localhost:11434",
		timeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithHostAndPort sets the host and port for the Ollama server.
func WithHostAndPort(host string, port int) OptionFunc {
	return func(o *options) {
		o.url = fmt.Sprintf("http://%s:%d", host, port)
	}
}

// WithURL sets the URL for the Ollama server.
func WithURL(url string) OptionFunc {
	return func(o *options) {
		o.url = url
	}
}

// WithAPIKey sets the API key for connecting to the Ollama cloud server.
func WithAPIKey(apiKey string) OptionFunc {
	return func(o *options) {
		o.apiKey = apiKey
	}
}

// WithTimeout sets the timeout for the Ollama client.
func WithTimeout(timeout time.Duration) OptionFunc {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithLogger sets the logger for the Ollama client.
func WithLogger(logger *slog.Logger) OptionFunc {
	return func(o *options) {
		o.logger = logger
	}
}

// WithHTTPClient sets the HTTP client for the Ollama client.
func WithHTTPClient(client *http.Client) OptionFunc {
	return func(o *options) {
		o.client = client
	}
}
