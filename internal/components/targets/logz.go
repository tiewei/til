package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
	"bridgedl/translation"
)

type Logz struct{}

var (
	_ translation.Decodable    = (*Logz)(nil)
	_ translation.Translatable = (*Logz)(nil)
	_ translation.Addressable  = (*Logz)(nil)
)

// Spec implements translation.Decodable.
func (*Logz) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"logs_listener_url": &hcldec.AttrSpec{
			Name:     "logs_listener_url",
			Type:     cty.String,
			Required: true,
		},
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Logz) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("LogzTarget")
	s.SetName(k8s.RFC1123Name(id))

	logsListenerURL := config.GetAttr("logs_listener_url").AsString()
	_ = unstructured.SetNestedField(s.Object, logsListenerURL, "spec", "logsListenerURL")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	apiTokenSecretRef := secrets.SecretKeyRefsLogz(authSecretName)
	_ = unstructured.SetNestedMap(s.Object, apiTokenSecretRef, "spec", "token", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Logz) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "LogzTarget", k8s.RFC1123Name(id))
}
