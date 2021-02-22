package build

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/graph"
)

// GraphTransformer operates transformations on a graph.
type GraphTransformer interface {
	Transform(*graph.DirectedGraph) hcl.Diagnostics
}
