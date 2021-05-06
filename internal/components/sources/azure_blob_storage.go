package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

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

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "AzureBlobStorageSource", id)

	storAccID := config.GetAttr("storage_account_id").AsString()
	s.SetNestedField(storAccID, "spec", "storageAccountID")

	eventHubID := config.GetAttr("event_hub_id").AsString()
	s.SetNestedField(eventHubID, "spec", "eventHubID")

	if v := config.GetAttr("event_types"); !v.IsNull() {
		eventTypesVals := v.AsValueSlice()
		eventTypes := make([]interface{}, 0, len(eventTypesVals))
		for _, v := range eventTypesVals {
			eventTypes = append(eventTypes, v.AsString())
		}
		s.SetNestedSlice(eventTypes, "spec", "eventTypes")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	tenantIDSecretRef, clientIDSecretRef, clientSecrSecretRef := secrets.SecretKeyRefsAzureSP(authSecretName)
	s.SetNestedMap(tenantIDSecretRef, "spec", "auth", "servicePrincipal", "tenantID", "valueFromSecret")
	s.SetNestedMap(clientIDSecretRef, "spec", "auth", "servicePrincipal", "clientID", "valueFromSecret")
	s.SetNestedMap(clientSecrSecretRef, "spec", "auth", "servicePrincipal", "clientSecret", "valueFromSecret")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
