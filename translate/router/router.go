package router

import "bridgedl/translate"

// AllRouters includes all supported "router" component types.
var AllRouters = &Provider{
	components: map[string]translate.BlockTranslator{
		"content_based": new(ContentBasedTranslator),
	},
}

// Provider provides translators for components of type "router".
type Provider struct {
	components map[string]translate.BlockTranslator
}

var _ translate.TranslatorProvider = (*Provider)(nil)

func (p *Provider) Translator(componentType string) translate.BlockTranslator {
	transl, ok := p.components[componentType]
	if !ok {
		return nil
	}
	return transl
}
