package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

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
}

var (
	_ ReferenceableVertex  = (*RouterVertex)(nil)
	_ ReferencerVertex     = (*RouterVertex)(nil)
	_ AttachableSpecVertex = (*RouterVertex)(nil)
	_ graph.DOTableVertex  = (*RouterVertex)(nil)
)

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

// FindSpec implements AttachableSpecVertex.
func (rtr *RouterVertex) FindSpec(s *Specs) (hcldec.Spec, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	spec := s.SpecFor(categoryRouters, rtr.Router.Type)
	if spec == nil {
		diags = diags.Append(noDecodeSpecDiagnostic(config.BlkRouter, rtr.Router.Type, rtr.Router.SourceRange))
	}

	return spec, diags
}

// AttachSpec implements AttachableSpecVertex.
func (rtr *RouterVertex) AttachSpec(s hcldec.Spec) {
	rtr.Spec = s
}

// Node implements graph.DOTableVertex.
func (rtr *RouterVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkRouter,
		Body:   rtr.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor2,
			HeaderTextColor: "white",
		},
	}
}
