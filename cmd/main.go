package main

import (
	"fmt"
	"iter"
)

func main() {
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// }))

	// logger.Info("Hello World!")

	// client, err := gollama.NewClient(gollama.WithLogger(logger))
	// if err != nil {
	// 	panic(err)
	// }

	// response, err := client.PullModel(context.Background(), gollama.PullModelRequest{
	// 	Model: "nomic-embed-text:latest",
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("response:", response)

	// fmt.Println("Done!")

	for i := range Count(10) {
		fmt.Println("range over - i:", i)

		if i > 5 {
			break
		}
	}
}

func Count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; i < n; i++ {
			fmt.Println("DEBUG - i inside the yield:", i)
			if !yield(i) {
				fmt.Println("DEBUG - not yielded")
				return
			}
		}
	}
}
