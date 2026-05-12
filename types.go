package gollama

import "time"

// ErrorResponse represents an error response from the Ollama API.
type ErrorResponse struct {
	// Message is the error message.
	Message string `json:"error"`
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	return e.Message
}

// GenerateResponseRequest represents the request to generate a response.
type GenerateResponseRequest struct {
	// Model is the model name.
	Model string `json:"model"`
}

// generateResponseRequestWithStream represents the request to generate a response with the stream field (to be filled internally based on calling function).
type generateResponseRequestWithStream struct {
	GenerateResponseRequest
	// Stream to stream progress updates.
	Stream bool `json:"stream"`
}

// GenerateResponseResponse represents the response of the generate response function.
type GenerateResponseResponse struct {
	// Model is the model name.
	Model string `json:"model"`
	// CreatedAt is the ISO 8601 timestamp of response creation.
	CreatedAt string `json:"created_at"`
	// Response is the model's generated text response.
	Response string `json:"response"`
	// Thinking is the model's generated thinking output.
	Thinking string `json:"thinking"`
	// Done indicates whether generation has finished.
	Done bool `json:"done"`
	// DoneReason is the reason the generation stopped.
	DoneReason string `json:"done_reason"`
	// TotalDuration is the time spent generating the response.
	TotalDuration time.Duration `json:"total_duration"`
	// LoadDuration is the time spent loading the model.
	LoadDuration time.Duration `json:"load_duration"`
	// PromptEvalCount is the number of input tokens in the prompt.
	PromptEvalCount int `json:"prompt_eval_count"`
	// PromptEvalDuration is the time spent evaluating the prompt.
	PromptEvalDuration time.Duration `json:"prompt_eval_duration"`
	// EvalCount is the number of output tokens generated in the response.
	EvalCount int `json:"eval_count"`
	// EvalDuration is the time spent generating tokens.
	EvalDuration time.Duration `json:"eval_duration"`
	// Logprobs is the log probability information for the generated tokens when logprobs are enabled.
	Logprobs []GenerateResponseResponseLogprobs `json:"logprobs"`
}

// GenerateResponseResponseLogprobs represents the log probability information.
type GenerateResponseResponseLogprobs struct {
	// Token is the text representation of the token.
	Token string `json:"token"`
	// Logprob is the log probability of this token.
	Logprob string `json:"logprob"`
	// Bytes is the raw byte representation of the token.
	Bytes []byte `json:"bytes"`
	// TopLogprobs are the most likely tokens and their log probabilities at this position.
	TopLogprobs []GenerateResponseResponseLogprobsTopLogprobs `json:"top_logprobs"`
}

// GenerateResponseResponseLogprobsTopLogprobs represents the most likely tokens and their log probabilities.
type GenerateResponseResponseLogprobsTopLogprobs struct {
	// Token is the text representation of the token.
	Token string `json:"token"`
	// Logprob is the log probability of this token.
	Logprob string `json:"logprob"`
	// Bytes is the raw byte representation of the token.
	Bytes []byte `json:"bytes"`
}

// GenerateResponseChunk represents a chunk of the generate response functionality.
type GenerateResponseChunk struct {
	// Model is the model name.
	Model string `json:"model"`
	// CreatedAt is the iSO 8601 timestamp of response creation.
	CreatedAt string `json:"created_at"`
	// Response is the the model's generated text response for this chunk.
	Response string `json:"response"`
	// Thinking is the the model's generated thinking output for this chunk.
	Thinking string `json:"thinking"`
	// Done indicates whether the stream has finished.
	Done bool `json:"done"`
	// DoneReason is the reason streaming finished.
	DoneReason string `json:"done_reason"`
	// TotalDuration is the time spent generating the response.
	TotalDuration time.Duration `json:"total_duration"`
	// LoadDuration is the time spent loading the model.
	LoadDuration time.Duration `json:"load_duration"`
	// PromptEvalCount is the number of input tokens in the prompt.
	PromptEvalCount int `json:"prompt_eval_count"`
	// PromptEvalDuration is the time spent evaluating the prompt.
	PromptEvalDuration time.Duration `json:"prompt_eval_duration"`
	// EvalCount is the number of output tokens generated in the response.
	EvalCount int `json:"eval_count"`
	// EvalDuration is the time spent generating tokens.
	EvalDuration int `json:"eval_duration"`
}

