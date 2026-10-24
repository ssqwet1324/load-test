package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ssqwet1324/load-test/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
