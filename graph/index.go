package graph

// Indexable is expected to be implemented by all vertices and edges of a graph
// so they can be indexed and uniquely referenced.
type Indexable interface {
	// Produces a unique key representing the caller. This key must be
	// comparable in order to be used as the key of a map.
	Key() interface{}
}

// indexKey returns the value produced by Key() if the given value implements
// Indexable, or the pointer to the original interface value otherwise.
func indexKey(v interface{}) interface{} {
	i, ok := v.(Indexable)
	if !ok {
		return v
	}
	return i.Key()
}

// IndexedVertices represents vertices of a graph indexed by key.
// The key is arbitrary but must be unique to each vertex.
type IndexedVertices map[interface{}]Vertex

// IndexedEdges represents edges of a graph indexed by key.
// The key is arbitrary but must be unique to each edge.
type IndexedEdges map[interface{}]*Edge

// HeadVerticesByTailVertex represents down edges of a graph by mapping lists
// of head vertices ("target nodes") to their associated tail vertex ("source node")
// represented by an index key unique to that vertex.
type HeadVerticesByTailVertex map[interface{}]IndexedVertices

// TailVerticesByHeadVertex represents up edges of a graph by mapping lists of
// tail vertices ("source nodes") to their associated head vertex ("target node")
// represented by an index key unique to that vertex.
type TailVerticesByHeadVertex map[interface{}]IndexedVertices

// Add indexes a Vertex.
func (i IndexedVertices) Add(v Vertex) {
	i[indexKey(v)] = v
}

// Add indexes an Edge.
func (i IndexedEdges) Add(e *Edge) {
	i[indexKey(e)] = e
}

// Connect indexes a head vertex for the given tail vertex.
func (i HeadVerticesByTailVertex) Connect(tail, head Vertex) {
	if i[indexKey(tail)] == nil {
		i[indexKey(tail)] = make(IndexedVertices)
	}
	i[indexKey(tail)].Add(head)
}

// Connect indexes a tail vertex for the given head vertex.
func (i TailVerticesByHeadVertex) Connect(head, tail Vertex) {
	if i[indexKey(head)] == nil {
		i[indexKey(head)] = make(IndexedVertices)
	}
	i[indexKey(head)].Add(tail)
}
