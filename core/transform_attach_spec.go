package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

	"bridgedl/graph"
)

// AttachableSpecVertex is implemented by all types used as graph.Vertex that
// can have a hcldec.Spec attached.
type AttachableSpecVertex interface {
	BridgeComponentVertex

	AttachSpec(hcldec.Spec)
	GetSpec() hcldec.Spec
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

		spec := t.Specs.SpecFor(attch.Category(), attch.Type())
		if spec == nil {
			diags = diags.Append(noDecodeSpecDiagnostic(attch.Category(), attch.Type(), attch.SourceRange()))
		}

		attch.AttachSpec(spec)
	}

	return diags
}
