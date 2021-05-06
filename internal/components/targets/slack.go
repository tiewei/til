package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Slack struct{}

var (
	_ translation.Decodable    = (*Slack)(nil)
	_ translation.Translatable = (*Slack)(nil)
	_ translation.Addressable  = (*Slack)(nil)
)

// Spec implements translation.Decodable.
func (*Slack) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Slack) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("SlackTarget")
	s.SetName(k8s.RFC1123Name(id))

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	apiTokenSecretRef := secrets.SecretKeyRefsSlack(authSecretName)
	_ = unstructured.SetNestedMap(s.Object, apiTokenSecretRef, "spec", "token", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Slack) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "SlackTarget", k8s.RFC1123Name(id))
}
