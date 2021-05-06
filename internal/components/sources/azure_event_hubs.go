package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("AzureEventHubSource")
	s.SetName(k8s.RFC1123Name(id))

	hubNamespace := config.GetAttr("hub_namespace").AsString()
	_ = unstructured.SetNestedField(s.Object, hubNamespace, "spec", "hubNamespace")

	hubName := config.GetAttr("hub_name").AsString()
	_ = unstructured.SetNestedField(s.Object, hubName, "spec", "hubName")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	tenantIDSecretRef, clientIDSecretRef, clientSecrSecretRef := secrets.SecretKeyRefsAzureSP(authSecretName)
	_ = unstructured.SetNestedMap(s.Object, tenantIDSecretRef, "spec", "auth", "servicePrincipal", "tenantID", "valueFromSecret")
	_ = unstructured.SetNestedMap(s.Object, clientIDSecretRef, "spec", "auth", "servicePrincipal", "clientID", "valueFromSecret")
	_ = unstructured.SetNestedMap(s.Object, clientSecrSecretRef, "spec", "auth", "servicePrincipal", "clientSecret", "valueFromSecret")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
