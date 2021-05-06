package routers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

	name := k8s.RFC1123Name(id)

	splitter := &unstructured.Unstructured{}
	splitter.SetAPIVersion("routing.triggermesh.io/v1alpha1")
	splitter.SetKind("Splitter")
	splitter.SetName(name)

	path := config.GetAttr("path").AsString()
	_ = unstructured.SetNestedField(splitter.Object, path, "spec", "path")

	typ := config.GetAttr("ce_context").GetAttr("type").AsString()
	_ = unstructured.SetNestedField(splitter.Object, typ, "spec", "ceContext", "type")

	src := config.GetAttr("ce_context").GetAttr("source").AsString()
	_ = unstructured.SetNestedField(splitter.Object, src, "spec", "ceContext", "source")

	if extsVal := config.GetAttr("ce_context").GetAttr("extensions"); !extsVal.IsNull() {
		exts := make(map[string]interface{}, extsVal.LengthInt())
		extsIter := extsVal.ElementIterator()
		for extsIter.Next() {
			attr, val := extsIter.Element()
			exts[attr.AsString()] = val.AsString()
		}
		_ = unstructured.SetNestedMap(splitter.Object, exts, "spec", "ceContext", "extensions")
	}

	sinkRef := config.GetAttr("to").GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(splitter.Object, sink, "spec", "sink", "ref")

	return append(manifests, splitter)
}

// Address implements translation.Addressable.
func (*Splitter) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("routing.triggermesh.io/v1alpha1", "Splitter", k8s.RFC1123Name(id))
}
