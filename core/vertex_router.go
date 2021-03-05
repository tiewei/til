package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"bridgedl/config"
	"bridgedl/config/addr"
	"bridgedl/graph"
	"bridgedl/lang"
	"bridgedl/translation"
)

// RouterVertex is an abstract representation of a Router component within a graph.
type RouterVertex struct {
	// Address of the Router component in the Bridge description.
	Addr addr.Router
	// Router block decoded from the Bridge description.
	Router *config.Router
	// Translator that can decode and translate a block configuration.
	Translator translation.BlockTranslator
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
func (rtr *RouterVertex) FindTranslator(tp *Translators) (translation.BlockTranslator, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	transl := tp.Routers.Translator(rtr.Router.Type)
	if transl == nil {
		diags = diags.Append(noTranslatorDiagnostic(config.BlkRouter, rtr.Router.Type, rtr.Router.SourceRange))
	}

	return transl, diags
}

// AttachBlockConfig implements AttachableTranslatorVertex.
func (rtr *RouterVertex) AttachTranslator(bt translation.BlockTranslator) {
	rtr.Translator = bt
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
