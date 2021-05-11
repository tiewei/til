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

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "TwilioTarget", name)

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

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "TwilioTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Twilio) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "TwilioTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
