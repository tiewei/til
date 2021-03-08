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

// TargetVertex is an abstract representation of a Target component within a graph.
type TargetVertex struct {
	// Address of the Target component in the Bridge description.
	Addr addr.Target
	// Target block decoded from the Bridge description.
	Target *config.Target
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
	// Address used as events destination.
	EventsAddr cty.Value
}

var (
	_ BridgeComponentVertex   = (*TargetVertex)(nil)
	_ ReferenceableVertex     = (*TargetVertex)(nil)
	_ ReferencerVertex        = (*TargetVertex)(nil)
	_ AttachableSpecVertex    = (*TargetVertex)(nil)
	_ AttachableAddressVertex = (*TargetVertex)(nil)
	_ graph.DOTableVertex     = (*TargetVertex)(nil)
)

// Category implements BridgeComponentVertex.
func (*TargetVertex) Category() config.ComponentCategory {
	return config.CategoryTargets
}

// Type implements BridgeComponentVertex.
func (trg *TargetVertex) Type() string {
	return trg.Target.Type
}

// Identifer implements BridgeComponentVertex.
func (trg *TargetVertex) Identifier() string {
	return trg.Target.Identifier
}

// SourceRange implements BridgeComponentVertex.
func (trg *TargetVertex) SourceRange() hcl.Range {
	return trg.Target.SourceRange
}

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

// AttachSpec implements AttachableSpecVertex.
func (trg *TargetVertex) AttachSpec(s hcldec.Spec) {
	trg.Spec = s
}

// AttachAddress implements AttachableAddressVertex.
func (trg *TargetVertex) AttachAddress(addr cty.Value) {
	trg.EventsAddr = addr
}

// Node implements graph.DOTableVertex.
func (trg *TargetVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryTargets.String(),
		Body:   trg.Target.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor5,
			HeaderTextColor: "white",
		},
	}
}
