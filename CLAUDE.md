# CLAUDE.md

## Project

Go client library for the Ollama REST API. Targets programmatic use in Go applications.

**Module path:** github.com/khchehab/gollama  
**Go version:** 1.22  
**Test framework:** stdlib `testing` only

## Developer Preferences

- Do not provide full implementations unless explicitly asked.
- Pseudocode is acceptable when illustrating a concept or approach.
- Do not oversimplify - assume strong Go proficiency.
- No filler words or unnecessary preamble in responses.
- When multiple design options exist, present them with trade-offs and let the developer decide.
- Do not add convenience wrappers or extra abstractions unless explicitly requested.
- Do not suggest external dependencies - stdlib only unless explicitly approved.

## Scope

Covered API surfaces:

- Generate (completions) - including streaming
- Chat - including streaming
- Embeddings
- Model management - list, pull, push, delete

Out of scope (to be tackled later): OpenAI compatibility, Anthropic compatibility, Blobs API.

## Project Structure

```
gollama/
├── go.mod
├── go.sum
├── CLAUDE.md
├── README.md
├── client.go             # Client struct, constructor, post-construction validation
├── options.go            # Functional options: WithHostAndPort, WithURL, WithAPIKey, WithTimeout, WithLogger, WithHTTPClient
├── types.go              # All request/response structs
├── errors.go             # Typed errors
├── endpoints.go          # Endpoint implementations
├── stream_test.go        # package gollama       - white-box testing for streamResponse[T]
└── gollama_test/         # package gollama_test  - black-box
   ├── client_test.go
   ├── generate_test.go
   ├── chat_test.go
   ├── embeddings_test.go
   └── models_test.go
```

## Client Design

- Constructor: `NewClient(opts ...OptionFunc) (*Client, error)`
- Options:
    - `WithHostAndPort(host string, port int)` - sets base URL to `http://{host}:{port}`
    - `WithURL(url string)` - full base URL override
    - `WithAPIKey(key string)` - sets Bearer token for Authorization header
    - `WithTimeout(d time.Duration)` - sets HTTP client timeout
    - `WithLogger(l *slog.Logger)` - sets logger for client
    - `WithHTTPClient(c *http.Client)` - replaces the default HTTP client (used for transport injection in tests)
- Defaults (no options provided): base URL `http://localhost:11434`, no auth.
- Cloud detection: if the resolved base URL is `https://ollama.com/api` and no API key is set, `NewClient` returns an
  error.
- API key on a local client is permitted (supports cloud model access via a local Ollama instance).

## Streaming

- Primary streaming API: `iter.Seq2[T, error]` (Go 1.22 range-over-func).
- Internal implementation: generic `streamResponse[T any](resp *http.Response) iter.Seq2[T, error]` in `endpoints.go`.
- Error contract: on error, yield `(zero, err)` and stop iteration.
- Early cancellation: caller returns `false` from the yield function; iterator exits cleanly.
- Response body lifecycle: closed via `defer` inside `streamResponse` - caller has no cleanup responsibility.
- Generate and Chat each have a dedicated streaming method returning `iter.Seq2[T, error]` where `T` is their respective
  response type.

## Conventions

- All exported types and functions must have godoc comments.
- Errors must be typed - no raw `errors.New` for API-level failures.
- No global state - all configuration lives on the `Client` struct.
- `context.Context` must be threaded through every public method.
- HTTP transport is configurable via `WithHTTPClient` to allow test injection.

## Testing

- White-box (`package gollama`): internal functions not reachable from outside the package - primarily `streamResponse`
  and any unexported helpers.
- Black-box (`package gollama_test`): all exported methods and types, tested as a consumer would use them.
- HTTP transport mocked via a custom `http.RoundTripper` - no live Ollama instance required.

## Commands

```bash
go test ./...           # Run all tests
go test -race ./...     # Run with race detector
go vet ./...            # Static analysis
```
