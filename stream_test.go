package gollama

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

// testLine is the intermediate line struct used in tests.
type testLine struct {
	Value  string `json:"value"`
	ErrMsg string `json:"error"`
}

// testItem is the final output struct used in tests.
type testItem struct {
	Value string
}

// testMapper maps a testLine to a *testItem, returning an error if the line has an error message.
func testMapper(line testLine) (*testItem, error) {
	if line.ErrMsg != "" {
		return nil, &ErrorResponse{Message: line.ErrMsg}
	}
	return &testItem{Value: line.Value}, nil
}

var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// makeResponse builds a fake *http.Response with the given status code and body.
func makeResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// collect iterates the sequence and returns all collected items and the first error encountered.
func collect(seq func(yield func(*testItem, error) bool)) ([]*testItem, error) {
	var items []*testItem
	var firstErr error
	for item, err := range seq {
		if err != nil {
			firstErr = err
			break
		}
		items = append(items, item)
	}
	return items, firstErr
}

// TestStreamResponse_HappyPath_SingleItem verifies a single valid ndjson line is yielded correctly.
func TestStreamResponse_HappyPath_SingleItem(t *testing.T) {
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, `{"value":"hello"}`+"\n"), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Value != "hello" {
		t.Errorf("expected value %q, got %q", "hello", items[0].Value)
	}
}

// TestStreamResponse_HappyPath_MultipleItems verifies multiple valid ndjson lines are all yielded.
func TestStreamResponse_HappyPath_MultipleItems(t *testing.T) {
	body := `{"value":"first"}` + "\n" + `{"value":"second"}` + "\n" + `{"value":"third"}` + "\n"
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, body), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	expected := []string{"first", "second", "third"}
	for i, item := range items {
		if item.Value != expected[i] {
			t.Errorf("item[%d]: expected %q, got %q", i, expected[i], item.Value)
		}
	}
}

// TestStreamResponse_EmptyBody verifies that an empty response body yields no items and no error.
func TestStreamResponse_EmptyBody(t *testing.T) {
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, ""), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

// TestStreamResponse_ExecutorError verifies that an error from the executor is yielded and iteration stops.
func TestStreamResponse_ExecutorError(t *testing.T) {
	execErr := errors.New("connection refused")
	exec := func() (*http.Response, error) {
		return nil, execErr
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, execErr) {
		t.Errorf("expected %v, got %v", execErr, err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

// TestStreamResponse_NonOKStatus verifies that a non-200 status yields an *ErrorResponse and stops.
func TestStreamResponse_NonOKStatus(t *testing.T) {
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusInternalServerError, `{"error":"internal server error"}`), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T", err)
	}
	if errResp.Message != "internal server error" {
		t.Errorf("expected message %q, got %q", "internal server error", errResp.Message)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

// TestStreamResponse_NonOKStatus_InvalidJSON verifies that a non-200 response with invalid JSON body yields a decode error.
func TestStreamResponse_NonOKStatus_InvalidJSON(t *testing.T) {
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusBadRequest, `not json`), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	_, err := collect(seq)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *ErrorResponse
	if errors.As(err, &errResp) {
		t.Error("expected a JSON decode error, not an *ErrorResponse")
	}
}

// TestStreamResponse_InvalidLineJSON verifies that a malformed ndjson line yields a decode error and stops iteration.
func TestStreamResponse_InvalidLineJSON(t *testing.T) {
	body := `{"value":"first"}` + "\n" + `not json` + "\n"
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, body), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// first item should have been yielded before the bad line
	if len(items) != 1 {
		t.Errorf("expected 1 item before error, got %d", len(items))
	}
}

// TestStreamResponse_MapperError verifies that a mapper error is yielded and iteration stops.
func TestStreamResponse_MapperError(t *testing.T) {
	body := `{"value":"first"}` + "\n" + `{"error":"model not found"}` + "\n"
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, body), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var errResp *ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T", err)
	}
	if errResp.Message != "model not found" {
		t.Errorf("expected message %q, got %q", "model not found", errResp.Message)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item before mapper error, got %d", len(items))
	}
}

// TestStreamResponse_EarlyCancellation verifies that returning false from yield stops iteration cleanly.
func TestStreamResponse_EarlyCancellation(t *testing.T) {
	body := `{"value":"first"}` + "\n" + `{"value":"second"}` + "\n" + `{"value":"third"}` + "\n"
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, body), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)

	var items []*testItem
	for item, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		items = append(items, item)
		break // cancel after first item
	}

	if len(items) != 1 {
		t.Errorf("expected 1 item after early cancellation, got %d", len(items))
	}
	if items[0].Value != "first" {
		t.Errorf("expected %q, got %q", "first", items[0].Value)
	}
}

// TestStreamResponse_TrailingNewlines verifies that extra blank lines in the body are silently skipped.
func TestStreamResponse_TrailingNewlines(t *testing.T) {
	body := `{"value":"only"}` + "\n\n\n"
	exec := func() (*http.Response, error) {
		return makeResponse(http.StatusOK, body), nil
	}

	seq := streamResponse[testItem, testLine](exec, testMapper, discardLogger)
	items, err := collect(seq)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Value != "only" {
		t.Errorf("expected %q, got %q", "only", items[0].Value)
	}
}
