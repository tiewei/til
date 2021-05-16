/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph_test

import (
	"testing"

	. "bridgedl/graph"
)

func TestDirectedGraph_Add_Connect(t *testing.T) {
	v1 := fakeVertex("vert1")
	v2 := fakeVertex("vert2")
	v3 := fakeVertex("vert3")

	g := NewDirectedGraph()

	g.Add(v1)
	g.Add(v2)

	// assert vertices

	const expectNumVerts = 2
	if numVerts := len(g.Vertices()); numVerts != expectNumVerts {
		t.Fatalf("Expected %d vertices, got %d", expectNumVerts, numVerts)
	}
	for _, vert := range []fakeVertex{v1, v2} {
		if indexed := g.Vertices()[vert.Key()]; indexed != vert {
			t.Errorf("Expected indexed vertex to equal %q, got %q", vert, indexed)
		}
	}

	g.Add(v3)

	g.Connect(v1, v2)
	g.Connect(v1, v3)

	// assert edges

	edges := g.Edges()

	const expectNumEdges = 2
	if numEdges := len(edges); numEdges != expectNumEdges {
		t.Fatalf("Expected %d edges, got %d", expectNumEdges, numEdges)
	}

	for _, edge := range []*Edge{{v1, v2}, {v1, v3}} {
		if indexedEdge := edges[edge.Key()]; *indexedEdge != *edge {
			t.Errorf("Expected indexed edge to equal %q, got %q", *edge, *indexedEdge)
		}
	}

	// assert down edges

	downEdges := g.DownEdges()

	const expectNumDownEdges = 1 // only v1 has down edges with other vertices
	if numDownEdges := len(downEdges); numDownEdges != expectNumDownEdges {
		t.Fatalf("Expected %d down edge, got %d", expectNumDownEdges, numDownEdges)
	}

	const expectNumDownVertices = 2 // v2 and v3 are down vertices of v1
	if numDownVertices := len(downEdges[v1.Key()]); numDownVertices != expectNumDownVertices {
		t.Fatalf("Expected %d down vertices, got %d", expectNumDownVertices, numDownVertices)
	}

	for _, head := range []Vertex{v2, v3} {
		indexedVert := downEdges[v1.Key()][head.(Indexable).Key()]
		if indexedVert != head {
			t.Errorf("Expected indexed down vertex to equal %q, got %q", head, indexedVert)
		}
	}

	// assert up edges

	upEdges := g.UpEdges()

	const expectNumUpEdges = 2 // both v2 and v3 have up edges with v1
	if numUpEdges := len(upEdges); numUpEdges != expectNumUpEdges {
		t.Fatalf("Expected %d up edges, got %d", expectNumDownEdges, numUpEdges)
	}

	const expectNumUpVertices = 1 // v1 is the only up vertex of both v2 and v3
	for _, head := range []Vertex{v2, v3} {
		if numUpVertices := len(upEdges[head.(Indexable).Key()]); numUpVertices != expectNumUpVertices {
			t.Fatalf("Expected %d up vertex, got %d", expectNumUpVertices, numUpVertices)
		}

		if indexedVert := upEdges[head.(Indexable).Key()][v1.Key()]; indexedVert != v1 {
			t.Errorf("Expected indexed up vertex to equal %q, got %q", v1, indexedVert)
		}
	}
}

// fakeVertex is used as a Vertex type in tests. It is backed by a string so
// that each instance can be easily printed in error messages.
type fakeVertex string

var _ Indexable = (*fakeVertex)(nil)

func (v fakeVertex) Key() interface{} {
	return "key{" + v + "}"
}
