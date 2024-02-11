package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/anyuan-chen/hackathon/configs"
	"github.com/gookit/config/v2"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	var handler http.Handler = mux
	return handler
}

func run(ctx context.Context, config configs.Config, args []string) error {
	srv := NewServer()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	//for safe shutdown
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	err := config.LoadFiles("./configs/basic.yml")
	if err != nil {
		fmt.Printf("failed to load the config file")
	}
	err = config.BindStruct()
	if err = run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
