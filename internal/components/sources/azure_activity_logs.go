package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AzureActivityLogs struct{}

var (
	_ translation.Decodable    = (*AzureActivityLogs)(nil)
	_ translation.Translatable = (*AzureActivityLogs)(nil)
)

// Spec implements translation.Decodable.
func (*AzureActivityLogs) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"event_hub_id": &hcldec.AttrSpec{
			Name:     "event_hub_id",
			Type:     cty.String,
			Required: true,
		},
		"event_hubs_sas_policy": &hcldec.AttrSpec{
			Name:     "event_hubs_sas_policy",
			Type:     cty.String,
			Required: false,
		},
		"categories": &hcldec.AttrSpec{
			Name:     "categories",
			Type:     cty.List(cty.String),
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
func (*AzureActivityLogs) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("AzureActivityLogsSource")
	s.SetName(k8s.RFC1123Name(id))

	eventHubID := config.GetAttr("event_hub_id").AsString()
	_ = unstructured.SetNestedField(s.Object, eventHubID, "spec", "eventHubID")

	if v := config.GetAttr("event_hubs_sas_policy"); !v.IsNull() {
		eventHubsSASPolicy := v.AsString()
		_ = unstructured.SetNestedField(s.Object, eventHubsSASPolicy, "spec", "eventHubsSASPolicy")
	}

	if v := config.GetAttr("categories"); !v.IsNull() {
		categoriesVals := v.AsValueSlice()
		categories := make([]string, 0, len(categoriesVals))
		for _, v := range categoriesVals {
			categories = append(categories, v.AsString())
		}
		_ = unstructured.SetNestedStringSlice(s.Object, categories, "spec", "categories")
	}

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
