package router

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/translate"
)

type ContentBasedTranslator struct{}

var _ translate.BlockTranslator = (*ContentBasedTranslator)(nil)

func (t *ContentBasedTranslator) ConcreteConfig() interface{} {
	return new(ContentBased)
}

func (t *ContentBasedTranslator) K8SManifests() []interface{} {
	return nil
}

type ContentBased struct {
	Routes []*Route `hcl:"route,block"`
}

type Route struct {
	Attrs map[string]string `hcl:"attributes,attr"`
	To    hcl.Expression    `hcl:"to,attr"`
}
