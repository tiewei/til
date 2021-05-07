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
			Required: false,
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

	f := k8s.NewObject("flow.triggermesh.io/v1alpha1", "Function", k8s.RFC1123Name(id))

	runtime := config.GetAttr("runtime").AsString()
	f.SetNestedField(runtime, "spec", "runtime")

	if runtime == "js" {

		s := &unstructured.Unstructured{}
		s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
		s.SetKind("InfraTarget")
		s.SetName(k8s.RFC1123Name(id))

		code := config.GetAttr("code").AsString()
		_ = unstructured.SetNestedField(s.Object, code, "spec", "script", "code")

		name := k8s.RFC1123Name(id)
		if !eventDst.IsNull() {
			ch := k8s.NewChannel(name)
			subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APIServing, "Service", name), eventDst)
			manifests = append(manifests, ch, subs)
		}

		return append(manifests, s)
	}

	public := config.GetAttr("public").True()
	f.SetNestedField(public, "spec", "public")

	entrypoint := config.GetAttr("entrypoint").AsString()
	f.SetNestedField(entrypoint, "spec", "entrypoint")

	code := config.GetAttr("code").AsString()
	f.SetNestedField(code, "spec", "code")

	sink := k8s.DecodeDestination(eventDst)
	f.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, f.Unstructured())
}

// Address implements translation.Addressable.
func (*Function) Address(id string, config, _ cty.Value) cty.Value {

	runtime := config.GetAttr("runtime").AsString()

	if runtime == "js" {
		return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "InfraTarget", k8s.RFC1123Name(id))
	}
	return k8s.NewDestination("flow.triggermesh.io/v1alpha1", "Function", k8s.RFC1123Name(id))
}
