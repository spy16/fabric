package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spy16/fabric"
)

func main() {
	var db string
	flag.StringVar(&db, "db", "fabric.db", "SQLite databade file path")
	flag.Parse()

	fab := fabric.New(setupStore(db))

	count, err := fab.Count(context.Background(), fabric.Query{})
	if err != nil {
		fmt.Printf("failed to fetch count: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("There are %d triples in the store\n", count)
}

func setupStore(path string) *fabric.SQLStore {
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
