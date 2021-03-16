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

// RouterVertex is an abstract representation of a Router component within a graph.
type RouterVertex struct {
	// Address of the Router component in the Bridge description.
	Addr addr.Router
	// Router block decoded from the Bridge description.
	Router *config.Router
	// Implementation of the Router component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*RouterVertex)(nil)
	_ ReferenceableVertex      = (*RouterVertex)(nil)
	_ ReferencerVertex         = (*RouterVertex)(nil)
	_ AttachableImplVertex     = (*RouterVertex)(nil)
	_ DecodableConfigVertex    = (*RouterVertex)(nil)
	_ graph.DOTableVertex      = (*RouterVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (rtr *RouterVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryRouters,
		Type:        rtr.Router.Type,
		Identifier:  rtr.Router.Identifier,
		SourceRange: rtr.Router.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (rtr *RouterVertex) Implementation() interface{} {
	return rtr.Impl
}

// Referenceable implements ReferenceableVertex.
func (rtr *RouterVertex) Referenceable() addr.Referenceable {
	return rtr.Addr
}

// EventAddress implements ReferenceableVertex.
func (rtr *RouterVertex) EventAddress() (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := rtr.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(rtr.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), diags
	}

	config, decDiags := lang.DecodeIgnoreVars(rtr.Router.Config, rtr.Spec)
	diags = diags.Extend(decDiags)

	eventDst := cty.NullVal(k8s.DestinationCty)
	dst := addr.Address(rtr.Router.Identifier, config, eventDst)

	return dst, diags
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

// AttachImpl implements AttachableImplVertex.
func (rtr *RouterVertex) AttachImpl(impl interface{}) {
	rtr.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (rtr *RouterVertex) DecodedConfig(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return hcldec.Decode(rtr.Router.Config, rtr.Spec, ctx)
}

// AttachSpec implements DecodableConfigVertex.
func (rtr *RouterVertex) AttachSpec(s hcldec.Spec) {
	rtr.Spec = s
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
