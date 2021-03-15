package routers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
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

	name := k8s.RFC1123Name(id)

	filter := &unstructured.Unstructured{}
	filter.SetAPIVersion("routing.triggermesh.io/v1alpha1")
	filter.SetKind("Filter")
	filter.SetName(name)

	expr := config.GetAttr("condition").AsString()
	_ = unstructured.SetNestedField(filter.Object, expr, "spec", "expression")

	sinkRef := config.GetAttr("to").GetAttr("ref")

	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(filter.Object, sink, "spec", "sink", "ref")

	return append(manifests, filter)
}

// Address implements translation.Addressable.
func (*DataExprFilter) Address(id string) cty.Value {
	return k8s.NewDestination("routing.triggermesh.io/v1alpha1", "Filter", k8s.RFC1123Name(id))
}
