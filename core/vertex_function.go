package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// FunctionVertex is an abstract representation of a Function component within a graph.
type FunctionVertex struct {
	// Address of the Function component in the Bridge description.
	Addr addr.Function
	// Function block decoded from the Bridge description.
	Function *config.Function
}

var (
	_ ReferenceableVertex = (*FunctionVertex)(nil)
	_ ReferencerVertex    = (*FunctionVertex)(nil)
	_ graph.DOTableVertex = (*FunctionVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (fn *FunctionVertex) Referenceable() addr.Referenceable {
	return fn.Addr
}

// References implements ReferencerVertex.
func (fn *FunctionVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if fn.Function == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	to, toDiags := lang.ParseBlockReference(fn.Function.ReplyTo)
	diags = diags.Extend(toDiags)

	if to != nil {
		refs = append(refs, to)
	}

	return refs, diags
}

// Node implements graph.DOTableVertex.
func (fn *FunctionVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkFunc,
		Body:   fn.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor6,
			HeaderTextColor: "white",
		},
	}
}
