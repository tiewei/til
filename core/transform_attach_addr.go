package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/graph"
)

// AttachableAddressVertex is implemented by all types used as graph.Vertex
// that can have an event address attached.
type AttachableAddressVertex interface {
	BridgeComponentVertex

	AttachAddress(cty.Value)
}

// AttachAddressesTransformer is a GraphTransformer that attaches an event
// address to all graph vertices that support it.
//
// This "address" is not to be confused with an address in terms of block
// location in a Bridge Description File, as expressed in the "addr" package.
type AttachAddressesTransformer struct {
	Addr *Addressables
}

var _ GraphTransformer = (*AttachAddressesTransformer)(nil)

// Transform implements GraphTransformer.
func (t *AttachAddressesTransformer) Transform(g *graph.DirectedGraph) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, v := range g.Vertices() {
		attch, ok := v.(AttachableAddressVertex)
		if !ok {
			continue
		}

		addressable := t.Addr.AddressableFor(attch.Category(), attch.Type())
		if addressable == nil {
			diags = diags.Append(noAddressableDiagnostic(attch.Category(), attch.Type(), attch.SourceRange()))
			continue
		}

		addr := addressable.Address(attch.Identifier())
		if !addr.Type().Equals(destinationCty) {
			diags = diags.Append(wrongAddressTypeDiagnostic(attch.Category(), attch.Identifier(), attch.SourceRange()))
		}

		attch.AttachAddress(addr)
	}

	return diags
}

// destinationCty is a non-primitive cty.Type that represents a duck Destination.
//
// It is redefined here because k8s.DestinationCty declares optional attributes
// which make comparisons using cty.Value.Type().Equals() always return "false".
var destinationCty = cty.Object(map[string]cty.Type{
	"ref": cty.Object(map[string]cty.Type{
		"apiVersion": cty.String,
		"kind":       cty.String,
		"name":       cty.String,
	}),
})
