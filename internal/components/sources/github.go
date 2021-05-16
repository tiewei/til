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

	"bridgedl/internal/sdk"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type GitHub struct{}

var (
	_ translation.Decodable    = (*GitHub)(nil)
	_ translation.Translatable = (*GitHub)(nil)
)

// Spec implements translation.Decodable.
func (*GitHub) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"event_types": &hcldec.AttrSpec{
			Name:     "event_types",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"owner_and_repository": &hcldec.AttrSpec{
			Name:     "owner_and_repository",
			Type:     cty.String,
			Required: true,
		},
		"tokens": &hcldec.AttrSpec{
			Name:     "tokens",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*GitHub) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.knative.dev/v1alpha1", "GitHubSource", k8s.RFC1123Name(id))

	eventTypes := sdk.DecodeStringSlice(config.GetAttr("event_types"))
	s.SetNestedSlice(eventTypes, "spec", "eventTypes")

	ownerAndRepository := config.GetAttr("owner_and_repository").AsString()
	s.SetNestedField(ownerAndRepository, "spec", "ownerAndRepository")

	tokens := config.GetAttr("tokens").GetAttr("name").AsString()
	accTokenSecretRef, webhookSecretRef := secrets.SecretKeyRefsGitHub(tokens)
	s.SetNestedMap(accTokenSecretRef, "spec", "accessToken", "secretKeyRef")
	s.SetNestedMap(webhookSecretRef, "spec", "secretToken", "secretKeyRef")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
