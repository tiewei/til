/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"til/config/globals"
	"til/internal/sdk"
	"til/internal/sdk/k8s"
	"til/internal/sdk/secrets"
	"til/translation"
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
func (*AzureBlobStorage) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "AzureBlobStorageSource", name)

	storAccID := config.GetAttr("storage_account_id").AsString()
	s.SetNestedField(storAccID, "spec", "storageAccountID")

	eventHubID := config.GetAttr("event_hub_id").AsString()
	s.SetNestedField(eventHubID, "spec", "eventHubID")

	if v := config.GetAttr("event_types"); !v.IsNull() {
		eventTypes := sdk.DecodeStringSlice(v)
		s.SetNestedSlice(eventTypes, "spec", "eventTypes")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	tenantIDSecretRef, clientIDSecretRef, clientSecrSecretRef := secrets.SecretKeyRefsAzureSP(authSecretName)
	s.SetNestedMap(tenantIDSecretRef, "spec", "auth", "servicePrincipal", "tenantID", "valueFromSecret")
	s.SetNestedMap(clientIDSecretRef, "spec", "auth", "servicePrincipal", "clientID", "valueFromSecret")
	s.SetNestedMap(clientSecrSecretRef, "spec", "auth", "servicePrincipal", "clientSecret", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
