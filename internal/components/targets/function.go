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
		"code": &hcldec.AttrSpec{
			Name:     "code",
			Type:     cty.String,
			Required: true,
		},
		"entrypoint": &hcldec.AttrSpec{
			Name:     "entrypoint",
			Type:     cty.String,
			Required: false,
		},
		"public": &hcldec.AttrSpec{
			Name:     "public",
			Type:     cty.Bool,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Function) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	f := k8s.NewObject(k8s.APIFlow, "Function", name)

	runtime := config.GetAttr("runtime").AsString()
	f.SetNestedField(runtime, "spec", "runtime")

	code := config.GetAttr("code").AsString()
	f.SetNestedField(code, "spec", "code")

	entrypoint := "main"
	if v := config.GetAttr("entrypoint"); !v.IsNull() {
		entrypoint = v.AsString()
	}
	f.SetNestedField(entrypoint, "spec", "entrypoint")

	public := config.GetAttr("public").True()
	f.SetNestedField(public, "spec", "public")

	exts := map[string]interface{}{
		"type": "io.triggermesh.target.functions." + name,
	}
	f.SetNestedMap(exts, "spec", "ceOverrides", "extensions")

	manifests = append(manifests, f.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APIFlow, "Function", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Function) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APIFlow, "Function", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
