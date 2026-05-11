package gollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"iter"
	"log/slog"
	"net/http"
)

// streamResponse is a helper function to stream a response. It takes a function to execute the request and get the
// raw HTTP response, a mapper function that maps from the intermediate line struct into the final struct, and finally
// the client's logger to log some information.
func streamResponse[T any, L any](responseExec func() (*http.Response, error), mapper func(L) (*T, error), logger *slog.Logger) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		res, err := responseExec()
		if err != nil {
			yield(nil, err)
			return
		}
		defer func(Body io.ReadCloser) {
			if closeErr := Body.Close(); closeErr != nil {
				logger.Error("error closing the response body", "error", closeErr)
			}
		}(res.Body)

		logger.Debug("stream response", "statusCode", res.StatusCode)

		if res.StatusCode != http.StatusOK {
			var b []byte
			if b, err = io.ReadAll(res.Body); err != nil {
				yield(nil, err)
				return
			}

			var errorResponse ErrorResponse
			if err = json.Unmarshal(b, &errorResponse); err != nil {
				yield(nil, err)
				return
			}

			yield(nil, &errorResponse)
			return
		}

		scanner := bufio.NewScanner(res.Body)

		for scanner.Scan() {
			b := scanner.Bytes()

			if len(bytes.TrimSpace(b)) == 0 {
				continue
			}

			var line L
			if err = json.Unmarshal(b, &line); err != nil {
				yield(nil, err)
				return
			}

			var r *T
			if r, err = mapper(line); err != nil {
				yield(nil, err)
				return
			}

			if !yield(r, nil) {
				return
			}
		}

		if err = scanner.Err(); err != nil {
			yield(nil, err)
			return
		}
	}
}
