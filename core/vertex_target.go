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

// TargetVertex is an abstract representation of a Target component within a graph.
type TargetVertex struct {
	// Address of the Target component in the Bridge description.
	Addr addr.Target
	// Target block decoded from the Bridge description.
	Target *config.Target
	// Implementation of the Target component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*TargetVertex)(nil)
	_ ReferenceableVertex      = (*TargetVertex)(nil)
	_ EventSenderVertex        = (*TargetVertex)(nil)
	_ AttachableImplVertex     = (*TargetVertex)(nil)
	_ DecodableConfigVertex    = (*TargetVertex)(nil)
	_ graph.DOTableVertex      = (*TargetVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (trg *TargetVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryTargets,
		Type:        trg.Target.Type,
		Identifier:  trg.Target.Identifier,
		SourceRange: trg.Target.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (trg *TargetVertex) Implementation() interface{} {
	return trg.Impl
}

// Referenceable implements ReferenceableVertex.
func (trg *TargetVertex) Referenceable() addr.Referenceable {
	return trg.Addr
}

// EventAddress implements ReferenceableVertex.
func (trg *TargetVertex) EventAddress(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := trg.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(trg.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), false, diags
	}

	cfg, cfgComplete, cfgDiags := trg.DecodedConfig(e)
	diags = diags.Extend(cfgDiags)

	dst, dstComplete, dstDiags := trg.EventDestination(e)
	diags = diags.Extend(dstDiags)

	evAddr := addr.Address(trg.Target.Identifier, cfg, dst)

	if !k8s.IsDestination(evAddr) {
		diags = diags.Append(wrongAddressTypeDiagnostic(trg.ComponentAddr()))
		evAddr = cty.UnknownVal(k8s.DestinationCty)
	}

	complete := cfgComplete && dstComplete

	return evAddr, complete, diags
}

// EventDestination implements EventSenderVertex.
func (trg *TargetVertex) EventDestination(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	if trg.Target.ReplyTo == nil {
		return cty.NullVal(k8s.DestinationCty), true, nil
	}
	return e.DecodeTraversal(trg.Target.ReplyTo)
}

// References implements EventSenderVertex.
func (trg *TargetVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if trg.Target == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	to, toDiags := lang.ParseBlockReference(trg.Target.ReplyTo)
	diags = diags.Extend(toDiags)

	if to != nil {
		refs = append(refs, to)
	}

	return refs, diags
}

// AttachImpl implements AttachableImplVertex.
func (trg *TargetVertex) AttachImpl(impl interface{}) {
	trg.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (trg *TargetVertex) DecodedConfig(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeBlock(trg.Target.Config, trg.Spec)
}

// AttachSpec implements DecodableConfigVertex.
func (trg *TargetVertex) AttachSpec(s hcldec.Spec) {
	trg.Spec = s
}

// Node implements graph.DOTableVertex.
func (trg *TargetVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryTargets.String(),
		Body:   trg.Target.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor5,
			HeaderTextColor: "white",
		},
	}
}
