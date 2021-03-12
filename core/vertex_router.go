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

// RouterVertex is an abstract representation of a Router component within a graph.
type RouterVertex struct {
	// Address of the Router component in the Bridge description.
	Addr addr.Router
	// Router block decoded from the Bridge description.
	Router *config.Router
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
	// Address used as events destination.
	EventsAddr cty.Value
}

var (
	_ BridgeComponentVertex   = (*RouterVertex)(nil)
	_ ReferenceableVertex     = (*RouterVertex)(nil)
	_ ReferencerVertex        = (*RouterVertex)(nil)
	_ AttachableSpecVertex    = (*RouterVertex)(nil)
	_ AttachableAddressVertex = (*RouterVertex)(nil)
	_ graph.DOTableVertex     = (*RouterVertex)(nil)
)

// Category implements BridgeComponentVertex.
func (*RouterVertex) Category() config.ComponentCategory {
	return config.CategoryRouters
}

// Type implements BridgeComponentVertex.
func (rtr *RouterVertex) Type() string {
	return rtr.Router.Type
}

// Identifer implements BridgeComponentVertex.
func (rtr *RouterVertex) Identifier() string {
	return rtr.Router.Identifier
}

// SourceRange implements BridgeComponentVertex.
func (rtr *RouterVertex) SourceRange() hcl.Range {
	return rtr.Router.SourceRange
}

// Referenceable implements ReferenceableVertex.
func (rtr *RouterVertex) Referenceable() addr.Referenceable {
	return rtr.Addr
}

// References implements ReferencerVertex.
func (rtr *RouterVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if rtr.Router == nil || rtr.Spec == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	refsInCfg, refDiags := lang.BlockReferencesInBody(rtr.Router.Config, rtr.Spec)
	diags = diags.Extend(refDiags)

	refs = append(refs, refsInCfg...)

	return refs, diags
}

// AttachSpec implements AttachableSpecVertex.
func (rtr *RouterVertex) AttachSpec(s hcldec.Spec) {
	rtr.Spec = s
}

// GetSpec implements AttachableSpecVertex.
func (rtr *RouterVertex) GetSpec() hcldec.Spec {
	return rtr.Spec
}

// AttachAddress implements AttachableAddressVertex.
func (rtr *RouterVertex) AttachAddress(addr cty.Value) {
	rtr.EventsAddr = addr
}

// GetAddress implements AttachableAddressVertex.
func (rtr *RouterVertex) GetAddress() cty.Value {
	return rtr.EventsAddr
}

// Node implements graph.DOTableVertex.
func (rtr *RouterVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryRouters.String(),
		Body:   rtr.Router.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor2,
			HeaderTextColor: "white",
		},
	}
}
