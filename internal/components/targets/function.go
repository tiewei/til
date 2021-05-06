package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

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

	f := k8s.NewObject("flow.triggermesh.io/v1alpha1", "Function", id)

	runtime := config.GetAttr("runtime").AsString()
	f.SetNestedField(runtime, "spec", "runtime")

	public := config.GetAttr("public").True()
	f.SetNestedField(public, "spec", "public")

	entrypoint := config.GetAttr("entrypoint").AsString()
	f.SetNestedField(entrypoint, "spec", "entrypoint")

	code := config.GetAttr("code").AsString()
	f.SetNestedField(code, "spec", "code")

	return append(manifests, f.Unstructured())
}

// Address implements translation.Addressable.
func (*Function) Address(id string, _, eventDst cty.Value) cty.Value {
	if eventDst.IsNull() {
		return k8s.NewDestination("flow.triggermesh.io/v1alpha1", "Function", k8s.RFC1123Name(id))
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", k8s.RFC1123Name(id))
}
