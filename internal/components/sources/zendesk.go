package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Zendesk struct{}

var (
	_ translation.Decodable    = (*Zendesk)(nil)
	_ translation.Translatable = (*Zendesk)(nil)
)

// Spec implements translation.Decodable.
func (*Zendesk) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"email": &hcldec.AttrSpec{
			Name:     "email",
			Type:     cty.String,
			Required: true,
		},
		"subdomain": &hcldec.AttrSpec{
			Name:     "subdomain",
			Type:     cty.String,
			Required: true,
		},
		"api_auth": &hcldec.AttrSpec{
			Name:     "api_auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
		"webhook_username": &hcldec.AttrSpec{
			Name:     "webhook_username",
			Type:     cty.String,
			Required: true,
		},
		"webhook_password": &hcldec.AttrSpec{
			Name:     "webhook_password",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Zendesk) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("ZendeskSource")
	s.SetName(k8s.RFC1123Name(id))

	email := config.GetAttr("email").AsString()
	_ = unstructured.SetNestedField(s.Object, email, "spec", "email")

	subdomain := config.GetAttr("subdomain").AsString()
	_ = unstructured.SetNestedField(s.Object, subdomain, "spec", "subdomain")

	apiAuthSecretName := config.GetAttr("api_auth").GetAttr("name").AsString()
	tokenSecretRef := secrets.SecretKeyRefsZendesk(apiAuthSecretName)
	_ = unstructured.SetNestedMap(s.Object, tokenSecretRef, "spec", "token", "valueFromSecret")

	webhookUsername := config.GetAttr("webhook_username").AsString()
	_ = unstructured.SetNestedField(s.Object, webhookUsername, "spec", "webhookUsername")

	webhookPassword := config.GetAttr("webhook_password").AsString()
	_ = unstructured.SetNestedField(s.Object, webhookPassword, "spec", "webhookPassword", "value")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
