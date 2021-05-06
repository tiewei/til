package routers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Splitter struct{}

var (
	_ translation.Decodable    = (*Splitter)(nil)
	_ translation.Translatable = (*Splitter)(nil)
	_ translation.Addressable  = (*Splitter)(nil)
)

// Spec implements translation.Decodable.
func (*Splitter) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"path": &hcldec.AttrSpec{
			Name:     "path",
			Type:     cty.String,
			Required: true,
		},
		"ce_context": &hcldec.BlockSpec{
			TypeName: "ce_context",
			Nested: &hcldec.ObjectSpec{
				"type": &hcldec.AttrSpec{
					Name:     "type",
					Type:     cty.String,
					Required: true,
				},
				"source": &hcldec.AttrSpec{
					Name:     "source",
					Type:     cty.String,
					Required: true,
				},
				"extensions": &hcldec.AttrSpec{
					Name:     "extensions",
					Type:     cty.Map(cty.String),
					Required: false,
				},
			},
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
func (*Splitter) Manifests(id string, config, _ cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("routing.triggermesh.io/v1alpha1", "Splitter", id)

	path := config.GetAttr("path").AsString()
	s.SetNestedField(path, "spec", "path")

	typ := config.GetAttr("ce_context").GetAttr("type").AsString()
	s.SetNestedField(typ, "spec", "ceContext", "type")

	src := config.GetAttr("ce_context").GetAttr("source").AsString()
	s.SetNestedField(src, "spec", "ceContext", "source")

	if extsVal := config.GetAttr("ce_context").GetAttr("extensions"); !extsVal.IsNull() {
		exts := make(map[string]interface{}, extsVal.LengthInt())
		extsIter := extsVal.ElementIterator()
		for extsIter.Next() {
			attr, val := extsIter.Element()
			exts[attr.AsString()] = val.AsString()
		}
		s.SetNestedMap(exts, "spec", "ceContext", "extensions")
	}

	sinkRef := config.GetAttr("to").GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}

// Address implements translation.Addressable.
func (*Splitter) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("routing.triggermesh.io/v1alpha1", "Splitter", k8s.RFC1123Name(id))
}
