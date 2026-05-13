package gollama_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/khchehab/gollama"
)

// TestListModels_Success verifies that a list of models is returned correctly.
func TestListModels_Success(t *testing.T) {
	body := `{"models":[
		{"name":"llama3","model":"llama3","size":4000000000,"digest":"abc123","details":{"format":"gguf","family":"llama","parameter_size":"7B","quantization_level":"Q4_0"}},
		{"name":"mistral","model":"mistral","size":3000000000,"digest":"def456","details":{"format":"gguf","family":"mistral","parameter_size":"7B","quantization_level":"Q4_0"}}
	]}`

	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/tags" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, body), nil
	}))

	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0].Name != "llama3" {
		t.Errorf("models[0].Name: expected %q, got %q", "llama3", models[0].Name)
	}
	if models[1].Name != "mistral" {
		t.Errorf("models[1].Name: expected %q, got %q", "mistral", models[1].Name)
	}
	if models[0].Details.Family != "llama" {
		t.Errorf("models[0].Details.Family: expected %q, got %q", "llama", models[0].Details.Family)
	}
}

// TestListModels_Empty verifies that an empty model list is returned without error.
func TestListModels_Empty(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, `{"models":[]}`), nil
	}))

	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 0 {
		t.Errorf("expected 0 models, got %d", len(models))
	}
}

// TestListModels_APIError verifies non-200 responses propagate as *ErrorResponse.
func TestListModels_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusInternalServerError, `{"error":"internal error"}`), nil
	}))

	_, err := c.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
}

// TestListRunningModels_Success verifies running models are returned correctly.
func TestListRunningModels_Success(t *testing.T) {
	body := `{"models":[{"name":"llama3","model":"llama3","size":4000000000,"digest":"abc","expires_at":"2024-01-01T01:00:00Z","size_vram":0,"context_length":4096}]}`

	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/ps" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, body), nil
	}))

	models, err := c.ListRunningModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 running model, got %d", len(models))
	}
	if models[0].Name != "llama3" {
		t.Errorf("expected name %q, got %q", "llama3", models[0].Name)
	}
	if models[0].ContextLength != 4096 {
		t.Errorf("expected context_length 4096, got %d", models[0].ContextLength)
	}
}

// TestShowModelDetails_Success verifies model details are returned correctly.
func TestShowModelDetails_Success(t *testing.T) {
	body := `{"parameters":"temperature 0.7","license":"MIT","modified_at":"2024-01-01T00:00:00Z","template":"{{ .Prompt }}","capabilities":["completion"]}`

	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/show" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, body), nil
	}))

	resp, err := c.ShowModelDetails(context.Background(), gollama.ShowModelDetailsRequest{Model: "llama3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.License != "MIT" {
		t.Errorf("expected license %q, got %q", "MIT", resp.License)
	}
	if len(resp.Capabilities) != 1 || resp.Capabilities[0] != "completion" {
		t.Errorf("unexpected capabilities: %v", resp.Capabilities)
	}
}

// TestShowModelDetails_APIError verifies non-200 responses propagate correctly.
func TestShowModelDetails_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusNotFound, `{"error":"model not found"}`), nil
	}))

	_, err := c.ShowModelDetails(context.Background(), gollama.ShowModelDetailsRequest{Model: "missing"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
}

// TestCreateModel_Success verifies non-streaming model creation sends stream=false.
func TestCreateModel_Success(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/create" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("could not decode request: %v", err)
		}
		if req["stream"] != false {
			t.Errorf("expected stream=false, got %v", req["stream"])
		}
		return respond(http.StatusOK, `{"status":"success"}`), nil
	}))

	resp, err := c.CreateModel(context.Background(), gollama.CreateModelRequest{
		Model: "mymodel",
		From:  "llama3",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status %q, got %q", "success", resp.Status)
	}
}

// TestCreateModelStream_Success verifies streaming model creation yields all status updates.
func TestCreateModelStream_Success(t *testing.T) {
	ndjson := `{"status":"pulling manifest"}` + "\n" +
		`{"status":"downloading","digest":"sha256:abc","total":1000,"completed":500}` + "\n" +
		`{"status":"success"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.CreateModelStream(context.Background(), gollama.CreateModelRequest{Model: "mymodel", From: "llama3"})
	var updates []*gollama.CreateModelStatusUpdate
	for update, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		updates = append(updates, update)
	}

	if len(updates) != 3 {
		t.Fatalf("expected 3 updates, got %d", len(updates))
	}
	if updates[0].Status != "pulling manifest" {
		t.Errorf("updates[0].Status: expected %q, got %q", "pulling manifest", updates[0].Status)
	}
	if updates[1].Total != 1000 || updates[1].Completed != 500 {
		t.Errorf("updates[1]: expected total=1000 completed=500, got total=%d completed=%d", updates[1].Total, updates[1].Completed)
	}
}

// TestCreateModelStream_InlineError verifies inline stream errors are propagated.
func TestCreateModelStream_InlineError(t *testing.T) {
	ndjson := `{"status":"pulling manifest"}` + "\n" +
		`{"error":"disk full"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.CreateModelStream(context.Background(), gollama.CreateModelRequest{Model: "mymodel"})
	var updates []*gollama.CreateModelStatusUpdate
	var gotErr error
	for update, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
		updates = append(updates, update)
	}
	if gotErr == nil {
		t.Fatal("expected inline error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(gotErr, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T", gotErr)
	}
	if len(updates) != 1 {
		t.Errorf("expected 1 update before error, got %d", len(updates))
	}
}

