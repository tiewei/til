package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/graph"
)

// GraphBuilder builds a graph by applying a series of sequential
// transformation steps.
type GraphBuilder struct {
	Bridge *config.Bridge
	Impls  *componentImpls
}

// Build iterates over the transformation steps of the GraphBuilder to build a graph.
func (b *GraphBuilder) Build() (*graph.DirectedGraph, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	steps := []GraphTransformer{
		// Add all blocks as graph vertices.
		&AddComponentsTransformer{
			Bridge: b.Bridge,
		},

		// Attach component implementations.
		&AttachImplementationsTransformer{
			Impls: b.Impls,
		},

		// Attach decode specs.
		// This needs to be done before trying to evaluate references
		// between vertices, because specs allow decoding the HCL
		// configurations which contain those resolvable references.
		&AttachSpecsTransformer{},

		// Resolve references and connect vertices.
		&ConnectReferencesTransformer{},
	}

	g := graph.NewDirectedGraph()

	for _, step := range steps {
		trsfDiags := step.Transform(g)
		diags = diags.Extend(trsfDiags)
	}

	return g, diags
}

// GraphTransformer operates transformations on a graph.
type GraphTransformer interface {
	Transform(*graph.DirectedGraph) hcl.Diagnostics
}

// Color codes used for representing Bridge components on a DOT graph.
//
// Those are from the Brewer palette "Set2" (https://www.graphviz.org/doc/info/colors.html).
const (
	dotNodeColor1 = "/set26/1"
	dotNodeColor2 = "/set26/2"
	dotNodeColor3 = "/set26/3"
	dotNodeColor4 = "/set26/4"
	dotNodeColor5 = "/set26/5"
	dotNodeColor6 = "/set26/6"
)