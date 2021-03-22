package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/k8s"
	"bridgedl/lang"
	"bridgedl/translation"
)

// FunctionVertex is an abstract representation of a Function component within a graph.
type FunctionVertex struct {
	// Address of the Function component in the Bridge description.
	Addr addr.Function
	// Function block decoded from the Bridge description.
	Function *config.Function
	// Implementation of the Function component.
	Impl interface{}
}

var (
	_ MessagingComponentVertex = (*FunctionVertex)(nil)
	_ ReferenceableVertex      = (*FunctionVertex)(nil)
	_ EventSenderVertex        = (*FunctionVertex)(nil)
	_ AttachableImplVertex     = (*FunctionVertex)(nil)
	_ graph.DOTableVertex      = (*FunctionVertex)(nil)
)

// ComponentAddr implements MessagingComponentVertex.
func (fn *FunctionVertex) ComponentAddr() addr.MessagingComponent {
	return addr.MessagingComponent{
		Category:    config.CategoryFunctions,
		Identifier:  fn.Function.Identifier,
		SourceRange: fn.Function.SourceRange,
	}
}

// Implementation implements MessagingComponentVertex.
func (fn *FunctionVertex) Implementation() interface{} {
	return fn.Impl
}

// Referenceable implements ReferenceableVertex.
func (fn *FunctionVertex) Referenceable() addr.Referenceable {
	return fn.Addr
}

// EventAddress implements ReferenceableVertex.
func (fn *FunctionVertex) EventAddress() (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addr, ok := fn.Impl.(translation.Addressable)
	if !ok {
		diags = diags.Append(noAddressableDiagnostic(fn.ComponentAddr()))
		return cty.NullVal(k8s.DestinationCty), diags
	}

	// functions of different types aren't supported yet, so there is no
	// body to decode
	config := cty.NullVal(cty.DynamicPseudoType)

	eventDst := cty.NullVal(k8s.DestinationCty)
	if fn.Function.ReplyTo != nil {
		// FIXME(antoineco): this is hacky. So far the only thing that
		// may influence the value of the event address is the presence
		// or not of a "reply_to" expression, not its actual value.
		// We should tackle this by revisiting our translation interfaces.
		eventDst = k8s.NewDestination("", "", "")
	}

	dst := addr.Address(fn.Function.Identifier, config, eventDst)

	if !k8s.IsDestination(dst) {
		diags = diags.Append(wrongAddressTypeDiagnostic(fn.ComponentAddr()))
		dst = cty.NullVal(k8s.DestinationCty)
	}

	return dst, diags
}

// EventDestination implements EventSenderVertex.
func (fn *FunctionVertex) EventDestination(ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	if fn.Function.ReplyTo == nil {
		return cty.NullVal(k8s.DestinationCty), nil
	}
	return fn.Function.ReplyTo.TraverseAbs(ctx)
}

// References implements EventSenderVertex.
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

// AttachImpl implements AttachableImplVertex.
func (fn *FunctionVertex) AttachImpl(impl interface{}) {
	fn.Impl = impl
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