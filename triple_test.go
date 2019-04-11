package fabric_test

import (
	"testing"

	"github.com/spy16/fabric"
)

func TestTriple_Validate(suite *testing.T) {
	suite.Parallel()

	cases := []struct {
		title     string
		triple    fabric.Triple
		expectErr bool
	}{
		{
			title: "InvalidSourceName",
			triple: fabric.Triple{
				Source: "?",
			},
			expectErr: true,
		},
		{
			title: "InvalidPredicateName",
			triple: fabric.Triple{
				Source:    "bob",
				Predicate: "",
			},
			expectErr: true,
		},
		{
			title: "InvalidTarget",
			triple: fabric.Triple{
				Source:    "bob",
				Predicate: "knows",
				Target:    "{",
			},
			expectErr: true,
		},
		{
			title: "Valid",
			triple: fabric.Triple{
				Source:    "bob",
				Predicate: "knows",
				Target:    "john",
			},
		},
	}

	for _, cs := range cases {
		suite.Run(cs.title, func(t *testing.T) {
			err := cs.triple.Validate()
			if err != nil {
				if !cs.expectErr {
					t.Errorf("unexpected error: %v", err)
					return
				}
				return
			}

			if cs.expectErr {
				t.Error("expecting error, got nil")
			}
		})
	}
}

func TestTriple_String(t *testing.T) {
	tri := fabric.Triple{
		Source:    "Bob",
		Predicate: "Knows",
		Target:    "John",
		Weight:    10,
	}

	expected := "Bob Knows John 10.000000"
	if tri.String() != expected {
		t.Errorf("expected string to be '%s', got '%s'", expected, tri.String())
	}
}
