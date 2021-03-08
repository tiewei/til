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

// ChannelVertex is an abstract representation of a Channel component within a graph.
type ChannelVertex struct {
	// Address of the Channel component in the Bridge description.
	Addr addr.Channel
	// Channel block decoded from the Bridge description.
	Channel *config.Channel
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
	// Address used as events destination.
	EventsAddr cty.Value
}

var (
	_ BridgeComponentVertex   = (*ChannelVertex)(nil)
	_ ReferenceableVertex     = (*ChannelVertex)(nil)
	_ ReferencerVertex        = (*ChannelVertex)(nil)
	_ AttachableSpecVertex    = (*ChannelVertex)(nil)
	_ AttachableAddressVertex = (*ChannelVertex)(nil)
	_ graph.DOTableVertex     = (*ChannelVertex)(nil)
)

// Category implements BridgeComponentVertex.
func (*ChannelVertex) Category() config.ComponentCategory {
	return config.CategoryChannels
}

// Type implements BridgeComponentVertex.
func (ch *ChannelVertex) Type() string {
	return ch.Channel.Type
}

// Identifer implements BridgeComponentVertex.
func (ch *ChannelVertex) Identifier() string {
	return ch.Channel.Identifier
}

// SourceRange implements BridgeComponentVertex.
func (ch *ChannelVertex) SourceRange() hcl.Range {
	return ch.Channel.SourceRange
}

// Referenceable implements ReferenceableVertex.
func (ch *ChannelVertex) Referenceable() addr.Referenceable {
	return ch.Addr
}

// References implements ReferencerVertex.
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

	// TODO(antoineco): channels can have multiple outbounds depending on
	// their type (e.g. dead letter destination). We need to decode the
	// config body using a provided schema in order to be able to determine
	// all the references.
	return refs, diags
}

// AttachSpec implements AttachableSpecVertex.
func (ch *ChannelVertex) AttachSpec(s hcldec.Spec) {
	ch.Spec = s
}

// AttachAddress implements AttachableAddressVertex.
func (ch *ChannelVertex) AttachAddress(addr cty.Value) {
	ch.EventsAddr = addr
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
