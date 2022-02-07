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

type AzureActivityLogs struct{}

var (
	_ translation.Decodable    = (*AzureActivityLogs)(nil)
	_ translation.Translatable = (*AzureActivityLogs)(nil)
)

// Spec implements translation.Decodable.
func (*AzureActivityLogs) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"subscription_id": &hcldec.AttrSpec{
			Name:     "subscription_id",
			Type:     cty.String,
			Required: true,
		},
		"event_hubs_namespace_id": &hcldec.AttrSpec{
			Name:     "event_hubs_namespace_id",
			Type:     cty.String,
			Required: true,
		},
		"event_hubs_instance_name": &hcldec.AttrSpec{
			Name:     "event_hubs_instance_name",
			Type:     cty.String,
			Required: false,
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
func (*AzureActivityLogs) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "AzureActivityLogsSource", name)

	subscriptionID := config.GetAttr("subscription_id").AsString()
	s.SetNestedField(subscriptionID, "spec", "subscriptionID")

	eventHubsNamespaceID := config.GetAttr("event_hubs_namespace_id").AsString()
	s.SetNestedField(eventHubsNamespaceID, "spec", "destination", "eventHubs", "namespaceID")

	if v := config.GetAttr("event_hubs_instance_name"); !v.IsNull() {
		eventHubsInstanceName := v.AsString()
		s.SetNestedField(eventHubsInstanceName, "spec", "destination", "eventHubs", "hubName")
	}

	if v := config.GetAttr("event_hubs_sas_policy"); !v.IsNull() {
		eventHubsSASPolicy := v.AsString()
		s.SetNestedField(eventHubsSASPolicy, "spec", "destination", "eventHubs", "sasPolicy")
	}

	if v := config.GetAttr("categories"); !v.IsNull() {
		categories := sdk.DecodeStringSlice(v)
		s.SetNestedSlice(categories, "spec", "categories")
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
