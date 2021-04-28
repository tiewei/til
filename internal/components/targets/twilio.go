package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
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

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("TwilioTarget")
	s.SetName(k8s.RFC1123Name(id))

	if v := config.GetAttr("default_phone_from"); !v.IsNull() {
		defaultPhoneFrom := config.GetAttr("default_phone_from").AsString()
		_ = unstructured.SetNestedField(s.Object, defaultPhoneFrom, "spec", "defaultPhoneFrom")
	}

	if v := config.GetAttr("default_phone_to"); !v.IsNull() {
		defaultPhoneTo := config.GetAttr("default_phone_to").AsString()
		_ = unstructured.SetNestedField(s.Object, defaultPhoneTo, "spec", "defaultPhoneTo")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	sidSecretRef, tokenSecretRef := secrets.SecretKeyRefsTwilio(authSecretName)
	_ = unstructured.SetNestedMap(s.Object, sidSecretRef, "spec", "sid", "secretKeyRef")
	_ = unstructured.SetNestedMap(s.Object, tokenSecretRef, "spec", "token", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Twilio) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "TwilioTarget", k8s.RFC1123Name(id))
}
