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

package graph

import "testing"

func TestStronglyConnectedComponents(t *testing.T) {
	g := NewDirectedGraph()

	// Test graph from the Wikipedia page about Tarjan's algorithm:
	// https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm

	const expectSCCs = 4

	g.Add(1)
	g.Add(2)
	g.Add(3)
	g.Add(4)
	g.Add(5)
	g.Add(6)
	g.Add(7)
	g.Add(8)
	g.Connect(1, 2)
	g.Connect(2, 3)
	g.Connect(3, 1)
	g.Connect(4, 2)
	g.Connect(4, 3)
	g.Connect(4, 5)
	g.Connect(5, 4)
	g.Connect(5, 6)
	g.Connect(6, 3)
	g.Connect(6, 7)
	g.Connect(7, 6)
	g.Connect(8, 5)
	g.Connect(8, 7)
	g.Connect(8, 8)

	sccs := g.StronglyConnectedComponents()

	if gotSCCs := len(sccs); gotSCCs != expectSCCs {
		t.Fatalf("Expected %d SCCs, got %d", expectSCCs, gotSCCs)
	}

	// verify that SCCs are returned as a reverse topological sort
	// (byproduct of Tarjan's algorithm)
	assertSCCContents(t, sccs[0], 1, 2, 3)
	assertSCCContents(t, sccs[1], 6, 7)
	assertSCCContents(t, sccs[2], 4, 5)
	assertSCCContents(t, sccs[3], 8)
}

// assertSCCContents asserts that a SCC contains all given vertices.
func assertSCCContents(t *testing.T, scc []Vertex, vs ...Vertex) {
	t.Helper()

	if expectLen, gotLen := len(vs), len(scc); gotLen != expectLen {
		t.Errorf("Expected SCC to contain %d vertices, got %d (%v)", expectLen, gotLen, scc)
	}

	sccVtxIdx := make(map[Vertex]struct{}, len(scc))
	for _, v := range scc {
		sccVtxIdx[v] = struct{}{}
	}

	for _, v := range vs {
		if _, ok := sccVtxIdx[v]; !ok {
			t.Errorf("Vertex %v is missing from the given SCC (%v)", v, scc)
		}
	}
}
