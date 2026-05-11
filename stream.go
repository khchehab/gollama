package gollama

import (
	"io"
	"iter"
)

func streamResponse[T any](r io.Reader) iter.Seq2[T, error] {
	return nil
}
