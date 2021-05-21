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

package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config/globals"
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
func (*Zendesk) Manifests(id string, config, eventDst cty.Value, _ globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "ZendeskTarget", name)

	subject := config.GetAttr("subject").AsString()
	t.SetNestedField(subject, "spec", "subject")

	subdomain := config.GetAttr("subdomain").AsString()
	t.SetNestedField(subdomain, "spec", "subdomain")

	email := config.GetAttr("email").AsString()
	t.SetNestedField(email, "spec", "email")

	apiAuthSecretName := config.GetAttr("api_auth").GetAttr("name").AsString()
	tokenSecretRef := secrets.SecretKeyRefsZendesk(apiAuthSecretName)
	t.SetNestedMap(tokenSecretRef, "spec", "token", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subscriber := k8s.NewDestination(k8s.APITargets, "ZendeskTarget", name)
		subs := k8s.NewSubscription(name, name, subscriber, k8s.ReplyDest(eventDst))
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Zendesk) Address(id string, _, eventDst cty.Value, _ globals.Accessor) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "ZendeskTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
