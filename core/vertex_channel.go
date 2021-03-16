package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/k8s"
	"bridgedl/lang"
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
func (ch *ChannelVertex) EventAddress() (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := ch.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(ch.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), diags
	}

	config, decDiags := lang.DecodeIgnoreVars(ch.Channel.Config, ch.Spec)
	diags = diags.Extend(decDiags)

	eventDst := cty.NullVal(k8s.DestinationCty)
	dst := addr.Address(ch.Channel.Identifier, config, eventDst)

	return dst, diags
}

// EventDestination implements EventSenderVertex.
func (ch *ChannelVertex) EventDestination(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return ch.Channel.To.TraverseAbs(ctx)
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
func (ch *ChannelVertex) DecodedConfig(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return hcldec.Decode(ch.Channel.Config, ch.Spec, ctx)
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
