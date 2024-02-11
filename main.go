package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
