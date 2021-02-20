package graph

import "fmt"

// DirectedSimpleGraph is an implementation of a directed graph.
//
// Fields are kept private to ensure the consistency of the graph's state.
type DirectedGraph struct {
	// Index of all vertices in the graph.
	vertices IndexedVertices

	// Index of all edges in the graph.
	edges IndexedEdges

	// Mapping of tail vertices in the graph to their corresponding head
	// vertices. Each combination represents a down edge.
	downEdges HeadVerticesByTailVertex
	// Mapping of head vertices in the graph to their corresponding tail
	// vertices. Each combination represents an up edge.
	upEdges TailVerticesByHeadVertex
}

// NewDirectedGraph returns an initialized DirectedGraph.
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		vertices:  make(IndexedVertices),
		edges:     make(IndexedEdges),
		downEdges: make(HeadVerticesByTailVertex),
		upEdges:   make(TailVerticesByHeadVertex),
	}
}

// Vertex represents a vertex ("node") of a graph.
type Vertex interface{}

// Edge represents a directional edge ("link") of a directed graph.
type Edge struct {
	Tail Vertex
	Head Vertex
}

var _ Indexable = (*Edge)(nil)

// Key implements Indexable.
func (e *Edge) Key() interface{} {
	return fmt.Sprintf("%p->%p", e.Tail, e.Head)
}

// Add adds a Vertex to the graph.
func (g *DirectedGraph) Add(v Vertex) {
	g.vertices.Add(v)
}

// Vertices is an accessor to the graph's vertices.
func (g *DirectedGraph) Vertices() IndexedVertices {
	return g.vertices
}

// Edges is an accessor to the graph's edges.
func (g *DirectedGraph) Edges() IndexedEdges {
	return g.edges
}

// DownEdges is an accessor to the graph's down edges.
func (g *DirectedGraph) DownEdges() HeadVerticesByTailVertex {
	return g.downEdges
}

// Edges is an accessor to the graph's up edges.
func (g *DirectedGraph) UpEdges() TailVerticesByHeadVertex {
	return g.upEdges
}

// Connect connects two vertices by a directional Edge.
func (g *DirectedGraph) Connect(tail, head Vertex) {
	e := &Edge{
		Tail: tail,
		Head: head,
	}

	g.edges.Add(e)

	g.downEdges.Connect(tail, head)
	g.upEdges.Connect(head, tail)
}

// SimpleDirectedGraph is a specialization of DirectedGraph which implements a
// simple directed graph. Such graph cannot have loops (an edge connecting a
// vertex to itself).
// TODO(antoineco): implement loop detection
type SimpleDirectedGraph struct {
	DirectedGraph
}
