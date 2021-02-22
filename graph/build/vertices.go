package build

import (
	"bridgedl/config"
	"bridgedl/config/addr"
)

// ChannelVertex is an abstract representation of a Channel component within a graph.
type ChannelVertex struct {
	// Address of the Channel component in the Bridge desription.
	Addr addr.Channel
	// Channel block decoded from the Bridge desription.
	Channel *config.Channel
}

var (
	_ ReferenceableVertex = (*ChannelVertex)(nil)
	_ ReferencerVertex    = (*ChannelVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (ch *ChannelVertex) Referenceable() addr.Referenceable {
	return ch.Addr
}

// Referenceable implements ReferencerVertex.
func (ch *ChannelVertex) References() []*addr.Reference {
	if ch.Channel == nil {
		return nil
	}

	var refs []*addr.Reference

	// only return parseable references, errors should be caught in a
	// validation step
	if to, _ := addr.ParseBlockRef(ch.Channel.To); to != nil {
		refs = append(refs, to)
	}

	// TODO(antoineco): channels can have multiple outbounds depending on
	// their type (e.g. dead letter destination). We need to decode the
	// config body using a provided schema in order to be able to determine
	// all the references.
	return refs
}

// RouterVertex is an abstract representation of a Router component within a graph.
type RouterVertex struct {
	// Address of the Router component in the Bridge desription.
	Addr addr.Router
	// Router block decoded from the Bridge desription.
	Router *config.Router
}

var (
	_ ReferenceableVertex = (*RouterVertex)(nil)
	_ ReferencerVertex    = (*RouterVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (rtr *RouterVertex) Referenceable() addr.Referenceable {
	return rtr.Addr
}

// Referenceable implements ReferencerVertex.
func (rtr *RouterVertex) References() []*addr.Reference {
	if rtr.Router == nil {
		return nil
	}

	// TODO(antoineco): routers can have multiple outbounds depending on
	// their type. We need to decode the config body using a provided
	// schema in order to be able to determine all the references.
	return nil
}

// TransformerVertex is an abstract representation of a Transformer component within a graph.
type TransformerVertex struct {
	// Address of the Transformer component in the Bridge desription.
	Addr addr.Transformer
	// Transformer block decoded from the Bridge desription.
	Transformer *config.Transformer
}

var (
	_ ReferenceableVertex = (*TransformerVertex)(nil)
	_ ReferencerVertex    = (*TransformerVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (trsf *TransformerVertex) Referenceable() addr.Referenceable {
	return trsf.Addr
}

// Referenceable implements ReferencerVertex.
func (trsf *TransformerVertex) References() []*addr.Reference {
	if trsf.Transformer == nil {
		return nil
	}

	var refs []*addr.Reference

	// only return parseable references, errors should be caught in a
	// validation step
	if to, _ := addr.ParseBlockRef(trsf.Transformer.To); to != nil {
		refs = append(refs, to)
	}

	return refs
}

// SourceVertex is an abstract representation of a Source component within a graph.
type SourceVertex struct {
	// Address of the Source component in the Bridge desription.
	Addr addr.Source
	// Source block decoded from the Bridge desription.
	Source *config.Source
}

var _ ReferencerVertex = (*SourceVertex)(nil)

// Referenceable implements ReferencerVertex.
func (src *SourceVertex) References() []*addr.Reference {
	if src.Source == nil {
		return nil
	}

	var refs []*addr.Reference

	// only return parseable references, errors should be caught in a
	// validation step
	if to, _ := addr.ParseBlockRef(src.Source.To); to != nil {
		refs = append(refs, to)
	}

	return refs
}

// TargetVertex is an abstract representation of a Target component within a graph.
type TargetVertex struct {
	// Address of the Target component in the Bridge desription.
	Addr addr.Target
	// Target block decoded from the Bridge desription.
	Target *config.Target
}

var (
	_ ReferenceableVertex = (*TargetVertex)(nil)
	_ ReferencerVertex    = (*TargetVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (trg *TargetVertex) Referenceable() addr.Referenceable {
	return trg.Addr
}

// Referenceable implements ReferencerVertex.
func (trg *TargetVertex) References() []*addr.Reference {
	if trg.Target == nil {
		return nil
	}

	var refs []*addr.Reference

	// only return parseable references, errors should be caught in a
	// validation step
	if to, _ := addr.ParseBlockRef(trg.Target.ReplyTo); to != nil {
		refs = append(refs, to)
	}

	return refs
}

// FunctionVertex is an abstract representation of a Function component within a graph.
type FunctionVertex struct {
	// Address of the Function component in the Bridge desription.
	Addr addr.Function
	// Function block decoded from the Bridge desription.
	Function *config.Function
}

var (
	_ ReferenceableVertex = (*FunctionVertex)(nil)
	_ ReferencerVertex    = (*FunctionVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (fn *FunctionVertex) Referenceable() addr.Referenceable {
	return fn.Addr
}

// Referenceable implements ReferencerVertex.
func (fn *FunctionVertex) References() []*addr.Reference {
	if fn.Function == nil {
		return nil
	}

	var refs []*addr.Reference

	// only return parseable references, errors should be caught in a
	// validation step
	if to, _ := addr.ParseBlockRef(fn.Function.ReplyTo); to != nil {
		refs = append(refs, to)
	}

	return refs
}
