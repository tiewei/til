package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// SourceVertex is an abstract representation of a Source component within a graph.
type SourceVertex struct {
	// Address of the Source component in the Bridge description.
	Addr addr.Source
	// Source block decoded from the Bridge description.
	Source *config.Source
}

var (
	_ ReferencerVertex    = (*SourceVertex)(nil)
	_ graph.DOTableVertex = (*SourceVertex)(nil)
)

// Referenceable implements ReferencerVertex.
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

// Node implements graph.DOTableVertex.
func (src *SourceVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkSource,
		Body:   src.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor4,
			HeaderTextColor: "white",
		},
	}
}
