package fabric

import (
	"fmt"
)

// Query represents a query to identify one or more triples.
type Query struct {
	Source    Clause `json:"source,omitempty"`
	Predicate Clause `json:"predicate,omitempty"`
	Target    Clause `json:"target,omitempty"`
	Weight    Clause `json:"weight,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// IsAny returns true if all clauses are any clauses.
func (q Query) IsAny() bool {
	return (q.Source.IsAny() && q.Predicate.IsAny() &&
		q.Target.IsAny() && q.Weight.IsAny())
}

// Map returns a map version of the query with all the any
// clauses removed.
func (q Query) Map() map[string]Clause {
	m := map[string]Clause{}
	if !q.Source.IsAny() {
		m["source"] = q.Source
	}

	if !q.Predicate.IsAny() {
		m["predicate"] = q.Predicate
	}

	if !q.Target.IsAny() {
		m["target"] = q.Target
	}

	if !q.Weight.IsAny() {
		m["weight"] = q.Weight
	}

	return m
}

// Clause represents a query clause. Zero value of this struct will
// be used as 'Any' clause which matches any value.
type Clause struct {
	// Type represents the operation that should be used. Examples include equal,
	// gt, lt etc. Supported operations are dictated by store implementations.
	Type string

	// Value that should be used as the right operand for the operation.
	Value string
}

// IsAny returns true if cl is a nil clause or both Op and Value are empty.
func (cl Clause) IsAny() bool {
	return cl.Type == "" && cl.Value == ""
}

func (cl Clause) String() string {
	return fmt.Sprintf("%s %s", cl.Type, cl.Value)
}
