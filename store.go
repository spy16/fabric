package fabric

import "context"

// Store implementation should provide functions for managing persistence
// of triples.
type Store interface {
	// Insert should insert the given triple into the store.
	Insert(ctx context.Context, tri Triple) error

	// Query should return triples from store that match the given clauses.
	// Possible keys of the clauses map are: source, target, predicate, weight
	Query(ctx context.Context, q Query) ([]Triple, error)

	// Delete should delete triples from store that match the given clauses.
	// Clauses will follow same format as used in Query() method.
	Delete(ctx context.Context, q Query) (int, error)
}

// ReWeighter can be implemented by Store implementations to support weight
// updates. In case, this interface is not implemented, update queries will
// not be supported.
type ReWeighter interface {
	// ReWeight should update all the triples matching the query as described
	// by delta and replace flag. If replace is true, weight of all the triples
	// should be set to delta. Otherwise, delta should be added to the current
	// weights.
	ReWeight(ctx context.Context, query Query, delta float64, replace bool) (int, error)
}

// Counter can be implemented by Store implementations to support count
// operations. In case, this interface is not implemented, count queries
// will not be supported.
type Counter interface {
	// Count should return the number of triples in the store that match the
	// given query.
	Count(ctx context.Context, query Query) (int, error)
}
