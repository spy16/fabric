package fabric

import (
	"fmt"
	"strings"
)

// ExportDOT exports the given set of triples in DOT format.
func ExportDOT(name string, triples []Triple) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "fabric"
	}

	out := fmt.Sprintf("digraph %s {\n", name)
	for _, tri := range triples {
		out += fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\" weight=%f];\n", tri.Source, tri.Target, tri.Predicate, tri.Weight)
	}
	out += "}\n"
	return out
}
