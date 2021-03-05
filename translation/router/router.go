package router

import "bridgedl/translation"

// AllRouters includes all "router" component types supported by TriggerMesh.
var AllRouters = &Provider{
	components: map[string]translation.BlockTranslator{
		"content_based": new(ContentBasedTranslator),
	},
}

// Provider provides translators for components of type "router".
type Provider struct {
	components map[string]translation.BlockTranslator
}

var _ translation.TranslatorProvider = (*Provider)(nil)

func (p *Provider) Translator(componentType string) translation.BlockTranslator {
	transl, ok := p.components[componentType]
	if !ok {
		return nil
	}
	return transl
}
