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

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Logz struct{}

var (
	_ translation.Decodable    = (*Logz)(nil)
	_ translation.Translatable = (*Logz)(nil)
	_ translation.Addressable  = (*Logz)(nil)
)

// Spec implements translation.Decodable.
func (*Logz) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"logs_listener_url": &hcldec.AttrSpec{
			Name:     "logs_listener_url",
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
func (*Logz) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "LogzTarget", name)

	logsListenerURL := config.GetAttr("logs_listener_url").AsString()
	t.SetNestedField(logsListenerURL, "spec", "logsListenerURL")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	shippingTokenSecretRef := secrets.SecretKeyRefsLogz(authSecretName)
	t.SetNestedMap(shippingTokenSecretRef, "spec", "shippingToken", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "LogzTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Logz) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "LogzTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
