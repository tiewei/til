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

	vs := g.Vertices()

	rm := NewReferenceMap(vs)

	for _, v := range vs {
		refs, refDiags := rm.References(v)
		diags = diags.Extend(refDiags)

		for _, ref := range refs {
			g.Connect(v, ref)
		}
	}

	return diags
}
