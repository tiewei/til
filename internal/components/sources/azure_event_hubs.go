package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AzureEventHubs struct{}

var (
	_ translation.Decodable    = (*AzureEventHubs)(nil)
	_ translation.Translatable = (*AzureEventHubs)(nil)
)

// Spec implements translation.Decodable.
func (*AzureEventHubs) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"hub_namespace": &hcldec.AttrSpec{
			Name:     "hub_namespace",
			Type:     cty.String,
			Required: true,
		},
		"hub_name": &hcldec.AttrSpec{
			Name:     "hub_name",
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
func (*AzureEventHubs) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "AzureEventHubSource", k8s.RFC1123Name(id))

	hubNamespace := config.GetAttr("hub_namespace").AsString()
	s.SetNestedField(hubNamespace, "spec", "hubNamespace")

	hubName := config.GetAttr("hub_name").AsString()
	s.SetNestedField(hubName, "spec", "hubName")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	tenantIDSecretRef, clientIDSecretRef, clientSecrSecretRef := secrets.SecretKeyRefsAzureSP(authSecretName)
	s.SetNestedMap(tenantIDSecretRef, "spec", "auth", "servicePrincipal", "tenantID", "valueFromSecret")
	s.SetNestedMap(clientIDSecretRef, "spec", "auth", "servicePrincipal", "clientID", "valueFromSecret")
	s.SetNestedMap(clientSecrSecretRef, "spec", "auth", "servicePrincipal", "clientSecret", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
