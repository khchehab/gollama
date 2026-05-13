package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/khchehab/gollama"
)

func main() {
	modelName := "fauxpaslife/squishy:150m"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	client, err := gollama.NewClient(gollama.WithLogger(logger))
	if err != nil {
		panic(err)
	}

	fmt.Println("===")

	stream := client.PullModelStream(context.Background(), gollama.PullModelRequest{
		Model: modelName,
	})

	for update, streamErr := range stream {
		if streamErr != nil {
			panic(streamErr)
		}

		fmt.Println("->", update)
	}

	fmt.Println("===")

	stream2 := client.GenerateChatMessageStream(context.Background(), gollama.GenerateChatMessageRequest{
		Model: modelName,
		Messages: []gollama.GenerateChatMessageRequestMessage{
			{Role: gollama.RoleUser, Content: "Why is the sky blue?"},
		},
	})

	for update2, streamErr2 := range stream2 {
		if streamErr2 != nil {
			panic(streamErr2)
		}

		fmt.Println("->", update2)
	}

	fmt.Println("Done!")
}
