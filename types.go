package gollama

// ErrorResponse represents an error response from the Ollama API.
type ErrorResponse struct {
	Message string `json:"error"`
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	return e.Message
}

// GenerateResponse

// GenerateChatMessage

// GenerateEmbeddingRequest represent the request to generate an embedding.
type GenerateEmbeddingRequest struct {
	// Model is the model name.
	Model string `json:"model"`
	// Input is the array of texts to generate embeddings for.
	Input []string `json:"input"`
	// Truncate if true, truncate inputs that exceed the context window. If false, returns an error.
	Truncate bool `json:"truncate"`
	// Dimensions is the number of dimensions to generate embeddings for.
	Dimensions int `json:"dimensions"`
	// KeepAlive is the model keep-alive duration.
	KeepAlive string `json:"keep_alive"`
	// Options is the runtime options that control text generation.
	Options GenerateEmbeddingRequestOption `json:"options"`
}

// GenerateEmbeddingRequestOption represents the runtime options that control the text generation.
type GenerateEmbeddingRequestOption struct {
	// Seed is the random seed used for reproducible outputs.
	Seed int `json:"seed"`
	// Temperature controls randomness in generation (higher = more random).
	Temperature float64 `json:"temperature"`
	// TopK limits next token selection to the K most likely.
	TopK int `json:"top_k"`
	// TopP is the cumulative probability threshold for nucleus sampling.
	TopP float64 `json:"top_p"`
	// MinP is the minimum probability threshold for token selection.
	MinP float64 `json:"min_p"`
	// Stop is the stop sequences that will halt generation.
	Stop []string `json:"stop"`
	// NumCtx is the context length size (number of tokens).
	NumCtx int `json:"num_ctx"`
	// NumPredict is the maximum number of tokens to generate.
	NumPredict int `json:"num_predict"`
}

// GenerateEmbeddingResponse represents the response from generating an embedding.
type GenerateEmbeddingResponse struct {
	// Model is the model that produced the embeddings
	Model string `json:"model"`
	// Embeddings is the array of vector embeddings.
	Embeddings [][]float64 `json:"embeddings"`
	// TotalDuration is the total time spent generating in nanoseconds.
	TotalDuration int `json:"total_duration"`
	// LoadDuration is the load time in nanoseconds.
	LoadDuration int `json:"load_duration"`
	// PromptEvalCount is the number of input tokens processed to generate embeddings.
	PromptEvalCount int `json:"prompt_eval_count"`
}

// listModelResponse represents the response from the list models.
type listModelResponse struct {
	// Models is the list of available models.
	Models []ModelInfo `json:"models"`
}

// ModelInfo represents information about a model.
type ModelInfo struct {
	// Name is the model name.
	Name string `json:"name"`
	// Model is the model name.
	Model string `json:"model"`
	// RemoteModel is the name of the upstream model, if the model is remote.
	RemoteModel string `json:"remote_model"`
	// RemoteHost is the URL of the upstream Ollama host, if the model is remote.
	RemoteHost string `json:"remote_host"`
	// ModifiedAt is the last modified timestamp in ISO 8601 format.
	ModifiedAt string `json:"modified_at"`
	// Size is the total size of the model on disk in bytes.
	Size int `json:"size"`
	// Digest is the SHA256 digest identifier of the model contents.
	Digest string `json:"digest"`
	// Details is the additional information about the model's format and family.
	Details ModelInfoDetail `json:"details"`
}

// ModelInfoDetail represents detailed information about a model.
type ModelInfoDetail struct {
	// Format is the model file format (e.g. gguf).
	Format string `json:"format"`
	// Family is the primary model family (e.g. llama).
	Family string `json:"family"`
	// Families are all the families the model belongs to, when applicable.
	Families []string `json:"families"`
	// ParameterSize is the approximate parameter count label (e.g. 7B, 13B).
	ParameterSize string `json:"parameter_size"`
	// QuantizationLevel is the quantization level used (e.g. Q4_0).
	QuantizationLevel string `json:"quantization_level"`
}

// runningModelsResponse represents the response from the list running models.
type runningModelsResponse struct {
	// Models are the currently running models
	Models []RunningModel `json:"models"`
}

// RunningModel represents a loaded model in memory.
type RunningModel struct {
	// Name is the name of the running model.
	Name string `json:"name"`
	// Model is the name of the running model.
	Model string `json:"model"`
	// Size is the size of the model in bytes.
	Size int `json:"size"`
	// Digest is the SHA256 digest of the model.
	Digest string `json:"digest"`
	// Details is the model details such as format and family.
	Details map[string]any `json:"details"`
	// ExpiresAt is the time when the model will be unloaded.
	ExpiresAt string `json:"expires_at"`
	// SizeVRAM is the VRAM usage in bytes.
	SizeVRAM int `json:"size_vram"`
	// ContextLength is the context length for the running model.
	ContextLength int `json:"context_length"`
}

// ShowModelDetailsRequest represents the request to get the model details.
type ShowModelDetailsRequest struct {
	// Model is the model name to show.
	Model string `json:"model"`
	// Verbose, if true, includes large verbose fields in the response.
	Verbose bool `json:"verbose"`
}

// ShowModelDetailsResponse represents the response from getting model details.
type ShowModelDetailsResponse struct {
	// Parameters is the model parameter settings serialized as text.
	Parameters string `json:"parameters"`
	// License is the license of the model.
	License string `json:"license"`
	// ModifiedAt is the last modified timestamp in ISO 8601 format.
	ModifiedAt string `json:"modified_at"`
	// Details is the high-level model details.
	Details map[string]any `json:"details"`
	// Template is the template used by the model to render prompts.
	Template string `json:"template"`
	// Capabilities is the list of supported features.
	Capabilities []string `json:"capabilities"`
	// ModelInfo is the additional model metadata.
	ModelInfo map[string]any `json:"model_info"`
}

// CreateModel

// CopyModelRequest represents the request to copy a model.
type CopyModelRequest struct {
	// Source is the existing model name to copy from.
	Source string `json:"source"`
	// Destination is the new model name to create.
	Destination string `json:"destination"`
}

// PullModelRequest represents the request to pull a model.
type PullModelRequest struct {
	// Model is the name of the model to download.
	Model string `json:"model"`
	// Insecure allows downloading over insecure connections.
	Insecure bool `json:"insecure"`
}

// pullModelRequestWithStream represents the request to pull a model witht he stream field (to be filled internally based on calling function).
type pullModelRequestWithStream struct {
	PullModelRequest
	// Stream to stream progress updates.
	Stream bool `json:"stream"`
}

// PullModelResponse represents the response from pulling a model.
type PullModelResponse struct {
	// Status is the current status message.
	Status string `json:"status"`
}

// PushModel

// DeleteModelRequest represents the request to delete a model.
type DeleteModelRequest struct {
	// Model is the model name to delete.
	Model string `json:"model"`
}

// versionResponse represents the response from getting the version of the server.
type versionResponse struct {
	// Version is the version of Ollama.
	Version string `json:"version"`
}
