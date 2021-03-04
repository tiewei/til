package core

import (
	"github.com/hashicorp/hcl/v2"

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
}

var (
	_ ReferenceableVertex = (*ChannelVertex)(nil)
	_ ReferencerVertex    = (*ChannelVertex)(nil)
	_ graph.DOTableVertex = (*ChannelVertex)(nil)
)

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

// Node implements graph.DOTableVertex.
func (ch *ChannelVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkChannel,
		Body:   ch.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor1,
			HeaderTextColor: "white",
		},
	}
}
