package transformers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Function struct{}

var (
	_ translation.Decodable    = (*Function)(nil)
	_ translation.Translatable = (*Function)(nil)
	_ translation.Addressable  = (*Function)(nil)
)

// Spec implements translation.Decodable.
func (*Function) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"runtime": &hcldec.AttrSpec{
			Name:     "runtime",
			Type:     cty.String,
			Required: true,
		},
		"public": &hcldec.AttrSpec{
			Name:     "public",
			Type:     cty.Bool,
			Required: false,
		},
		"entrypoint": &hcldec.AttrSpec{
			Name:     "entrypoint",
			Type:     cty.String,
			Required: true,
		},
		"code": &hcldec.AttrSpec{
			Name:     "code",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Function) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("flow.triggermesh.io/v1alpha1")
	s.SetKind("Function")
	s.SetName(k8s.RFC1123Name(id))

	runtime := config.GetAttr("runtime").AsString()
	_ = unstructured.SetNestedField(s.Object, runtime, "spec", "runtime")
	public := config.GetAttr("public").True()
	_ = unstructured.SetNestedField(s.Object, public, "spec", "public")
	entrypoint := config.GetAttr("entrypoint").AsString()
	_ = unstructured.SetNestedField(s.Object, entrypoint, "spec", "entrypoint")
	code := config.GetAttr("code").AsString()
	_ = unstructured.SetNestedField(s.Object, code, "spec", "code")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Function) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("flow.triggermesh.io/v1alpha1", "Function", k8s.RFC1123Name(id))
}
