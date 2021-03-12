package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
)

// SourceVertex is an abstract representation of a Source component within a graph.
type SourceVertex struct {
	// Source block decoded from the Bridge description.
	Source *config.Source
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ BridgeComponentVertex = (*SourceVertex)(nil)
	_ ReferencerVertex      = (*SourceVertex)(nil)
	_ AttachableSpecVertex  = (*SourceVertex)(nil)
	_ graph.DOTableVertex   = (*SourceVertex)(nil)
)

// Category implements BridgeComponentVertex.
func (*SourceVertex) Category() config.ComponentCategory {
	return config.CategorySources
}

// Type implements BridgeComponentVertex.
func (src *SourceVertex) Type() string {
	return src.Source.Type
}

// Identifer implements BridgeComponentVertex.
func (src *SourceVertex) Identifier() string {
	return src.Source.Identifier
}

// SourceRange implements BridgeComponentVertex.
func (src *SourceVertex) SourceRange() hcl.Range {
	return src.Source.SourceRange
}

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

// AttachSpec implements AttachableSpecVertex.
func (src *SourceVertex) AttachSpec(s hcldec.Spec) {
	src.Spec = s
}

// GetSpec implements AttachableSpecVertex.
func (src *SourceVertex) GetSpec() hcldec.Spec {
	return src.Spec
}

// Node implements graph.DOTableVertex.
func (src *SourceVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategorySources.String(),
		Body:   src.Source.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor4,
			HeaderTextColor: "white",
		},
	}
}
