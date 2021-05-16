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

// ChannelVertex is an abstract representation of a Channel component within a graph.
type ChannelVertex struct {
	// Address of the Channel component in the Bridge description.
	Addr addr.Channel
	// Channel block decoded from the Bridge description.
	Channel *config.Channel
	// Implementation of the Channel component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*ChannelVertex)(nil)
	_ ReferenceableVertex      = (*ChannelVertex)(nil)
	_ EventSenderVertex        = (*ChannelVertex)(nil)
	_ AttachableImplVertex     = (*ChannelVertex)(nil)
	_ DecodableConfigVertex    = (*ChannelVertex)(nil)
	_ graph.DOTableVertex      = (*ChannelVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (ch *ChannelVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryChannels,
		Type:        ch.Channel.Type,
		Identifier:  ch.Channel.Identifier,
		SourceRange: ch.Channel.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (ch *ChannelVertex) Implementation() interface{} {
	return ch.Impl
}

// Referenceable implements ReferenceableVertex.
func (ch *ChannelVertex) Referenceable() addr.Referenceable {
	return ch.Addr
}

// EventAddress implements ReferenceableVertex.
func (ch *ChannelVertex) EventAddress(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := ch.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(ch.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), false, diags
	}

	cfg, cfgComplete, cfgDiags := ch.DecodedConfig(e)
	diags = diags.Extend(cfgDiags)

	dst, dstComplete, dstDiags := ch.EventDestination(e)
	diags = diags.Extend(dstDiags)

	evAddr := addr.Address(ch.Channel.Identifier, cfg, dst)

	if !k8s.IsDestination(evAddr) {
		diags = diags.Append(wrongAddressTypeDiagnostic(ch.ComponentAddr()))
		evAddr = cty.UnknownVal(k8s.DestinationCty)
	}

	complete := cfgComplete && dstComplete

	return evAddr, complete, diags
}

// EventDestination implements EventSenderVertex.
func (ch *ChannelVertex) EventDestination(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeTraversal(ch.Channel.To)
}

// References implements EventSenderVertex.
func (ch *ChannelVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if ch.Channel == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	to, toDiags := lang.ParseBlockReference(ch.Channel.To)
	diags = diags.Extend(toDiags)

	if to != nil {
		refs = append(refs, to)
	}

	refsInCfg, refDiags := lang.BlockReferencesInBody(ch.Channel.Config, ch.Spec)
	diags = diags.Extend(refDiags)

	refs = append(refs, refsInCfg...)

	return refs, diags
}

// AttachImpl implements AttachableImplVertex.
func (ch *ChannelVertex) AttachImpl(impl interface{}) {
	ch.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (ch *ChannelVertex) DecodedConfig(e *Evaluator) (cty.Value, bool, hcl.Diagnostics) {
	return e.DecodeBlock(ch.Channel.Config, ch.Spec)
}

// AttachSpec implements DecodableConfigVertex.
func (ch *ChannelVertex) AttachSpec(s hcldec.Spec) {
	ch.Spec = s
}

// Node implements graph.DOTableVertex.
func (ch *ChannelVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryChannels.String(),
		Body:   ch.Channel.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor1,
			HeaderTextColor: "white",
		},
	}
}
