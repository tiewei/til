package build

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
	"bridgedl/translate"
)

// Color codes from the set26 palette at https://www.graphviz.org/doc/info/colors.html.
const (
	set26Color1 = "#66c2a5"
	set26Color2 = "#fc8d62"
	set26Color3 = "#8da0cb"
	set26Color4 = "#e78ac3"
	set26Color5 = "#a6d854"
	set26Color6 = "#ffd92f"
)

// ChannelVertex is an abstract representation of a Channel component within a graph.
type ChannelVertex struct {
	// Address of the Channel component in the Bridge description.
	Addr addr.Channel
	// Channel block decoded from the Bridge description.
	Channel *config.Channel
}

var (
	_ ReferenceableVertex = (*ChannelVertex)(nil)
	_ ReferencerVertex    = (*ChannelVertex)(nil)
	_ graph.DOTableVertex = (*ChannelVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (ch *ChannelVertex) Referenceable() addr.Referenceable {
	return ch.Addr
}

// References implements ReferencerVertex.
func (ch *ChannelVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if ch.Channel == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	to, toDiags := lang.ParseBlockReference(ch.Channel.To)
	diags = diags.Extend(toDiags)

	if to != nil {
		refs = append(refs, to)
	}

	// TODO(antoineco): channels can have multiple outbounds depending on
	// their type (e.g. dead letter destination). We need to decode the
	// config body using a provided schema in order to be able to determine
	// all the references.
	return refs, diags
}

// Node implements graph.DOTableVertex.
func (ch *ChannelVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkChannel,
		Body:   ch.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     set26Color1,
			HeaderTextColor: "white",
		},
	}
}

// RouterVertex is an abstract representation of a Router component within a graph.
type RouterVertex struct {
	// Address of the Router component in the Bridge description.
	Addr addr.Router
	// Router block decoded from the Bridge description.
	Router *config.Router
	// Translator that can decode and translate a block configuration.
	Translator translate.BlockTranslator
}

var (
	_ ReferenceableVertex        = (*RouterVertex)(nil)
	_ ReferencerVertex           = (*RouterVertex)(nil)
	_ AttachableTranslatorVertex = (*RouterVertex)(nil)
	_ graph.DOTableVertex        = (*RouterVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (rtr *RouterVertex) Referenceable() addr.Referenceable {
	return rtr.Addr
}

// References implements ReferencerVertex.
func (rtr *RouterVertex) References() ([]*addr.Reference, hcl.Diagnostics) {
	if rtr.Router == nil || rtr.Translator == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	var refs []*addr.Reference

	goConfig := rtr.Translator.ConcreteConfig()
	decodeDiags := gohcl.DecodeBody(rtr.Router.Config, nil, goConfig)
	diags = diags.Extend(decodeDiags)

	refsInCfg, refDiags := lang.BlockReferences(goConfig)
	diags = diags.Extend(refDiags)
	refs = append(refs, refsInCfg...)

	return refs, diags
}

// FindTranslator implements AttachableTranslatorVertex.
func (rtr *RouterVertex) FindTranslator(tp *translate.TranslatorProviders) (translate.BlockTranslator, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	transl := tp.Routers.Translator(rtr.Router.Type)
	if transl == nil {
		diags = diags.Append(noTranslatorDiagnostic(config.BlkRouter, rtr.Router.Type, rtr.Router.SourceRange))
	}

	return transl, diags
}

// AttachBlockConfig implements AttachableTranslatorVertex.
func (rtr *RouterVertex) AttachTranslator(tr translate.BlockTranslator) {
	rtr.Translator = tr
}

// Node implements graph.DOTableVertex.
func (rtr *RouterVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkRouter,
		Body:   rtr.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     set26Color2,
			HeaderTextColor: "white",
		},
	}
}

// TransformerVertex is an abstract representation of a Transformer component within a graph.
type TransformerVertex struct {
	// Address of the Transformer component in the Bridge description.
	Addr addr.Transformer
	// Transformer block decoded from the Bridge description.
	Transformer *config.Transformer
}

var (
	_ ReferenceableVertex = (*TransformerVertex)(nil)
	_ ReferencerVertex    = (*TransformerVertex)(nil)
	_ graph.DOTableVertex = (*TransformerVertex)(nil)
)

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

// Node implements graph.DOTableVertex.
func (trsf *TransformerVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkTransf,
		Body:   trsf.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     set26Color3,
			HeaderTextColor: "white",
		},
	}
}

// SourceVertex is an abstract representation of a Source component within a graph.
type SourceVertex struct {
	// Address of the Source component in the Bridge description.
	Addr addr.Source
	// Source block decoded from the Bridge description.
	Source *config.Source
}

var (
	_ ReferencerVertex    = (*SourceVertex)(nil)
	_ graph.DOTableVertex = (*SourceVertex)(nil)
)

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

// Node implements graph.DOTableVertex.
func (src *SourceVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkSource,
		Body:   src.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     set26Color4,
			HeaderTextColor: "white",
		},
	}
}

// TargetVertex is an abstract representation of a Target component within a graph.
type TargetVertex struct {
	// Address of the Target component in the Bridge description.
	Addr addr.Target
	// Target block decoded from the Bridge description.
	Target *config.Target
}

var (
	_ ReferenceableVertex = (*TargetVertex)(nil)
	_ ReferencerVertex    = (*TargetVertex)(nil)
	_ graph.DOTableVertex = (*TargetVertex)(nil)
)

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

// Node implements graph.DOTableVertex.
func (trg *TargetVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkTarget,
		Body:   trg.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     set26Color5,
			HeaderTextColor: "white",
		},
	}
}

// FunctionVertex is an abstract representation of a Function component within a graph.
type FunctionVertex struct {
	// Address of the Function component in the Bridge description.
	Addr addr.Function
	// Function block decoded from the Bridge description.
	Function *config.Function
}

var (
	_ ReferenceableVertex = (*FunctionVertex)(nil)
	_ ReferencerVertex    = (*FunctionVertex)(nil)
	_ graph.DOTableVertex = (*FunctionVertex)(nil)
)

// Referenceable implements ReferenceableVertex.
func (fn *FunctionVertex) Referenceable() addr.Referenceable {
	return fn.Addr
}

// References implements ReferencerVertex.
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

// Node implements graph.DOTableVertex.
func (fn *FunctionVertex) Node() graph.DOTNode {
	return graph.DOTNode{
		Header: config.BlkFunc,
		Body:   fn.Addr.Identifier,
		Style: &graph.DOTNodeStyle{
			AccentColor:     set26Color6,
			HeaderTextColor: "white",
		},
	}
}
