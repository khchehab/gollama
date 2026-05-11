package gollama

import (
	"context"
	"net/http"
)

// GenerateResponse

// GenerateChatMessage

// GenerateEmbeddings creates vector embeddings representing the input text.
func (c *Client) GenerateEmbeddings(ctx context.Context, request GenerateEmbeddingRequest) (*GenerateEmbeddingResponse, error) {
	var response GenerateEmbeddingResponse
	if err := c.do(ctx, http.MethodPost, "/embed", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ListModels fetch a list of models and their details.
func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error) {
	var response listModelResponse
	if err := c.do(ctx, http.MethodGet, "/tags", nil, &response); err != nil {
		return nil, err
	}
	return response.Models, nil
}

// ListRunningModels retrieve a list of models that are currently running.
func (c *Client) ListRunningModels(ctx context.Context) ([]RunningModel, error) {
	var response runningModelsResponse
	if err := c.do(ctx, http.MethodGet, "/ps", nil, &response); err != nil {
		return nil, err
	}
	return response.Models, nil
}

// ShowModelDetails retrieves the model details.
func (c *Client) ShowModelDetails(ctx context.Context, request ShowModelDetailsRequest) (*ShowModelDetailsResponse, error) {
	var response ShowModelDetailsResponse
	if err := c.do(ctx, http.MethodPost, "/show", request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateModel

// CopyModel copies a model.
func (c *Client) CopyModel(ctx context.Context, request CopyModelRequest) error {
	return c.do(ctx, http.MethodPost, "/copy", request, nil)
}

// PullModel

// PushModel

// DeleteModel deletes a model.
func (c *Client) DeleteModel(ctx context.Context, request DeleteModelRequest) error {
	return c.do(ctx, http.MethodDelete, "/delete", request, nil)
}

// GetVersion retrieve the version of the Ollama.
func (c *Client) GetVersion(ctx context.Context) (string, error) {
	var response versionResponse
	if err := c.do(ctx, http.MethodGet, "/version", nil, &response); err != nil {
		return "", err
	}
	return response.Version, nil
}
