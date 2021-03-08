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

// TransformerVertex is an abstract representation of a Transformer component within a graph.
type TransformerVertex struct {
	// Address of the Transformer component in the Bridge description.
	Addr addr.Transformer
	// Transformer block decoded from the Bridge description.
	Transformer *config.Transformer
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
	// Address used as events destination.
	EventsAddr cty.Value
}

var (
	_ BridgeComponentVertex   = (*TransformerVertex)(nil)
	_ ReferenceableVertex     = (*TransformerVertex)(nil)
	_ ReferencerVertex        = (*TransformerVertex)(nil)
	_ AttachableSpecVertex    = (*TransformerVertex)(nil)
	_ AttachableAddressVertex = (*TransformerVertex)(nil)
	_ graph.DOTableVertex     = (*TransformerVertex)(nil)
)

// Category implements BridgeComponentVertex.
func (*TransformerVertex) Category() config.ComponentCategory {
	return config.CategoryTransformers
}

// Type implements BridgeComponentVertex.
func (trsf *TransformerVertex) Type() string {
	return trsf.Transformer.Type
}

// Identifer implements BridgeComponentVertex.
func (trsf *TransformerVertex) Identifier() string {
	return trsf.Transformer.Identifier
}

// SourceRange implements BridgeComponentVertex.
func (trsf *TransformerVertex) SourceRange() hcl.Range {
	return trsf.Transformer.SourceRange
}

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

// AttachSpec implements AttachableSpecVertex.
func (trsf *TransformerVertex) AttachSpec(s hcldec.Spec) {
	trsf.Spec = s
}

// AttachAddress implements AttachableAddressVertex.
func (trsf *TransformerVertex) AttachAddress(addr cty.Value) {
	trsf.EventsAddr = addr
}

// Node implements graph.DOTableVertex.
func (trsf *TransformerVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.CategoryTransformers.String(),
		Body:   trsf.Transformer.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     dotNodeColor3,
			HeaderTextColor: "white",
		},
	}
}
