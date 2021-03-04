package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// TransformerVertex is an abstract representation of a Transformer component within a graph.
type TransformerVertex struct {
	// Address of the Transformer component in the Bridge description.
	Addr addr.Transformer
	// Transformer block decoded from the Bridge description.
	Transformer *config.Transformer
}

var (
	_ ReferenceableVertex = (*TransformerVertex)(nil)
	_ ReferencerVertex    = (*TransformerVertex)(nil)
	_ graph.DOTableVertex = (*TransformerVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (trsf *TransformerVertex) Referenceable() addr.Referenceable {
	return trsf.Addr
}

// References implements ReferencerVertex.
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

// Node implements graph.DOTableVertex.
func (trsf *TransformerVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkTransf,
		Body:   trsf.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor3,
			HeaderTextColor: "white",
		},
	}
}