// generateResponseChunkLine represents a line of progress during a streamed generate response (can be either a chunk or error response).
type generateResponseChunkLine struct {
	// Model is the model name.
	Model string `json:"model"`
	// CreatedAt is the iSO 8601 timestamp of response creation.
	CreatedAt string `json:"created_at"`
	// Response is the the model's generated text response for this chunk.
	Response string `json:"response"`
	// Thinking is the the model's generated thinking output for this chunk.
	Thinking string `json:"thinking"`
	// Done indicates whether the stream has finished.
	Done bool `json:"done"`
	// DoneReason is the reason streaming finished.
	DoneReason string `json:"done_reason"`
	// TotalDuration is the time spent generating the response.
	TotalDuration time.Duration `json:"total_duration"`
	// LoadDuration is the time spent loading the model.
	LoadDuration time.Duration `json:"load_duration"`
	// PromptEvalCount is the number of input tokens in the prompt.
	PromptEvalCount int `json:"prompt_eval_count"`
	// PromptEvalDuration is the time spent evaluating the prompt.
	PromptEvalDuration time.Duration `json:"prompt_eval_duration"`
	// EvalCount is the number of output tokens generated in the response.
	EvalCount int `json:"eval_count"`
	// EvalDuration is the time spent generating tokens.
	EvalDuration int `json:"eval_duration"`
	// Message is the error message.
	Message string `json:"error"`
}

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
	// TopK limits the next token selection to the K most likely.
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
	// TotalDuration is the total time spent generating.
	TotalDuration time.Duration `json:"total_duration"`
	// LoadDuration is the load time.
	LoadDuration time.Duration `json:"load_duration"`
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
	// Size is the total size of the model on the disk in bytes.
	Size int64 `json:"size"`
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
	Size int64 `json:"size"`
	// Digest is the SHA256 digest of the model.
	Digest string `json:"digest"`
	// Details is the model details such as format and family.
	Details map[string]any `json:"details"`
	// ExpiresAt is the time when the model will be unloaded.
	ExpiresAt string `json:"expires_at"`
	// SizeVRAM is the VRAM usage in bytes.
	SizeVRAM int64 `json:"size_vram"`
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

// CreateModelRequest represents the request to create a model.
type CreateModelRequest struct {
	// Model is the name of the model to create.
	Model string `json:"model"`
	// From is the existing model to create from.
	From string `json:"from"`
	// Template is the prompt template to use for the model.
	Template string `json:"template"`
	// License is the list of licenses for the model.
	License []string `json:"license"`
	// System is the system prompt to embed in the model.
	System string `json:"system"`
	// Parameters is the key-value parameters for the model.
	Parameters map[string]any `json:"parameters"`
	// Messages is the message history to use for the model.
	Messages []CreateModelRequestMessage `json:"messages"`
	// Quantize is the quantization level to apply (e.g. q4_K_M, q8_0).
	Quantize string `json:"quantize"`
}

// CreateModelRequestMessage represents a message history to use for the model.
type CreateModelRequestMessage struct {
	// Role is the author of the message.
	Role MessageRole `json:"role"`
	// Content is the message text content.
	Content string `json:"content"`
	// Images is an optional list of inline images for multimodal models.
	// Each item will be a base64-encoded image content.
	Images []string `json:"images"`
	// ToolCalls is the tool call requests produced by the model.
	ToolCalls []CreateModelRequestMessageToolCall `json:"tool_calls"`
}

// CreateModelRequestMessageToolCall represents a message tool call.
type CreateModelRequestMessageToolCall struct {
	// Function is the function call details.
	Function CreateModelRequestMessageToolCallFunction `json:"function"`
}

// CreateModelRequestMessageToolCallFunction represents the function call details.
type CreateModelRequestMessageToolCallFunction struct {
	// Name of the function to call.
	Name string `json:"name"`
	// Description is what the function does.
	Description string `json:"description"`
	// Arguments is a map of arguments to pass to the function.
	Arguments map[string]any `json:"arguments"`
}

