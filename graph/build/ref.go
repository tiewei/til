package build

import (
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
	References() []*addr.Reference
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
func (rm ReferenceMap) References(v graph.Vertex) []graph.Vertex {
	rfr, ok := v.(ReferencerVertex)
	if !ok {
		return nil
	}

	var vs []graph.Vertex

	for _, ref := range rfr.References() {
		key := ref.Subject.Addr()
		v, exists := rm[key]
		if !exists {
			continue
		}

		vs = append(vs, v)
	}

	return vs
}
