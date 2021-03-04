package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// TargetVertex is an abstract representation of a Target component within a graph.
type TargetVertex struct {
	// Address of the Target component in the Bridge description.
	Addr addr.Target
	// Target block decoded from the Bridge description.
	Target *config.Target
}

var (
	_ ReferenceableVertex = (*TargetVertex)(nil)
	_ ReferencerVertex    = (*TargetVertex)(nil)
	_ graph.DOTableVertex = (*TargetVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (trg *TargetVertex) Referenceable() addr.Referenceable {
	return trg.Addr
}

// References implements ReferencerVertex.
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

// Node implements graph.DOTableVertex.
func (trg *TargetVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkTarget,
		Body:   trg.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor5,
			HeaderTextColor: "white",
		},
	}
}
