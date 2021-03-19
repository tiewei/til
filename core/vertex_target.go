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

// TargetVertex is an abstract representation of a Target component within a graph.
type TargetVertex struct {
	// Address of the Target component in the Bridge description.
	Addr addr.Target
	// Target block decoded from the Bridge description.
	Target *config.Target
	// Implementation of the Target component.
	Impl interface{}
	// Spec used to decode the block configuration.
	Spec hcldec.Spec
}

var (
	_ MessagingComponentVertex = (*TargetVertex)(nil)
	_ ReferenceableVertex      = (*TargetVertex)(nil)
	_ EventSenderVertex        = (*TargetVertex)(nil)
	_ AttachableImplVertex     = (*TargetVertex)(nil)
	_ DecodableConfigVertex    = (*TargetVertex)(nil)
	_ graph.DOTableVertex      = (*TargetVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (trg *TargetVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryTargets,
		Type:        trg.Target.Type,
		Identifier:  trg.Target.Identifier,
		SourceRange: trg.Target.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (trg *TargetVertex) Implementation() interface{} {
	return trg.Impl
}

// Referenceable implements ReferenceableVertex.
func (trg *TargetVertex) Referenceable() addr.Referenceable {
	return trg.Addr
}

// EventAddress implements ReferenceableVertex.
func (trg *TargetVertex) EventAddress() (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := trg.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(trg.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), diags
	}

	config, decDiags := lang.DecodeIgnoreVars(trg.Target.Config, trg.Spec)
	diags = diags.Extend(decDiags)

	eventDst := cty.NullVal(k8s.DestinationCty)
	if trg.Target.ReplyTo != nil {
		// FIXME(antoineco): this is hacky. So far the only thing that
		// may influence the value of the event address is the presence
		// or not of a "reply_to" expression, not its actual value.
		// We should tackle this by revisiting our translation interfaces.
		eventDst = k8s.NewDestination("", "", "")
	}

	dst := addr.Address(trg.Target.Identifier, config, eventDst)

	if !k8s.IsDestination(dst) {
		diags = diags.Append(wrongAddressTypeDiagnostic(trg.ComponentAddr()))
		dst = cty.NullVal(k8s.DestinationCty)
	}

	return dst, diags
}

// EventDestination implements EventSenderVertex.
func (trg *TargetVertex) EventDestination(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	if trg.Target.ReplyTo == nil {
		return cty.NullVal(k8s.DestinationCty), nil
	}
	return trg.Target.ReplyTo.TraverseAbs(ctx)
}

// References implements EventSenderVertex.
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

// AttachImpl implements AttachableImplVertex.
func (trg *TargetVertex) AttachImpl(impl interface{}) {
	trg.Impl = impl
}

// DecodedConfig implements DecodableConfigVertex.
func (trg *TargetVertex) DecodedConfig(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	return hcldec.Decode(trg.Target.Config, trg.Spec, ctx)
}

// AttachSpec implements DecodableConfigVertex.
func (trg *TargetVertex) AttachSpec(s hcldec.Spec) {
	trg.Spec = s
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
