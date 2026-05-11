package gollama

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// ollamaCloudURL is the default base URL for the Ollama cloud API.
const ollamaCloudURL = "https://ollama.com/api"

// Client represents a client for the Ollama API.
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
	logger  *slog.Logger
}

// NewClient creates a new Ollama client.
func NewClient(opts ...OptionFunc) (*Client, error) {
	o := applyOptions(opts...)

	// validate the base URL
	u, err := url.Parse(o.url)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, &InvalidBaseURLError{BaseURL: o.url}
	}

	// remove any trailing slashes
	u.Path = strings.TrimRight(u.Path, "/")

	// normalize the base URL
	if !strings.HasSuffix(u.Path, "/api") {
		u.Path += "/api"
	}

	baseURL := u.String()

	// using Ollama cloud API requires an API key
	if baseURL == ollamaCloudURL && o.apiKey == "" {
		return nil, ErrCloudMissingAPIKey
	}

	if strings.ContainsAny(o.apiKey, "\r\n") {
		return nil, ErrAPIKeyHasLineBreaks
	}

	httpClient := o.client
	if httpClient == nil {
		httpClient = &http.Client{Timeout: o.timeout}
	}

	var logger *slog.Logger
	if o.logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	} else {
		logger = o.logger
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  o.apiKey,
		client:  httpClient,
		logger:  logger,
	}, nil
}

// do executes an HTTP request against the given endpoint, marshalling request as a JSON body if non-nil, and
// unmarshalling the response body into response if non-nil. Returns an *ErrorResponse if the server responds with a
// non-200 status.
func (c *Client) do(ctx context.Context, method, endpoint string, request any, response any) error {
	endpointURL, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return err
	}

	contentType := ""

	var body io.Reader
	if request != nil {
		var reqBytes []byte
		if reqBytes, err = json.Marshal(request); err != nil {
			return err
		}
		body = bytes.NewReader(reqBytes)

		c.logger.Debug("request body", "body", string(reqBytes))

		contentType = "application/json"
	}

	c.logger.Debug("making request", "method", method, "url", endpointURL, "contentType", contentType)

	req, err := http.NewRequestWithContext(ctx, method, endpointURL, body)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			c.logger.Error("error closing the response body", "error", closeErr)
		}
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	c.logger.Debug("response", "statusCode", res.StatusCode)

	if res.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err = json.Unmarshal(b, &errorResponse); err != nil {
			return err
		}
		return &errorResponse
	}

	if response != nil {
		if err = json.Unmarshal(b, response); err != nil {
			return err
		}
	}

	return nil
}
