package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Twilio struct{}

var (
	_ translation.Decodable    = (*Twilio)(nil)
	_ translation.Translatable = (*Twilio)(nil)
	_ translation.Addressable  = (*Twilio)(nil)
)

// Spec implements translation.Decodable.
func (*Twilio) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"default_phone_from": &hcldec.AttrSpec{
			Name:     "default_phone_from",
			Type:     cty.String,
			Required: false,
		},
		"default_phone_to": &hcldec.AttrSpec{
			Name:     "default_phone_to",
			Type:     cty.String,
			Required: false,
		},
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Twilio) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	t := k8s.NewObject("targets.triggermesh.io/v1alpha1", "TwilioTarget", k8s.RFC1123Name(id))

	if v := config.GetAttr("default_phone_from"); !v.IsNull() {
		defaultPhoneFrom := config.GetAttr("default_phone_from").AsString()
		t.SetNestedField(defaultPhoneFrom, "spec", "defaultPhoneFrom")
	}

	if v := config.GetAttr("default_phone_to"); !v.IsNull() {
		defaultPhoneTo := config.GetAttr("default_phone_to").AsString()
		t.SetNestedField(defaultPhoneTo, "spec", "defaultPhoneTo")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	sidSecretRef, tokenSecretRef := secrets.SecretKeyRefsTwilio(authSecretName)
	t.SetNestedMap(sidSecretRef, "spec", "sid", "secretKeyRef")
	t.SetNestedMap(tokenSecretRef, "spec", "token", "secretKeyRef")

	return append(manifests, t.Unstructured())
}

// Address implements translation.Addressable.
func (*Twilio) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "TwilioTarget", k8s.RFC1123Name(id))
}
