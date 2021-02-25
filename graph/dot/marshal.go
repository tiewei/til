package dot

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"bridgedl/graph"
)

const (
	defaultAccentColor    = "#dcdcdc"
	defaultHeaderTxtColor = "#000000"
)

// Marshal serializes a graph to DOT.
func Marshal(g *graph.DirectedGraph) ([]byte, error) {
	var b bytes.Buffer

	// Static graph configuration attributes

	b.WriteString("strict digraph bridge {\n")
	b.WriteByte('\n')
	b.WriteString("graph [\n")
	b.WriteString("    rankdir=LR\n")
	b.WriteString("]\n")
	b.WriteByte('\n')
	b.WriteString("node [\n")
	b.WriteString("    fontname=\"Helvetica\"\n")
	b.WriteString("    shape=plain\n")
	b.WriteString("]\n")
	b.WriteByte('\n')

	// Vertices

	// Index of vertices already converted to node, for faster access
	// during sorting of edges.
	// The keys used in the map are also the ones used in the graph.DirectedGraph
	vertIndex := make(map[interface{}]*node)

	sortedNodes := make(nodeList, 0, len(g.Vertices()))
	for k, v := range g.Vertices() {
		n := graphVertexToNode(v)
		vertIndex[k] = n

		sortedNodes = append(sortedNodes, n)
		sort.Sort(sortedNodes)
	}

	for _, n := range sortedNodes {
		dotN, err := n.marshalDOT()
		if err != nil {
			return nil, fmt.Errorf("marshaling node to DOT: %w", err)
		}
		b.Write(dotN)
	}

	b.WriteByte('\n')

	// Edges

	sortedDownEdges := make(downEdgesList, 0, len(g.DownEdges()))
	for tailVertKey, headVerts := range g.DownEdges() {
		tailNode := vertIndex[tailVertKey]

		sortedHeadNodes := make(nodeList, 0, len(headVerts))
		for headVertKey := range headVerts {
			sortedHeadNodes = append(sortedHeadNodes, vertIndex[headVertKey])
		}
		sort.Sort(sortedHeadNodes)

		sortedDownEdges = append(sortedDownEdges, downEdges{
			tail:  tailNode,
			heads: sortedHeadNodes,
		})
	}
	sort.Sort(sortedDownEdges)

	for _, e := range sortedDownEdges {
		b.WriteString(e.tail.id() + " -> {")
		for _, h := range e.heads {
			b.WriteByte(' ')
			b.WriteString(h.id())
		}
		b.WriteString(" }\n")
	}

	b.WriteByte('\n')

	b.WriteString("}\n")

	return b.Bytes(), nil
}

// node represents a "node" statement in a DOT graph.
type node struct {
	graph.DOTNode
}

// marshalDOT serializes the node to a "node" DOT statement.
func (n *node) marshalDOT() ([]byte, error) {
	htmlShape, err := n.htmlLabel()
	if err != nil {
		return nil, fmt.Errorf("generating HTML-like label for node: %w", err)
	}

	var b bytes.Buffer

	b.WriteString(n.id() + " [\n")
	b.WriteString("    label=<\n")
	b.WriteString("        " + htmlShape + "\n")
	b.WriteString("    >\n")
	b.WriteString("]\n")

	return b.Bytes(), nil
}

// id returns the DOT "node_id" of the node.
func (n *node) id() string {
	return strconv.Quote(n.Header + "." + n.Body)
}

// htmlLabel returns the DOT "label" attribute of the node formatted as a
// HTML-string. It represents the shape of the node on the graph.
func (n *node) htmlLabel() (string, error) {
	accentColor, headerTxtColor := nodeColors(n)
	shape := newCustomShape(n.Header, n.Body, accentColor, headerTxtColor)

	// TODO(antoineco): Formatting a HTML-like label as a one-liner is not
	// very readable in the generated graph.
	// Although DOT is not necessarily meant to be read by humans, we could
	// at least try to improve this by writing each table row on its own line:
	//
	// label=<
	//     <table>
	//     <tr><td> header </td></tr>
	//     <tr><td> body </td></tr>
	//     </table>
	// >
	//
	htmlShape, err := xml.Marshal(shape)
	if err != nil {
		return "", fmt.Errorf("marshaling shape to XML: %w", err)
	}

	return string(htmlShape), nil
}

// nodeColors returns colors for the accent and header of the given node.
func nodeColors(n *node) (accent, headerTxt string) {
	accent, headerTxt = defaultAccentColor, defaultHeaderTxtColor

	if n.Style == nil {
		return
	}

	if col := n.Style.AccentColor; col != "" {
		accent = col
	}
	if col := n.Style.HeaderTextColor; col != "" {
		headerTxt = col
	}

	return
}

// graphVertexToNode converts a graph.Vertex to a marshalable DOT node.
func graphVertexToNode(v graph.Vertex) *node {
	if dv, ok := v.(graph.DOTableVertex); ok {
		return &node{DOTNode: dv.Node()}
	}

	reflectType := reflect.TypeOf(v)
	switch reflectType.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		reflectType = reflectType.Elem()
	}

	n := &node{}
	n.Header = reflectType.Name()
	n.Body = fmt.Sprint(v)

	return n
}
