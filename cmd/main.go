package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/khchehab/gollama"
)

func main() {
	modelName := "nomic-embed-text:latest"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	client, err := gollama.NewClient(gollama.WithLogger(logger))
	if err != nil {
		panic(err)
	}

	stream := client.PullModelStream(context.Background(), gollama.PullModelRequest{
		Model: modelName,
	})

	for update, streamErr := range stream {
		if streamErr != nil {
			panic(streamErr)
		}

		fmt.Println("->", update)
	}

	fmt.Println("Done!")
}
