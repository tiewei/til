package graph

// DOTableVertex can be implemented by a Vertex in order to give DOT marshalers
// a hint about its expected representation in the DOT language (Graphviz).
type DOTableVertex interface {
	Node() DOTNode
}

// DOTNode represents a "node" statement in a DOT graph.
// It wraps attributes marshalers can use to render a Vertex.
type DOTNode struct {
	Header string
	Body   string
	Style  *DOTNodeStyle
}
type DOTNodeStyle struct {
	AccentColor     string
	HeaderTextColor string
}
