package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/graph"
)

// AttachableImplVertex is implemented by all types used as graph.Vertex that
// can have a component implementation attached.
type AttachableImplVertex interface {
	MessagingComponentVertex

	AttachImpl(interface{})
}

// AttachSpecsTransformer is a GraphTransformer that attaches a component
// implementation to all graph vertices that support it.
type AttachImplementationsTransformer struct {
	Impls *componentImpls
}

var _ GraphTransformer = (*AttachImplementationsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachImplementationsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		attch, ok := v.(AttachableImplVertex)
		if !ok {
			continue
		}

		cmpAddr := attch.ComponentAddr()

		impl := t.Impls.ImplementationFor(cmpAddr.Category, cmpAddr.Type)
		if impl == nil {
			diags = diags.Append(noComponentImplDiagnostic(cmpAddr))
		}

		attch.AttachImpl(impl)
	}

	return diags
}
