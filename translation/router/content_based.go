package router

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/translation"
)

type ContentBasedTranslator struct{}

var _ translation.BlockTranslator = (*ContentBasedTranslator)(nil)

func (t *ContentBasedTranslator) ConcreteConfig() interface{} {
	return new(ContentBased)
}

func (t *ContentBasedTranslator) Manifests(interface{}) []interface{} {
	panic("not implemented")
}

type ContentBased struct {
	Routes []*Route `hcl:"route,block"`
}

type Route struct {
	Attrs map[string]string `hcl:"attributes,attr"`
	To    hcl.Expression    `hcl:"to,attr"`
}
