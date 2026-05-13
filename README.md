# gollama

[![Go Reference](https://pkg.go.dev/badge/github.com/khchehab/gollama.svg)](https://pkg.go.dev/github.com/khchehab/gollama)
[![Go Report Card](https://goreportcard.com/badge/github.com/khchehab/gollama)](https://goreportcard.com/report/github.com/khchehab/gollama)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/khchehab/gollama/main)](https://go.dev/doc/install)

A Go client library for the [Ollama](https://ollama.com) REST API.

## Requirements

- Go 1.26

## Installation

```bash
go get github.com/khchehab/gollama
```

## Usage

### Creating a client

```go
// Default: connects to http://localhost:11434
client, err := gollama.NewClient()

// Custom host and port
client, err := gollama.NewClient(gollama.WithHostAndPort("myhost", 9999))

// Full URL override
client, err := gollama.NewClient(gollama.WithURL("http://myhost:9999"))

// Ollama cloud (API key required)
client, err := gollama.NewClient(
    gollama.WithURL("https://ollama.com"),
    gollama.WithAPIKey("sk-..."),
)
```

### Generate a response

```go
resp, err := client.GenerateResponse(ctx, gollama.GenerateResponseRequest{
    Model:  "llama3",
    Prompt: "Why is the sky blue?",
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Response)
```

### Stream a response

```go
for chunk, err := range client.GenerateResponseStream(ctx, gollama.GenerateResponseRequest{
    Model:  "llama3",
    Prompt: "Why is the sky blue?",
}) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(chunk.Response)
}
```

### Chat

```go
resp, err := client.GenerateChatMessage(ctx, gollama.GenerateChatMessageRequest{
    Model: "llama3",
    Messages: []gollama.GenerateChatMessageRequestMessage{
        {Role: gollama.RoleUser, Content: "Hello!"},
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Message.Content)
```

### Stream a chat response

```go
for chunk, err := range client.GenerateChatMessageStream(ctx, gollama.GenerateChatMessageRequest{
    Model: "llama3",
    Messages: []gollama.GenerateChatMessageRequestMessage{
        {Role: gollama.RoleUser, Content: "Hello!"},
    },
}) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(chunk.Message.Content)
}
```

### Generate embeddings

```go
resp, err := client.GenerateEmbeddings(ctx, gollama.GenerateEmbeddingRequest{
    Model: "nomic-embed-text",
    Input: []string{"hello world", "foo bar"},
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Embeddings)
```

### Model management

```go
// List available models
models, err := client.ListModels(ctx)

// List running models
running, err := client.ListRunningModels(ctx)

// Show model details
details, err := client.ShowModelDetails(ctx, gollama.ShowModelDetailsRequest{Model: "llama3"})

// Pull a model (non-streaming)
resp, err := client.PullModel(ctx, gollama.PullModelRequest{Model: "llama3"})

// Pull a model (streaming progress)
for update, err := range client.PullModelStream(ctx, gollama.PullModelRequest{Model: "llama3"}) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s %d/%d\n", update.Status, update.Completed, update.Total)
}

// Push a model
resp, err := client.PushModel(ctx, gollama.PushModelRequest{Model: "mymodel"})

// Copy a model
err := client.CopyModel(ctx, gollama.CopyModelRequest{Source: "llama3", Destination: "llama3-custom"})

// Create a model
resp, err := client.CreateModel(ctx, gollama.CreateModelRequest{
    Model: "mymodel",
    From:  "llama3",
    System: "You are a helpful assistant.",
})

// Delete a model
err := client.DeleteModel(ctx, gollama.DeleteModelRequest{Model: "llama3"})

// Get server version
version, err := client.GetVersion(ctx)
```

## Client options

| Option                        | Description                                                                      |
|-------------------------------|----------------------------------------------------------------------------------|
| `WithHostAndPort(host, port)` | Set base URL to `http://{host}:{port}`                                           |
| `WithURL(url)`                | Full base URL override                                                           |
| `WithAPIKey(key)`             | Set Bearer token for the `Authorization` header                                  |
| `WithTimeout(d)`              | Set HTTP client timeout (default: 30s; streaming requests always use no timeout) |
| `WithLogger(l)`               | Set a `*slog.Logger` for debug/error logging                                     |
| `WithHTTPClient(c)`           | Replace the default HTTP client (useful for transport injection in tests)        |

## Error handling

API-level errors are returned as `*gollama.ErrorResponse`, which implements `error`. Client construction errors are either sentinel values or typed errors:

```go
var errResp *gollama.ErrorResponse
if errors.As(err, &errResp) {
    fmt.Println("API error:", errResp.Error())
}

// Sentinel errors
errors.Is(err, gollama.ErrCloudMissingAPIKey)
errors.Is(err, gollama.ErrAPIKeyHasLineBreaks)

// Typed construction error
var urlErr *gollama.InvalidBaseURLError
errors.As(err, &urlErr)
```

## Streaming

Streaming methods return `iter.Seq2[*T, error]` (Go 1.23 range-over-func). Iteration stops on the first error. Early cancellation is supported by returning `false` from the loop body (i.e. using `break`). The response body is closed automatically — the caller has no cleanup responsibility.

## Development

```bash
go test ./...        # run all tests
go test -race ./...  # run with race detector
go vet ./...         # static analysis
```

## License

MIT
