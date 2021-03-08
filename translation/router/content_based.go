package router

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/translation"
)

type ContentBasedTranslator struct{}

var _ translation.BlockTranslator = (*ContentBasedTranslator)(nil)

func (t *ContentBasedTranslator) Spec() hcldec.Spec {
	// NOTE(antoineco): see the following implementation to get a sense of
	// how HCL blocks map to hcldec.Specs and cty.Types:
	// https://pkg.go.dev/github.com/hashicorp/terraform@v0.14.7/configs/configschema#Block.DecoderSpec
	return &hcldec.ObjectSpec{
		"route": &hcldec.BlockSetSpec{
			TypeName: "route",
			Nested: &hcldec.ObjectSpec{
				"attributes": &hcldec.AttrSpec{
					Name:     "attributes",
					Type:     cty.Map(cty.String),
					Required: true,
				},
				"to": &hcldec.AttrSpec{
					Name:     "to",
					Type:     cty.String,
					Required: true,
				},
			},
			MinItems: 1,
		},
	}
}

func (t *ContentBasedTranslator) Manifests(interface{}) []interface{} {
	panic("not implemented")
}
