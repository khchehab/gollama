package gollama

import (
	"errors"
	"fmt"
)

var (
	// ErrCloudMissingAPIKey is returned when the Ollama cloud URL is used without an API key.
	ErrCloudMissingAPIKey = errors.New("an API key is required when using the Ollama cloud URL")
	// ErrAPIKeyHasLineBreaks is returned when the API key contains line breaks.
	ErrAPIKeyHasLineBreaks = errors.New("apiKey must not contain line breaks")
)

// InvalidBaseURLError is returned when the base URL is invalid.
type InvalidBaseURLError struct {
	BaseURL string
}

// Error implements the error interface.
func (e *InvalidBaseURLError) Error() string {
	return fmt.Sprintf("invalid base URL %q: must be an absolute http/https URL", e.BaseURL)
}
