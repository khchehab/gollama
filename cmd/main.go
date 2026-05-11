package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/khchehab/gollama"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("Hello World!")

	client, err := gollama.NewClient(gollama.WithLogger(logger))
	if err != nil {
		panic(err)
	}

	response, err := client.GenerateEmbeddings(context.Background(), gollama.GenerateEmbeddingRequest{
		Model: "nomic-embed-text:latest",
		Input: []string{"Query: Hello"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("response:", response)

	fmt.Println("Done!")
}
