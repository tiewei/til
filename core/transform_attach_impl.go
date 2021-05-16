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

	"bridgedl/graph"
)

// AttachableImplVertex is implemented by all types used as graph.Vertex that
// can have a component implementation attached.
type AttachableImplVertex interface {
	MessagingComponentVertex

	AttachImpl(interface{})
}

// AttachImplementationsTransformer is a GraphTransformer that attaches a
// component implementation to all graph vertices that support it.
type AttachImplementationsTransformer struct {
	Impls *componentImpls
}

var _ GraphTransformer = (*AttachImplementationsTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachImplementationsTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		attch, ok := v.(AttachableImplVertex)
		if !ok {
			continue
		}

		cmpAddr := attch.ComponentAddr()

		impl := t.Impls.ImplementationFor(cmpAddr.Category, cmpAddr.Type)
		if impl == nil {
			diags = diags.Append(noComponentImplDiagnostic(cmpAddr))
		}

		attch.AttachImpl(impl)
	}

	return diags
}
