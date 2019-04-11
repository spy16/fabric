> WIP

# Fabric

Fabric is a triple-store written in `Golang`. Fabric provides simple functions
and store options to deal with "Subject->Predicate->Object" relations or so called
triples.

## Usage

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

fab := fabric.New(store)
```

> Fabric `SQLStore` uses Golang standard `database/sql` package. So any SQL database supported
> through this interface (includes most major SQL databases) should work.
