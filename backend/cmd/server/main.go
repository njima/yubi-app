package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "stop server: %v\n", err)
		os.Exit(1)
	}
}
