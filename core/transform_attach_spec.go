package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

	"bridgedl/graph"
)

// AttachableSpecVertex is implemented by all types used as graph.Vertex that
// can have a hcldec.Spec attached.
type AttachableSpecVertex interface {
	SpecFinder
	AttachSpec(hcldec.Spec)
}

// SpecFinder can look up a suitable hcldec.Spec in a collection of
// Translators.
type SpecFinder interface {
	FindSpec(*Translators) (hcldec.Spec, hcl.Diagnostics)
}

// AttachTranslatorsTransformer is a GraphTransformer that attaches a block
// translator to all graph vertices that support it.
type AttachTranslatorsTransformer struct {
	Translators *Translators
}

var _ GraphTransformer = (*AttachTranslatorsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachTranslatorsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		attch, ok := v.(AttachableSpecVertex)
		if !ok {
			continue
		}

		spec, translDiags := attch.FindSpec(t.Translators)
		diags = diags.Extend(translDiags)
		attch.AttachSpec(spec)
	}

	return diags
}
