package routers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type DataExprFilter struct{}

var (
	_ translation.Decodable    = (*DataExprFilter)(nil)
	_ translation.Translatable = (*DataExprFilter)(nil)
	_ translation.Addressable  = (*DataExprFilter)(nil)
)

// Spec implements translation.Decodable.
func (*DataExprFilter) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"condition": &hcldec.AttrSpec{
			Name:     "condition",
			Type:     cty.String,
			Required: true,
		},
		"to": &hcldec.AttrSpec{
			Name:     "to",
			Type:     k8s.DestinationCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*DataExprFilter) Manifests(id string, config, _ cty.Value) []interface{} {
	var manifests []interface{}

	f := k8s.NewObject("routing.triggermesh.io/v1alpha1", "Filter", k8s.RFC1123Name(id))

	expr := config.GetAttr("condition").AsString()
	f.SetNestedField(expr, "spec", "expression")

	sink := k8s.DecodeDestination(config.GetAttr("to").GetAttr("ref"))
	f.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, f.Unstructured())
}

// Address implements translation.Addressable.
func (*DataExprFilter) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("routing.triggermesh.io/v1alpha1", "Filter", k8s.RFC1123Name(id))
}
