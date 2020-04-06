package fabric

import (
	"context"
	"errors"
)

// New returns a new instance of fabric with given store implementation.
func New(store Store) *Fabric {
	if f, ok := store.(*Fabric); ok {
		return f
	}

	f := &Fabric{}
	f.store = store
	return f
}

// Fabric provides functions to query and manage triples.
type Fabric struct {
	store Store
}

// Insert validates the triple and persists it to the store.
func (f *Fabric) Insert(ctx context.Context, tri Triple) error {
	if err := tri.Validate(); err != nil {
		return err
	}

	return f.store.Insert(ctx, tri)
}

// Query finds all the triples matching the given query.
func (f *Fabric) Query(ctx context.Context, query Query) ([]Triple, error) {
	query.normalize()
	return f.store.Query(ctx, query)
}

// Count returns the number of triples matching the query. If the store does
// not implement the Counter interface, standard Query method will be used to
// fetch all triples and the result set length is returned.
func (f *Fabric) Count(ctx context.Context, query Query) (int, error) {
	counter, ok := f.store.(Counter)
	if ok {
		return counter.Count(ctx, query)
	}

	arr, err := f.store.Query(ctx, query)
	if err != nil {
		return 0, err
	}

	return len(arr), nil
}

// Delete removes all the triples from the store matching the given query and
// returns the number of items deleted.
func (f *Fabric) Delete(ctx context.Context, query Query) (int, error) {
	return f.store.Delete(ctx, query)
}

// ReWeight performs weight updates on all triples matching the query, if the
// store implements ReWeighter interface. Otherwise, returns ErrNotSupported.
func (f *Fabric) ReWeight(ctx context.Context, query Query, delta float64, replace bool) (int, error) {
	if delta == 0 && !replace {
		// adding delta has no effect since it is zero
		return 0, errors.New("update has no effect since delta is zero and replace is false")
	}

	rew, ok := f.store.(ReWeighter)
	if !ok {
		return 0, ErrNotSupported
	}

	return rew.ReWeight(ctx, query, delta, replace)
}

// ErrNotSupported is returned when an operation is not supported.
var ErrNotSupported = errors.New("not supported")
