package dot_test

import (
	"testing"

	"bridgedl/graph"
	. "bridgedl/graph/dot"
)

func TestMarshal(t *testing.T) {
	g := graph.NewDirectedGraph()

	// graph vertices;
	// any Go type is valid from the "graph" package's perspective
	v1 := "some_node"
	v2 := 42
	v3 := fakeDOTableVertex{header: "some_header", body: "some_body"}

	g.Add(v1)
	g.Add(v2)
	g.Add(v3)

	g.Connect(v1, v2)
	g.Connect(v1, v3)
	g.Connect(v2, v3)

	b, err := Marshal(g)
	if err != nil {
		t.Fatal("Error marshaling graph:", err)
	}

	if string(b) != referenceDOTGraph {
		t.Error("DOT graph differs from reference:\n" + string(b))
	}
}

const (
	testAccentColor    = "goldenrod"
	testHeaderTxtColor = "floralwhite"
)

// fakeDOTableVertex is meant to be used as a graph.Vertex in tests.
type fakeDOTableVertex struct {
	header string
	body   string
}

var _ graph.DOTableVertex = (*fakeDOTableVertex)(nil)

// Node implements graph.DOTableVertex.
func (v fakeDOTableVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: v.header,
		Body:   v.body,
		Style: &graph.DOTNodeStyle{
			AccentColor:     testAccentColor,
			HeaderTextColor: testHeaderTxtColor,
		},
	}
}

const referenceDOTGraph = `strict digraph bridge {

graph [
    rankdir=LR
]

node [
    fontname="Helvetica"
    shape=plain
]

"int.42" [
    label=<
        <table border="0" color="#dcdcdc" cellborder="1" cellspacing="0" cellpadding="4"><tr><td bgcolor="#dcdcdc"><font color="#000000">int</font></td></tr><tr><td><font>42</font></td></tr></table>
    >
]
"some_header.some_body" [
    label=<
        <table border="0" color="goldenrod" cellborder="1" cellspacing="0" cellpadding="4"><tr><td bgcolor="goldenrod"><font color="floralwhite">some_header</font></td></tr><tr><td><font>some_body</font></td></tr></table>
    >
]
"string.some_node" [
    label=<
        <table border="0" color="#dcdcdc" cellborder="1" cellspacing="0" cellpadding="4"><tr><td bgcolor="#dcdcdc"><font color="#000000">string</font></td></tr><tr><td><font>some_node</font></td></tr></table>
    >
]

"int.42" -> { "some_header.some_body" }
"string.some_node" -> { "int.42" "some_header.some_body" }

}
`
