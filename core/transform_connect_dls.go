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

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// ConnectDeadLetterSinkTransformer is a GraphTransformer that connects
// vertices representing event senders to a global dead-letter sink.
type ConnectDeadLetterSinkTransformer struct {
	BridgeDeliveryOpts *config.Delivery
}

var _ GraphTransformer = (*ConnectDeadLetterSinkTransformer)(nil)

// Transform implements GraphTransformer.
func (t *ConnectDeadLetterSinkTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	if t.BridgeDeliveryOpts == nil {
		return nil
	}

	var diags hcl.Diagnostics

	dlsRef := t.BridgeDeliveryOpts.DeadLetterSink

	dlsAddr, parseDiags := lang.ParseBlockReference(dlsRef)
	diags = diags.Extend(parseDiags)
	if dlsAddr == nil {
		return diags
	}

	dlsV := findVertexWithAddr(g.Vertices(), dlsAddr)
	if dlsV == nil {
		return diags.Append(unknownReferenceDiagnostic(dlsAddr.Subject, dlsRef.SourceRange()))
	}

	// all referenceable vertices are messaging components, so this type
	// assertion is assumed to be safe
	dlsCmp := dlsV.(MessagingComponentVertex)

	downEdges := g.DownEdges()

	// TODO(antoineco): allow a dead-letter sink to be connected to other nodes,
	// as long as the subgraph starting at the dead-letter sink is acyclic.
	// See "Scenario 2" at triggermesh/bridgedl#137
	if _, dlsHasDownEdges := downEdges[dlsV]; dlsHasDownEdges {
		return diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Dead-letter sink is connected",
			Detail: "The selected dead-letter sink has connections to other components. This is " +
				"currently not supported.",
			Subject: dlsCmp.ComponentAddr().SourceRange.Ptr(),
		})
	}

	for _, v := range g.Vertices() {
		if _, ok := v.(ReferencerVertex); !ok {
			continue
		}

		if _, hasDownEdges := downEdges[v]; !hasDownEdges {
			continue
		}

		// NOTE(antoineco): leaving this conditional commented on purpose because it is
		// covered by the "!hasDownEdges" condition above, but we will want to uncomment it
		// as soon as the "Scenario 2" mentioned in the above TODO is enabled.
		//
		// do not connect dead-letter sink to itself
		//if v == dlsV {
		//	continue
		//}

		g.Connect(v, dlsV)
	}

	return diags
}

// findVertexWithAddr attempts to find the referenceable vertex with the given
// address in a list of vertices. Returns nil if no such vertex is found.
func findVertexWithAddr(vs graph.IndexedVertices, addr *addr.Reference) graph.Vertex {
	for _, v := range vs {
		ref, ok := v.(ReferenceableVertex)
		if !ok {
			continue
		}

		if ref.Referenceable().Addr() == addr.Subject.Addr() {
			return v
		}
	}

	return nil
}
