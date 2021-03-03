package build

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config/addr"
	"bridgedl/graph"
)

// ReferenceableVertex must be implemented by all types used as graph.Vertex
// that can be referenced.
type ReferenceableVertex interface {
	Referenceable() addr.Referenceable
}

// ReferencerVertex must be implemented by all types used as graph.Vertex that
// can reference addr.Referenceables.
type ReferencerVertex interface {
	References() ([]*addr.Reference, hcl.Diagnostics)
}

// ReferenceMap is a lookup map of Referenceable vertices indexed by address.
type ReferenceMap map[string]graph.Vertex

// NewReferenceMap returns a ReferenceMap initialized from the given vertices.
func NewReferenceMap(vs graph.IndexedVertices) ReferenceMap {
	rm := make(ReferenceMap)

	for _, v := range vs {
		ref, ok := v.(ReferenceableVertex)
		if !ok {
			continue
		}

		key := ref.Referenceable().Addr()
		rm[key] = v
	}

	return rm
}

// References returns all the graph vertices the given vertex refers to.
func (rm ReferenceMap) References(v graph.Vertex) ([]graph.Vertex, hcl.Diagnostics) {
	rfr, ok := v.(ReferencerVertex)
	if !ok {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var vs []graph.Vertex

	refs, refDiags := rfr.References()
	diags = diags.Extend(refDiags)

	for _, ref := range refs {
		key := ref.Subject.Addr()
		v, exists := rm[key]
		if !exists {
			diags = diags.Append(unknownReferenceDiagnostic(ref.Subject.Addr(), ref.SourceRange))
		}

		if vs != nil {
			vs = append(vs, v)
		}
	}

	return vs, diags
}
