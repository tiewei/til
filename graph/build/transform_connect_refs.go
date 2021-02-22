package build

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/graph"
)

// ConnectReferencesTransformer is a GraphTransformer that connects vertices of
// a graph based on how they reference each other.
type ConnectReferencesTransformer struct{}

var _ GraphTransformer = (*ConnectReferencesTransformer)(nil)

// Transform implements GraphTransformer.
func (t *ConnectReferencesTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	// TODO(antoineco): fail on unknown references.

	vs := g.Vertices()

	rm := NewReferenceMap(vs)

	for _, v := range vs {
		for _, ref := range rm.References(v) {
			g.Connect(v, ref)
		}
	}

	return diags
}