// createModelRequestWithStream represents the request to create a model with the stream field (to be filled internally based on calling function).
type createModelRequestWithStream struct {
	CreateModelRequest
	// Stream to stream progress updates.
	Stream bool `json:"stream"`
}

// CreateModelResponse represents the response from creating a model.
type CreateModelResponse struct {
	// Status is the current status message.
	Status string `json:"status"`
}

// CreateModelStatusUpdate represents a status update during a streamed create model.
type CreateModelStatusUpdate struct {
	// Status is a human-readable status message.
	Status string `json:"status"`
	// Digest is the content digest associated with the status, if applicable.
	Digest string `json:"digest"`
	// Total is the total number of bytes expected for the operation.
	Total int64 `json:"total"`
	// Completed is the number of bytes transferred so far.
	Completed int64 `json:"completed"`
}

// createModelLine represents a line of progress during a streamed create model (can be either a status update or error response).
type createModelLine struct {
	// Status is a human-readable status message.
	Status string `json:"status"`
	// Digest is the content digest associated with the status, if applicable.
	Digest string `json:"digest"`
	// Total is the total number of bytes expected for the operation.
	Total int64 `json:"total"`
	// Completed is the number of bytes transferred so far.
	Completed int64 `json:"completed"`
	// Message is the error message.
	Message string `json:"error"`
}

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

// pullModelRequestWithStream represents the request to pull a model with the stream field (to be filled internally based on calling function).
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

// PullModelStatusUpdate represents a status update during a streamed pull model.
type PullModelStatusUpdate struct {
	// Status is a human-readable status message.
	Status string `json:"status"`
	// Digest is the content digest associated with the status, if applicable.
	Digest string `json:"digest"`
	// Total is the total number of bytes expected for the operation.
	Total int64 `json:"total"`
	// Completed is the number of bytes transferred so far.
	Completed int64 `json:"completed"`
}

// pullModelLine represents a line of progress during a streamed pull model (can be either a status update or error response).
type pullModelLine struct {
	// Status is a human-readable status message.
	Status string `json:"status"`
	// Digest is the content digest associated with the status, if applicable.
	Digest string `json:"digest"`
	// Total is the total number of bytes expected for the operation.
	Total int64 `json:"total"`
	// Completed is the number of bytes transferred so far.
	Completed int64 `json:"completed"`
	// Message is the error message.
	Message string `json:"error"`
}

// PushModelRequest represents the request to push a model.
type PushModelRequest struct {
	// Model is the name of the model to publish.
	Model string `json:"model"`
	// Insecure allows publishing over insecure connections.
	Insecure bool `json:"insecure"`
}

// pushModelRequestWithStream represents the request to push a model with the stream field (to be filled internally based on calling function).
type pushModelRequestWithStream struct {
	PushModelRequest
	// Stream to stream progress updates.
	Stream bool `json:"stream"`
}

// PushModelResponse represents the response from pushing a model.
type PushModelResponse struct {
	// Status is the current status message.
	Status string `json:"status"`
}

// PushModelStatusUpdate represents a status update during a streamed push model.
type PushModelStatusUpdate struct {
	// Status is a human-readable status message.
	Status string `json:"status"`
	// Digest is the content digest associated with the status, if applicable.
	Digest string `json:"digest"`
	// Total is the total number of bytes expected for the operation.
	Total int64 `json:"total"`
	// Completed is the number of bytes transferred so far.
	Completed int64 `json:"completed"`
}

// pushModelLine represents a line of progress during a streamed push model (can be either a status update or error response).
type pushModelLine struct {
	// Status is a human-readable status message.
	Status string `json:"status"`
	// Digest is the content digest associated with the status, if applicable.
	Digest string `json:"digest"`
	// Total is the total number of bytes expected for the operation.
	Total int64 `json:"total"`
	// Completed is the number of bytes transferred so far.
	Completed int64 `json:"completed"`
	// Message is the error message.
	Message string `json:"error"`
}

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
