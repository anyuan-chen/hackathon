package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/anyuan-chen/hackathon/configs"
	"github.com/anyuan-chen/hackathon/src/hackers"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/*
* new server, add routes
 */
func NewServer(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	hackers.AddRoutes(r, db)
	var handler http.Handler = r
	return handler
}

func run(ctx context.Context) error {
	//set up safe interrupt
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	//load config file
	config.AddDriver(yamlv3.Driver)
	err := config.LoadFiles("../configs/basic.yml")
	if err != nil {
		print("failed to load the config file", err.Error())
		return err
	}
	conf := configs.Config{}
	err = config.Decode(&conf)
	if err != nil {
		fmt.Printf("failed to decode config file")
		return err
	}
	// set up sql db
	db, err := sql.Open("postgres", conf.DatabaseUrl)
	if err != nil {
		fmt.Printf("unable to open database")
		return err
	}

	// set up http server
	srv := NewServer(db)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(conf.Host, conf.Port),
		Handler: srv,
	}

	// serve http server
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
	//for safe shutdown
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
