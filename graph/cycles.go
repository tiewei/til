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

// StronglyConnectedComponents finds all strongly connected components[1] in
// the directed graph.
//
// This method is particularly suitable for finding cycles in a directed graph.
// The algorithm used is Tarjan's strongly connected components algorithm[2],
// which has the interesting property of identifying SCCs in an order that
// consitutes a reverse topological sort of the DAG formed by the SCCs.
//
// [1] https://en.wikipedia.org/wiki/Strongly_connected_component
// [2] https://youtu.be/wUgWX0nc4NY
func (g *DirectedGraph) StronglyConnectedComponents() [][]Vertex {
	const unvisited = -1

	var sccs [][]Vertex

	nv := len(g.vertices)
	vIDs := make(map[interface{}]int, nv)
	lowLinks := make(map[interface{}]int, nv)
	onStack := make(map[interface{}]bool, nv)
	for vIdx := range g.vertices {
		vIDs[vIdx] = unvisited
		lowLinks[vIdx] = 0
		onStack[vIdx] = false
	}

	var s stack
	nextID := 0

	var tarjan func(vIdx interface{})

	tarjan = func(vIdx interface{}) {
		s.push(vIdx)
		onStack[vIdx] = true
		vIDs[vIdx] = nextID
		lowLinks[vIdx] = nextID
		nextID++

		for adjIdx := range g.downEdges[vIdx] {
			if vIDs[adjIdx] == unvisited {
				tarjan(adjIdx)
			}
			if onStack[adjIdx] {
				lowLinks[vIdx] = min(lowLinks[vIdx], lowLinks[adjIdx])
			}
		}

		if vIDs[vIdx] == lowLinks[vIdx] {
			var scc []Vertex

			for v := s.pop(); ; v = s.pop() {
				if v == nil {
					break
				}

				onStack[v] = false
				lowLinks[v] = vIDs[vIdx]

				scc = append(scc, g.vertices[v])

				if v == vIdx {
					break
				}
			}

			sccs = append(sccs, scc)
		}
	}

	for vIdx, vID := range vIDs {
		if vID == unvisited {
			tarjan(vIdx)
		}
	}

	return sccs
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// stack is a simple stack implementation.
type stack []interface{}

// push adds an element on top of the stack.
func (s *stack) push(v interface{}) {
	*s = append(*s, v)
}

// pop removes and returns the most recently added element from the stack.
func (s *stack) pop() interface{} {
	if len(*s) == 0 {
		return nil
	}

	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}
