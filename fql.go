package fabric

import (
	"context"
	"errors"
)

// FQL provides a simple query language on top of Fabric.
type FQL struct {
	*Fabric
}

// Exec dispatches FQL query string into appropriate Fabric query methods.
func (fql *FQL) Exec(ctx context.Context, query string) ([]Triple, error) {
	if fql == nil || fql.Fabric == nil {
		return nil, errors.New("not initialized")
	}

	// get bob* knows ?

	return nil, errors.New("not implemented")
}

var (
	// ErrNotInitialized is returned when an un-initialized FQL instance is
	// used.
	ErrNotInitialized = errors.New("not initialized")
)
