package fabric_test

import (
	"testing"

	"github.com/spy16/fabric"
)

func TestExportDOT(suite *testing.T) {
	suite.Parallel()

	suite.Run("WithNoName", func(t *testing.T) {
		expected := "digraph fabric {\n}\n"
		out := fabric.ExportDOT("", []fabric.Triple{})

		if out != expected {
			t.Errorf("expected '%s', got '%s'", expected, out)
		}
	})

	suite.Run("Normal", func(t *testing.T) {
		expected := "digraph hello {\n  \"s\" -> \"t\" [label=\"p\" weight=0.000000];\n}\n"
		out := fabric.ExportDOT("hello", []fabric.Triple{
			{
				Source:    "s",
				Target:    "t",
				Predicate: "p",
			},
		})

		if out != expected {
			t.Errorf("expected '%s', got '%s'", expected, out)
		}
	})
}
