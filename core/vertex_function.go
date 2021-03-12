package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

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
	// Address used as events destination.
	EventsAddr cty.Value
}

var (
	_ BridgeComponentVertex   = (*FunctionVertex)(nil)
	_ ReferenceableVertex     = (*FunctionVertex)(nil)
	_ ReferencerVertex        = (*FunctionVertex)(nil)
	_ AttachableAddressVertex = (*FunctionVertex)(nil)
	_ graph.DOTableVertex     = (*FunctionVertex)(nil)
)

// Category implements BridgeComponentVertex.
func (*FunctionVertex) Category() config.ComponentCategory {
	return config.CategoryFunctions
}

// Type implements BridgeComponentVertex.
func (fn *FunctionVertex) Type() string {
	// "function" types are not supported yet
	return "undefined"
}

// Identifer implements BridgeComponentVertex.
func (fn *FunctionVertex) Identifier() string {
	return fn.Function.Identifier
}

// SourceRange implements BridgeComponentVertex.
func (fn *FunctionVertex) SourceRange() hcl.Range {
	return fn.Function.SourceRange
}

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

// AttachAddress implements AttachableAddressVertex.
func (fn *FunctionVertex) AttachAddress(addr cty.Value) {
	fn.EventsAddr = addr
}

// GetAddress implements AttachableAddressVertex.
func (fn *FunctionVertex) GetAddress() cty.Value {
	return fn.EventsAddr
}

// Node implements graph.DOTableVertex.
func (fn *FunctionVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryFunctions.String(),
		Body:   fn.Function.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor6,
			HeaderTextColor: "white",
		},
	}
}
