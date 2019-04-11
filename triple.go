package fabric

import (
	"errors"
	"fmt"
	"strings"
)

// Triple represents a subject-predicate-object.
type Triple struct {
	Source    string  `json:"source" yaml:"source" db:"source"`
	Predicate string  `json:"predicate" yaml:"predicate" db:"predicate"`
	Target    string  `json:"target" yaml:"target" db:"target"`
	Weight    float64 `json:"weight" yaml:"weight" db:"weight"` // extension field
}

// Validate ensures the entity names are valid.
func (tri Triple) Validate() error {
	if strings.ContainsAny(tri.Source, forbiddenChars) || tri.Source == "" {
		return errors.New("invalid source name")
	}

	if strings.ContainsAny(tri.Predicate, forbiddenChars) || tri.Predicate == "" {
		return errors.New("invalid predicate")
	}

	if strings.ContainsAny(tri.Target, forbiddenChars) || tri.Target == "" {
		return errors.New("invalid target name")
	}

	return nil
}

func (tri Triple) String() string {
	return fmt.Sprintf("%s %s %s %f", tri.Source, tri.Predicate, tri.Target, tri.Weight)
}

var forbiddenChars = "? {}()"
