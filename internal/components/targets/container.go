package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Container struct{}

var (
	_ translation.Decodable    = (*Container)(nil)
	_ translation.Translatable = (*Container)(nil)
	_ translation.Addressable  = (*Container)(nil)
)

// Spec implements translation.Decodable.
func (*Container) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"image": &hcldec.AttrSpec{
			Name:     "image",
			Type:     cty.String,
			Required: true,
		},
		"public": &hcldec.AttrSpec{
			Name:     "public",
			Type:     cty.Bool,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Container) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	img := config.GetAttr("image").AsString()
	public := config.GetAttr("public").True()
	ksvc := k8s.NewKnService(name, img, public)
	manifests = append(manifests, ksvc)

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APIServing, "Service", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Container) Address(id string, _, eventDst cty.Value) cty.Value {
	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APIServing, "Service", k8s.RFC1123Name(id))
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", k8s.RFC1123Name(id))
}
