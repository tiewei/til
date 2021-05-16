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

package dot

import "sort"

// nodeList is a sortable list of nodes.
type nodeList []*node

var _ sort.Interface = (nodeList)(nil)

// Len implements sort.Interface.
func (n nodeList) Len() int { return len(n) }

// Swap implements sort.Interface.
func (n nodeList) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// Less implements sort.Interface.
func (n nodeList) Less(i, j int) bool { return n[i].id() < n[j].id() }

// downEdges associates a list of down (head) nodes to a tail node.
type downEdges struct {
	tail  *node
	heads []*node
}

// downEdgesList is a sortable list of downEdges associations.
type downEdgesList []downEdges

var _ sort.Interface = (downEdgesList)(nil)

// Len implements sort.Interface.
func (e downEdgesList) Len() int { return len(e) }

// Swap implements sort.Interface.
func (e downEdgesList) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// Less implements sort.Interface.
func (e downEdgesList) Less(i, j int) bool { return e[i].tail.id() < e[j].tail.id() }
