package fabric_test

import (
	"reflect"
	"testing"

	"github.com/spy16/fabric"
)

func TestQuery_IsAny(suite *testing.T) {
	suite.Parallel()

	suite.Run("WhenTrue", func(t *testing.T) {
		query := fabric.Query{}
		if !query.IsAny() {
			t.Errorf("expecting IsAny() to be true, got false")
		}
	})

	suite.Run("WhenFalse", func(t *testing.T) {
		query := fabric.Query{
			Source: fabric.Clause{
				Type:  "==",
				Value: "Bob",
			},
		}

		if query.IsAny() {
			t.Errorf("expecting IsAny() to be false, got true")
		}
	})
}

func TestQuery_Map(suite *testing.T) {
	suite.Parallel()

	suite.Run("WhenSourceIsAny", func(t *testing.T) {
		query := fabric.Query{
			Predicate: fabric.Clause{"like", "knows"},
			Target:    fabric.Clause{"like", "Bob"},
			Weight:    fabric.Clause{"gt", "10"},
		}

		src, present := query.Map()["source"]
		if present {
			t.Errorf("expecting source to be not present but found '%s'", src)
		}
	})

	suite.Run("AllPresent", func(t *testing.T) {
		query := fabric.Query{
			Source:    fabric.Clause{"~", "John"},
			Predicate: fabric.Clause{"~", "knows"},
			Target:    fabric.Clause{"!", "Bob"},
			Weight:    fabric.Clause{">", "10"},
		}

		expected := map[string]fabric.Clause{
			"source":    query.Source,
			"predicate": query.Predicate,
			"target":    query.Target,
			"weight":    query.Weight,
		}
		m := query.Map()
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("not expected: %v != %v", expected, m)
		}
	})
}
