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

	client, err := gollama.NewClient(gollama.WithLogger(logger))
	if err != nil {
		panic(err)
	}

	response, err := client.ListRunningModels(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("response:", response)
	fmt.Println("response:", len(response))

	fmt.Println("Done!")
}
