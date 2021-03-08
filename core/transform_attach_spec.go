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

// SpecFinder can look up a suitable hcldec.Spec in a collection of specs.
type SpecFinder interface {
	FindSpec(*Specs) (hcldec.Spec, hcl.Diagnostics)
}

// AttachSpecsTransformer is a GraphTransformer that attaches a decode spec to
// all graph vertices that support it.
type AttachSpecsTransformer struct {
	Specs *Specs
}

var _ GraphTransformer = (*AttachSpecsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachSpecsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		attch, ok := v.(AttachableSpecVertex)
		if !ok {
			continue
		}

		spec, specDiags := attch.FindSpec(t.Specs)
		diags = diags.Extend(specDiags)
		attch.AttachSpec(spec)
	}

	return diags
}
