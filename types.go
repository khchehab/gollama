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

// GenerateEmbeddingRequest represent the request to the /embed endpoint.
type GenerateEmbeddingRequest struct {
	// Model is the model name.
	Model string `json:"model"`
	// Input is the array of texts to generate embeddings for.
	Input []string `json:"input"`
	// Truncate if true, truncate inputs that exceed the context window. If false, returns an error.
	Truncate bool `json:"truncate,omitempty"`
	// Dimensions is the number of dimensions to generate embeddings for.
	Dimensions int `json:"dimensions,omitempty"`
	// KeepAlive is the model keep-alive duration.
	KeepAlive string `json:"keep_alive,omitempty"`
	// Options is the runtime options that control text generation.
	Options GenerateEmbeddingRequestOption `json:"options,omitempty"`
}

type GenerateEmbeddingRequestOption struct {
	Seed        int      `json:"seed,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
	TopK        int      `json:"top_k,omitempty"`
	TopP        float64  `json:"top_p,omitempty"`
	MinP        float64  `json:"min_p,omitempty"`
	Stop        []string `json:"stop,omitempty"`
	NumCtx      int      `json:"num_ctx,omitempty"`
	NumPredict  int      `json:"num_predict,omitempty"`
}

// GenerateEmbeddingResponse represents the response from the /embed endpoint.
type GenerateEmbeddingResponse struct {
}

// listModelResponse represents the response from the /tags endpoint.
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

// runningModelsResponse represents the response from the /ps endpoint.
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

// ShowModelDetailsRequest represents the request to the /show endpoint.
type ShowModelDetailsRequest struct {
	// Model is the model name to show.
	Model string `json:"model"`
	// Verbose, if true, includes large verbose fields in the response.
	Verbose bool `json:"verbose,omitempty"`
}

// ShowModelDetailsResponse represents the response from the /show endpoint.
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

// PullModel

// PushModel

// DeleteModelRequest represents the request to delete a model.
type DeleteModelRequest struct {
	// Model is the model name to delete.
	Model string `json:"model"`
}

// versionResponse represents the response from the /version endpoint.
type versionResponse struct {
	// Version is the version of Ollama.
	Version string `json:"version"`
}
