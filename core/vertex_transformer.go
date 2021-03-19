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

// TransformerVertex is an abstract representation of a Transformer component within a graph.
type TransformerVertex struct {
	// Address of the Transformer component in the Bridge description.
	Addr addr.Transformer
	// Transformer block decoded from the Bridge description.
	Transformer *config.Transformer
	// Implementation of the Transformer component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*TransformerVertex)(nil)
	_ ReferenceableVertex      = (*TransformerVertex)(nil)
	_ EventSenderVertex        = (*TransformerVertex)(nil)
	_ AttachableImplVertex     = (*TransformerVertex)(nil)
	_ DecodableConfigVertex    = (*TransformerVertex)(nil)
	_ graph.DOTableVertex      = (*TransformerVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (trsf *TransformerVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryTransformers,
		Type:        trsf.Transformer.Type,
		Identifier:  trsf.Transformer.Identifier,
		SourceRange: trsf.Transformer.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (trsf *TransformerVertex) Implementation() interface{} {
	return trsf.Impl
}

// Referenceable implements ReferenceableVertex.
func (trsf *TransformerVertex) Referenceable() addr.Referenceable {
	return trsf.Addr
}

// EventAddress implements ReferenceableVertex.
func (trsf *TransformerVertex) EventAddress() (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := trsf.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(trsf.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), diags
	}

	config, decDiags := lang.DecodeIgnoreVars(trsf.Transformer.Config, trsf.Spec)
	diags = diags.Extend(decDiags)

	eventDst := cty.NullVal(k8s.DestinationCty)
	dst := addr.Address(trsf.Transformer.Identifier, config, eventDst)

	if !k8s.IsDestination(dst) {
		diags = diags.Append(wrongAddressTypeDiagnostic(trsf.ComponentAddr()))
		dst = cty.NullVal(k8s.DestinationCty)
	}

	return dst, diags
}

// EventDestination implements EventSenderVertex.
func (trsf *TransformerVertex) EventDestination(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return trsf.Transformer.To.TraverseAbs(ctx)
}

// References implements EventSenderVertex.
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

// AttachImpl implements AttachableImplVertex.
func (trsf *TransformerVertex) AttachImpl(impl interface{}) {
	trsf.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (trsf *TransformerVertex) DecodedConfig(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return hcldec.Decode(trsf.Transformer.Config, trsf.Spec, ctx)
}

// AttachSpec implements DecodableConfigVertex.
func (trsf *TransformerVertex) AttachSpec(s hcldec.Spec) {
	trsf.Spec = s
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
