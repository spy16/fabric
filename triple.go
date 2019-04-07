package fabric

import (
	"errors"
	"fmt"
	"regexp"
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
	if !namePattern.MatchString(tri.Source) {
		return errors.New("invalid source name")
	}

	if !namePattern.MatchString(tri.Target) {
		return errors.New("invalid target name")
	}

	if !predicatePattern.MatchString(tri.Predicate) {
		return errors.New("invalid predicate")
	}

	return nil
}

func (tri Triple) String() string {
	return fmt.Sprintf("%s %s %s %f", tri.Source, tri.Predicate, tri.Target, tri.Weight)
}

var (
	namePattern      = regexp.MustCompile("^[A-Za-z0-9-\\.:]*$")
	predicatePattern = regexp.MustCompile("^[A-Za-z0-9-\\.:]*$")
)
