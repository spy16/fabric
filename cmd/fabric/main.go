package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spy16/fabric"
	"github.com/spy16/fabric/server"
)

var (
	db       = flag.String("store", ":memory:", "Storage location")
	httpAddr = flag.String("http", ":8080", "HTTP Server Address")
)

func main() {
	flag.Parse()

	fab := fabric.New(setupStore(*db))
	mux := server.NewHTTP(fab)
	log.Printf("starting HTTP API server on '%s'...", *httpAddr)
	log.Fatalf("server exiting: %v", http.ListenAndServe(*httpAddr, mux))
}

func setupStore(path string) fabric.Store {
	if path == ":memory:" {
		return &fabric.InMemoryStore{}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to open db: %v\n", err)
	}

	store := &fabric.SQLStore{
		DB: db,
	}

	if err := store.Setup(context.Background()); err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	return store
}
