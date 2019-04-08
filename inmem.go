package fabric

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var _ Store = &InMemoryStore{}
var _ ReWeighter = &InMemoryStore{}
var _ Counter = &InMemoryStore{}

// InMemoryStore implements the Store interface using the golang
// map type.
type InMemoryStore struct {
	mu   *sync.RWMutex
	data map[string]Triple
}

// Count returns the number of triples in the store matching the given query.
func (mem *InMemoryStore) Count(ctx context.Context, query Query) (int, error) {
	if query.IsAny() {
		return len(mem.data), nil
	}

	triples, err := mem.Query(ctx, query)
	if err != nil {
		return 0, err
	}

	return len(triples), nil
}

// Insert stores the triple into the in-memory map.
func (mem *InMemoryStore) Insert(ctx context.Context, tri Triple) error {
	mem.ensureInit()

	mem.mu.Lock()
	defer mem.mu.Unlock()

	if mem.data == nil {
		mem.data = map[string]Triple{}
	}

	if _, ok := mem.data[tri.String()]; ok {
		return errors.New("triple already exists")
	}
	mem.data[tri.String()] = tri

	return nil
}

// Query returns all the triples matching the given query.
func (mem *InMemoryStore) Query(ctx context.Context, query Query) ([]Triple, error) {
	mem.ensureInit()

	mem.mu.RLock()
	defer mem.mu.RUnlock()

	triples := []Triple{}

	for _, tri := range mem.data {
		if query.IsAny() {
			triples = append(triples, tri)
			continue
		}

		match, err := isMatch(tri, query)
		if err != nil {
			return nil, err
		}

		if match {
			triples = append(triples, tri)
		}
	}

	return triples, nil
}

// Delete removes all the triples that match the given query.
func (mem *InMemoryStore) Delete(ctx context.Context, query Query) (int, error) {
	if mem.data == nil {
		return 0, nil
	}

	triples, err := mem.Query(ctx, query)
	if err != nil {
		return 0, err
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	for _, tri := range triples {
		delete(mem.data, tri.String())
	}

	return len(triples), nil
}

// ReWeight re-weights all the triples matching the query.
func (mem *InMemoryStore) ReWeight(ctx context.Context, query Query, delta float64, replace bool) (int, error) {
	triples, err := mem.Query(ctx, query)
	if err != nil {
		return 0, err
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	for _, tri := range triples {
		if replace {
			tri.Weight = delta
		} else {
			tri.Weight += delta
		}

		mem.data[tri.String()] = tri
	}

	return len(triples), nil
}

func (mem *InMemoryStore) ensureInit() {
	if mem.mu == nil {
		mem.mu = new(sync.RWMutex)
	}
}

func isMatch(tri Triple, query Query) (bool, error) {
	matchers := []matcher{
		matchClause(tri.Source, query.Source),
		matchClause(tri.Predicate, query.Predicate),
		matchClause(tri.Target, query.Target),
	}

	for _, matcher := range matchers {
		match, err := matcher()
		if err != nil {
			return false, err
		}

		if !match {
			return false, nil
		}
	}

	match, err := isWeightMatch(tri.Weight, query.Weight)
	if err != nil {
		return false, err
	}

	return match, nil
}

func matchClause(actual string, clause Clause) matcher {
	if clause.IsAny() {
		return func() (bool, error) {
			return true, nil
		}
	}

	switch clause.Type {
	case "=", "==", "equal":
		return func() (bool, error) {
			return clause.Value == actual, nil
		}

	case "~", "~=", "like":
		return func() (bool, error) {
			exp := strings.Replace(clause.Value, "*", ".*", -1)
			re, err := regexp.Compile(exp)
			if err != nil {
				return false, err
			}

			return re.MatchString(actual), nil
		}

	}

	return func() (bool, error) {
		return false, fmt.Errorf("clause type '%s' not supported", clause.Type)
	}
}

func isWeightMatch(actual float64, clause Clause) (bool, error) {
	if clause.IsAny() {
		return true, nil
	}

	w, err := strconv.ParseFloat(clause.Value, 64)
	if err != nil {
		return false, err
	}

	switch clause.Type {
	case "=", "==", "equal":
		return actual == w, nil

	case ">=", "gte":
		return actual >= w, nil

	case "<=", "lte":
		return actual <= w, nil

	case ">", "gt":
		return actual > w, nil

	case "<", "lt":
		return actual < w, nil
	}

	return false, nil
}

type matcher func() (bool, error)
