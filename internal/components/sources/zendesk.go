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
	"til/internal/sdk/k8s"
	"til/internal/sdk/secrets"
	"til/translation"
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
func (*Zendesk) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "ZendeskSource", name)

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
