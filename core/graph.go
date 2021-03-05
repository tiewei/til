package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/graph"
)

// GraphBuilder builds a graph by applying a series of sequential
// transformation steps.
type GraphBuilder struct {
	Bridge      *config.Bridge
	Translators *Translators
}

// Build iterates over the transformation steps of the GraphBuilder to build a graph.
func (b *GraphBuilder) Build() (*graph.DirectedGraph, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	steps := []GraphTransformer{
		// Add all blocks as graph vertices.
		&AddComponentsTransformer{
			Bridge: b.Bridge,
		},

		// Attach block translators.
		// This needs to be done before trying to evaluate references
		// between vertices, because block translators provide access
		// to the concrete Go types containing those resolvable
		// references.
		&AttachTranslatorsTransformer{
			Translators: b.Translators,
		},

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
	dotNodeColor1 = "#66c2a5"
	dotNodeColor2 = "#fc8d62"
	dotNodeColor3 = "#8da0cb"
	dotNodeColor4 = "#e78ac3"
	dotNodeColor5 = "#a6d854"
	dotNodeColor6 = "#ffd92f"
)
