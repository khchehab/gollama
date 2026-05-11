package gollama

import (
	"context"
	"iter"
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

// CreateModel creates a model.
func (c *Client) CreateModel(ctx context.Context, request CreateModelRequest) (*CreateModelResponse, error) {
	internalRequest := createModelRequestWithStream{
		CreateModelRequest: request,
		Stream:             false,
	}
	var response CreateModelResponse
	if err := c.do(ctx, http.MethodPost, "/create", internalRequest, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateModelStream streams creating a model.
func (c *Client) CreateModelStream(ctx context.Context, request CreateModelRequest) iter.Seq2[*CreateModelStatusUpdate, error] {
	return streamResponse[CreateModelStatusUpdate](func() (*http.Response, error) {
		internalRequest := createModelRequestWithStream{
			CreateModelRequest: request,
			Stream:             true,
		}
		return c.executeStream(ctx, http.MethodPost, "/create", internalRequest)
	}, func(line createModelLine) (*CreateModelStatusUpdate, error) {
		if line.Message != "" {
			return nil, &ErrorResponse{Message: line.Message}
		}

		return &CreateModelStatusUpdate{
			Status:    line.Status,
			Digest:    line.Digest,
			Total:     line.Total,
			Completed: line.Completed,
		}, nil
	}, c.logger)
}

// CopyModel copies a model.
func (c *Client) CopyModel(ctx context.Context, request CopyModelRequest) error {
	return c.do(ctx, http.MethodPost, "/copy", request, nil)
}

// PullModel pulls a model.
func (c *Client) PullModel(ctx context.Context, request PullModelRequest) (*PullModelResponse, error) {
	internalRequest := pullModelRequestWithStream{
		PullModelRequest: request,
		Stream:           false,
	}
	var response PullModelResponse
	if err := c.do(ctx, http.MethodPost, "/pull", internalRequest, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// PullModelStream streams pulling a model.
func (c *Client) PullModelStream(ctx context.Context, request PullModelRequest) iter.Seq2[*PullModelStatusUpdate, error] {
	return streamResponse[PullModelStatusUpdate](func() (*http.Response, error) {
		internalRequest := pullModelRequestWithStream{
			PullModelRequest: request,
			Stream:           true,
		}
		return c.executeStream(ctx, http.MethodPost, "/pull", internalRequest)
	}, func(line pullModelLine) (*PullModelStatusUpdate, error) {
		if line.Message != "" {
			return nil, &ErrorResponse{Message: line.Message}
		}

		return &PullModelStatusUpdate{
			Status:    line.Status,
			Digest:    line.Digest,
			Total:     line.Total,
			Completed: line.Completed,
		}, nil
	}, c.logger)
}

// PushModel pushes a model.
func (c *Client) PushModel(ctx context.Context, request PushModelRequest) (*PushModelResponse, error) {
	internalRequest := pushModelRequestWithStream{
		PushModelRequest: request,
		Stream:           false,
	}
	var response PushModelResponse
	if err := c.do(ctx, http.MethodPost, "/push", internalRequest, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// PushModelStream streams pushing a model.
func (c *Client) PushModelStream(ctx context.Context, request PushModelRequest) iter.Seq2[*PushModelStatusUpdate, error] {
	return streamResponse[PushModelStatusUpdate](func() (*http.Response, error) {
		internalRequest := pushModelRequestWithStream{
			PushModelRequest: request,
			Stream:           true,
		}
		return c.executeStream(ctx, http.MethodPost, "/push", internalRequest)
	}, func(line pushModelLine) (*PushModelStatusUpdate, error) {
		if line.Message != "" {
			return nil, &ErrorResponse{Message: line.Message}
		}

		return &PushModelStatusUpdate{
			Status:    line.Status,
			Digest:    line.Digest,
			Total:     line.Total,
			Completed: line.Completed,
		}, nil
	}, c.logger)
}

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
