package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

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

	s := k8s.NewObject(k8s.APISources, "ZendeskSource", k8s.RFC1123Name(id))

	email := config.GetAttr("email").AsString()
	s.SetNestedField(email, "spec", "email")

	subdomain := config.GetAttr("subdomain").AsString()
	s.SetNestedField(subdomain, "spec", "subdomain")

	apiAuthSecretName := config.GetAttr("api_auth").GetAttr("name").AsString()
	tokenSecretRef := secrets.SecretKeyRefsZendesk(apiAuthSecretName)
	s.SetNestedMap(tokenSecretRef, "spec", "token", "valueFromSecret")

	webhookUsername := config.GetAttr("webhook_username").AsString()
	s.SetNestedField(webhookUsername, "spec", "webhookUsername")

	webhookPassword := config.GetAttr("webhook_password").AsString()
	s.SetNestedField(webhookPassword, "spec", "webhookPassword", "value")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
