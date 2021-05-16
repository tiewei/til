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
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// SourceVertex is an abstract representation of a Source component within a graph.
type SourceVertex struct {
	// Source block decoded from the Bridge description.
	Source *config.Source
	// Implementation of the Source component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*SourceVertex)(nil)
	_ EventSenderVertex        = (*SourceVertex)(nil)
	_ AttachableImplVertex     = (*SourceVertex)(nil)
	_ DecodableConfigVertex    = (*SourceVertex)(nil)
	_ graph.DOTableVertex      = (*SourceVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (src *SourceVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategorySources,
		Type:        src.Source.Type,
		Identifier:  src.Source.Identifier,
		SourceRange: src.Source.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (src *SourceVertex) Implementation() interface{} {
	return src.Impl
}

// EventDestination implements EventSenderVertex.
func (src *SourceVertex) EventDestination(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeTraversal(src.Source.To)
}

// References implements EventSenderVertex.
func (src *SourceVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if src.Source == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	to, toDiags := lang.ParseBlockReference(src.Source.To)
	diags = diags.Extend(toDiags)

	if to != nil {
		refs = append(refs, to)
	}

	return refs, diags
}

// AttachImpl implements AttachableImplVertex.
func (src *SourceVertex) AttachImpl(impl interface{}) {
	src.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (src *SourceVertex) DecodedConfig(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeBlock(src.Source.Config, src.Spec)
}

// AttachSpec implements DecodableConfigVertex.
func (src *SourceVertex) AttachSpec(s hcldec.Spec) {
	src.Spec = s
}

// Node implements graph.DOTableVertex.
func (src *SourceVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategorySources.String(),
		Body:   src.Source.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor4,
			HeaderTextColor: "white",
		},
	}
}