// TestCopyModel_Success verifies copying a model succeeds with a 200 response.
func TestCopyModel_Success(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/copy" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, ``), nil
	}))

	err := c.CopyModel(context.Background(), gollama.CopyModelRequest{
		Source:      "llama3",
		Destination: "llama3-copy",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestCopyModel_APIError verifies non-200 responses propagate as *ErrorResponse.
func TestCopyModel_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusNotFound, `{"error":"source model not found"}`), nil
	}))

	err := c.CopyModel(context.Background(), gollama.CopyModelRequest{Source: "missing", Destination: "copy"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
}

// TestPullModel_Success verifies non-streaming model pull sends stream=false.
func TestPullModel_Success(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/pull" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("could not decode request: %v", err)
		}
		if req["stream"] != false {
			t.Errorf("expected stream=false, got %v", req["stream"])
		}
		return respond(http.StatusOK, `{"status":"success"}`), nil
	}))

	resp, err := c.PullModel(context.Background(), gollama.PullModelRequest{Model: "llama3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status %q, got %q", "success", resp.Status)
	}
}

// TestPullModelStream_Success verifies streaming pull yields all status updates.
func TestPullModelStream_Success(t *testing.T) {
	ndjson := `{"status":"pulling manifest"}` + "\n" +
		`{"status":"downloading","digest":"sha256:abc","total":2000,"completed":1000}` + "\n" +
		`{"status":"success"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.PullModelStream(context.Background(), gollama.PullModelRequest{Model: "llama3"})
	var updates []*gollama.PullModelStatusUpdate
	for update, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		updates = append(updates, update)
	}

	if len(updates) != 3 {
		t.Fatalf("expected 3 updates, got %d", len(updates))
	}
	if updates[1].Digest != "sha256:abc" {
		t.Errorf("updates[1].Digest: expected %q, got %q", "sha256:abc", updates[1].Digest)
	}
}

// TestPullModelStream_InlineError verifies inline stream errors are propagated.
func TestPullModelStream_InlineError(t *testing.T) {
	ndjson := `{"status":"pulling manifest"}` + "\n" +
		`{"error":"model not found"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.PullModelStream(context.Background(), gollama.PullModelRequest{Model: "missing"})
	var gotErr error
	for _, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
	}
	if gotErr == nil {
		t.Fatal("expected inline error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(gotErr, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T", gotErr)
	}
}

// TestPushModel_Success verifies non-streaming model push sends stream=false.
func TestPushModel_Success(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/push" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("could not decode request: %v", err)
		}
		if req["stream"] != false {
			t.Errorf("expected stream=false, got %v", req["stream"])
		}
		return respond(http.StatusOK, `{"status":"success"}`), nil
	}))

	resp, err := c.PushModel(context.Background(), gollama.PushModelRequest{Model: "mymodel"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status %q, got %q", "success", resp.Status)
	}
}

// TestPushModelStream_Success verifies streaming push yields all status updates.
func TestPushModelStream_Success(t *testing.T) {
	ndjson := `{"status":"pushing manifest"}` + "\n" +
		`{"status":"uploading","digest":"sha256:abc","total":5000,"completed":2500}` + "\n" +
		`{"status":"success"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.PushModelStream(context.Background(), gollama.PushModelRequest{Model: "mymodel"})
	var updates []*gollama.PushModelStatusUpdate
	for update, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		updates = append(updates, update)
	}

	if len(updates) != 3 {
		t.Fatalf("expected 3 updates, got %d", len(updates))
	}
	if updates[1].Total != 5000 {
		t.Errorf("updates[1].Total: expected 5000, got %d", updates[1].Total)
	}
}

// TestPushModelStream_InlineError verifies inline stream errors are propagated.
func TestPushModelStream_InlineError(t *testing.T) {
	ndjson := `{"status":"pushing manifest"}` + "\n" +
		`{"error":"unauthorized"}` + "\n"

	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusOK, ndjson), nil
	}))

	seq := c.PushModelStream(context.Background(), gollama.PushModelRequest{Model: "mymodel"})
	var gotErr error
	for _, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
	}
	if gotErr == nil {
		t.Fatal("expected inline error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(gotErr, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T", gotErr)
	}
	if errResp.Error() != "unauthorized" {
		t.Errorf("unexpected message %q", errResp.Error())
	}
}

// TestDeleteModel_Success verifies deleting a model succeeds with a 200 response.
func TestDeleteModel_Success(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/delete" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, ``), nil
	}))

	err := c.DeleteModel(context.Background(), gollama.DeleteModelRequest{Model: "llama3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestDeleteModel_APIError verifies non-200 responses propagate as *ErrorResponse.
func TestDeleteModel_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusNotFound, `{"error":"model not found"}`), nil
	}))

	err := c.DeleteModel(context.Background(), gollama.DeleteModelRequest{Model: "missing"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
}

// TestGetVersion_Success verifies that the server version is returned correctly.
func TestGetVersion_Success(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/version" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		return respond(http.StatusOK, `{"version":"0.6.5"}`), nil
	}))

	version, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "0.6.5" {
		t.Errorf("expected version %q, got %q", "0.6.5", version)
	}
}

// TestGetVersion_APIError verifies non-200 responses propagate correctly.
func TestGetVersion_APIError(t *testing.T) {
	c := newClientWithTransport(t, roundTripFunc(func(*http.Request) (*http.Response, error) {
		return respond(http.StatusServiceUnavailable, `{"error":"service unavailable"}`), nil
	}))

	_, err := c.GetVersion(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *gollama.ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
}
