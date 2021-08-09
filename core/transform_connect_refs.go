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

package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"til/config/addr"
	"til/graph"
)

// AddressableVertex is implemented by all types used as graph.Vertex that can
// expose an address for receiving events.
type AddressableVertex interface {
	// Address for receiving events.
	// The returned boolean value indicates whether all expressions from
	// the component's configuration could be decoded without injecting
	// placeholders into the evaluation context.
	EventAddress(*Evaluator) (cty.Value, bool, hcl.Diagnostics)

	// This interface is embedded for convenience. During a graph
	// evaluation, knowledge about the category of the component
	// represented by the vertex is required to be able to populate the
	// Evaluator (in addition to the event address).
	MessagingComponentVertex
}

// ReferenceableVertex is implemented by all types used as graph.Vertex that
// can be referenced by other vertices.
type ReferenceableVertex interface {
	Referenceable() addr.Referenceable

	// In the current version of the TriggerMesh Integration Language, a
	// Referenceable vertex must also expose an address and accept events.
	// This may change in the future.
	AddressableVertex
}

// ReferencerVertex is implemented by all types used as graph.Vertex that can
// reference other vertices.
type ReferencerVertex interface {
	References() ([]*addr.Reference, hcl.Diagnostics)
}

// EventSenderVertex is implemented by all types used as graph.Vertex that may
// have a main event destination configured ("to" top-level HCL attribute).
type EventSenderVertex interface {
	// Event destination.
	// The returned boolean value indicates whether all expressions from
	// the component's configuration could be decoded without injecting
	// placeholders into the evaluation context.
	EventDestination(*Evaluator) (cty.Value, bool, hcl.Diagnostics)

	// If a component can send events, it can also have at least one
	// reference to other components.
	ReferencerVertex
}

// ConnectReferencesTransformer is a GraphTransformer that connects vertices of
// a graph based on how they reference each other.
type ConnectReferencesTransformer struct{}

var _ GraphTransformer = (*ConnectReferencesTransformer)(nil)

// Transform implements GraphTransformer.
func (t *ConnectReferencesTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	vs := g.Vertices()

	rm := NewReferenceMap(vs)

	for _, v := range vs {
		refs, refDiags := rm.References(v)
		diags = diags.Extend(refDiags)

		for _, ref := range refs {
			g.Connect(v, ref)
		}
	}

	return diags
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
func (rm ReferenceMap) References(v graph.Vertex) ([]graph.Vertex, hcl.Diagnostics) {
	rfr, ok := v.(ReferencerVertex)
	if !ok {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var vs []graph.Vertex

	refs, refDiags := rfr.References()
	diags = diags.Extend(refDiags)

	for _, ref := range refs {
		key := ref.Subject.Addr()
		v, exists := rm[key]
		if !exists {
			diags = diags.Append(unknownReferenceDiagnostic(ref.Subject, ref.SourceRange))
		}

		if v != nil {
			vs = append(vs, v)
		}
	}

	return vs, diags
}
