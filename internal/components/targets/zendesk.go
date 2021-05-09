package targets

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
	_ translation.Addressable  = (*Zendesk)(nil)
)

// Spec implements translation.Decodable.
func (*Zendesk) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"subject": &hcldec.AttrSpec{
			Name:     "subject",
			Type:     cty.String,
			Required: true,
		},
		"subdomain": &hcldec.AttrSpec{
			Name:     "subdomain",
			Type:     cty.String,
			Required: true,
		},
		"email": &hcldec.AttrSpec{
			Name:     "email",
			Type:     cty.String,
			Required: true,
		},
		"api_auth": &hcldec.AttrSpec{
			Name:     "api_auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Zendesk) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	t := k8s.NewObject("targets.triggermesh.io/v1alpha1", "ZendeskTarget", k8s.RFC1123Name(id))

	subject := config.GetAttr("subject").AsString()
	t.SetNestedField(subject, "spec", "subject")

	subdomain := config.GetAttr("subdomain").AsString()
	t.SetNestedField(subdomain, "spec", "subdomain")

	email := config.GetAttr("email").AsString()
	t.SetNestedField(email, "spec", "email")

	apiAuthSecretName := config.GetAttr("api_auth").GetAttr("name").AsString()
	tokenSecretRef := secrets.SecretKeyRefsZendesk(apiAuthSecretName)
	t.SetNestedMap(tokenSecretRef, "spec", "token", "secretKeyRef")

	return append(manifests, t.Unstructured())
}

// Address implements translation.Addressable.
func (*Zendesk) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "ZendeskTarget", k8s.RFC1123Name(id))
}
