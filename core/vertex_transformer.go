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
	"bridgedl/lang/k8s"
	"bridgedl/translation"
)

// TransformerVertex is an abstract representation of a Transformer component within a graph.
type TransformerVertex struct {
	// Address of the Transformer component in the Bridge description.
	Addr addr.Transformer
	// Transformer block decoded from the Bridge description.
	Transformer *config.Transformer
	// Implementation of the Transformer component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*TransformerVertex)(nil)
	_ ReferenceableVertex      = (*TransformerVertex)(nil)
	_ EventSenderVertex        = (*TransformerVertex)(nil)
	_ AttachableImplVertex     = (*TransformerVertex)(nil)
	_ DecodableConfigVertex    = (*TransformerVertex)(nil)
	_ graph.DOTableVertex      = (*TransformerVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (trsf *TransformerVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryTransformers,
		Type:        trsf.Transformer.Type,
		Identifier:  trsf.Transformer.Identifier,
		SourceRange: trsf.Transformer.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (trsf *TransformerVertex) Implementation() interface{} {
	return trsf.Impl
}

// Referenceable implements ReferenceableVertex.
func (trsf *TransformerVertex) Referenceable() addr.Referenceable {
	return trsf.Addr
}

// EventAddress implements ReferenceableVertex.
func (trsf *TransformerVertex) EventAddress(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := trsf.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(trsf.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), false, diags
	}

	cfg, cfgComplete, cfgDiags := trsf.DecodedConfig(e)
	diags = diags.Extend(cfgDiags)

	dst, dstComplete, dstDiags := trsf.EventDestination(e)
	diags = diags.Extend(dstDiags)

	evAddr := addr.Address(trsf.Transformer.Identifier, cfg, dst, e.Globals())

	if !k8s.IsDestination(evAddr) {
		diags = diags.Append(wrongAddressTypeDiagnostic(trsf.ComponentAddr()))
		evAddr = cty.UnknownVal(k8s.DestinationCty)
	}

	complete := cfgComplete && dstComplete

	return evAddr, complete, diags
}

// EventDestination implements EventSenderVertex.
func (trsf *TransformerVertex) EventDestination(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeTraversal(trsf.Transformer.To)
}

// References implements EventSenderVertex.
func (trsf *TransformerVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if trsf.Transformer == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	to, toDiags := lang.ParseBlockReference(trsf.Transformer.To)
	diags = diags.Extend(toDiags)

	if to != nil {
		refs = append(refs, to)
	}

	return refs, diags
}

// AttachImpl implements AttachableImplVertex.
func (trsf *TransformerVertex) AttachImpl(impl interface{}) {
	trsf.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (trsf *TransformerVertex) DecodedConfig(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeBlock(trsf.Transformer.Config, trsf.Spec)
}

// AttachSpec implements DecodableConfigVertex.
func (trsf *TransformerVertex) AttachSpec(s hcldec.Spec) {
	trsf.Spec = s
}

// Node implements graph.DOTableVertex.
func (trsf *TransformerVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryTransformers.String(),
		Body:   trsf.Transformer.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor3,
			HeaderTextColor: "white",
		},
	}
}
