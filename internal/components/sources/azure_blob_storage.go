package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AzureBlobStorage struct{}

var (
	_ translation.Decodable    = (*AzureBlobStorage)(nil)
	_ translation.Translatable = (*AzureBlobStorage)(nil)
)

// Spec implements translation.Decodable.
func (*AzureBlobStorage) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"storage_account_id": &hcldec.AttrSpec{
			Name:     "storage_account_id",
			Type:     cty.String,
			Required: true,
		},
		"event_hub_id": &hcldec.AttrSpec{
			Name:     "event_hub_id",
			Type:     cty.String,
			Required: true,
		},
		"event_types": &hcldec.AttrSpec{
			Name:     "event_types",
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
func (*AzureBlobStorage) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("AzureBlobStorageSource")
	s.SetName(k8s.RFC1123Name(id))

	storAccID := config.GetAttr("storage_account_id").AsString()
	_ = unstructured.SetNestedField(s.Object, storAccID, "spec", "storageAccountID")

	eventHubID := config.GetAttr("event_hub_id").AsString()
	_ = unstructured.SetNestedField(s.Object, eventHubID, "spec", "eventHubID")

	if v := config.GetAttr("event_types"); !v.IsNull() {
		eventTypesVals := v.AsValueSlice()
		eventTypes := make([]string, 0, len(eventTypesVals))
		for _, v := range eventTypesVals {
			eventTypes = append(eventTypes, v.AsString())
		}
		_ = unstructured.SetNestedStringSlice(s.Object, eventTypes, "spec", "eventTypes")
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
