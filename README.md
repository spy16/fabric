> WIP

# Fabric


[![GoDoc](https://godoc.org/github.com/spy16/fabric?status.svg)](https://godoc.org/github.com/spy16/fabric) [![Go Report Card](https://goreportcard.com/badge/github.com/spy16/fabric)](https://goreportcard.com/report/github.com/spy16/fabric)

Fabric is a triple-store written in `Go`. Fabric provides simple functions
and store options to deal with "Subject->Predicate->Object" relations or so called
triples.

## Usage

Get fabric by using `go get -u github.com/spy16/fabric` (Fabric as a library has no external dependencies)

```go
mem := &fabric.InMemoryStore{}

fab := fabric.New(mem)

fab.Insert(context.Background(), fabric.Triple{
    Source: "Bob",
    Predicate: "Knows",
    Target: "John",
})

fab.Query(context.Background(), fabric.Query{
    Source: fabric.Clause{
        Type: "equal",
        Value: "Bob",
    },
})
```

To use a SQL database for storing the triples, use the following snippet:

```go
db, err := sql.Open("sqlite3", "fabric.db")
if err != nil {
    panic(err)
}

store := &fabric.SQLStore{
    DB: db,
}
store.Setup(context.Background()) // to create required tables

fab := fabric.New(store)
```

> Fabric `SQLStore` uses Go's standard `database/sql` package. So any SQL database supported
> through this interface (includes most major SQL databases) should work.

Additional store support can be added by implementing the `Store` interface.

```go
type Store interface {
	Insert(ctx context.Context, tri Triple) error
	Query(ctx context.Context, q Query) ([]Triple, error)
	Delete(ctx context.Context, q Query) (int, error)
}
```

Optional `Counter` and `ReWeighter` can be implemented by the store implementations
to support extended query options.
