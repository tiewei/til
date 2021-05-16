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

type Splunk struct{}

var (
	_ translation.Decodable    = (*Splunk)(nil)
	_ translation.Translatable = (*Splunk)(nil)
	_ translation.Addressable  = (*Splunk)(nil)
)

// Spec implements translation.Decodable.
func (*Splunk) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"endpoint": &hcldec.AttrSpec{
			Name:     "endpoint",
			Type:     cty.String,
			Required: true,
		},
		"index": &hcldec.AttrSpec{
			Name:     "index",
			Type:     cty.String,
			Required: false,
		},
		"skip_tls_verify": &hcldec.AttrSpec{
			Name:     "skip_tls_verify",
			Type:     cty.Bool,
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
func (*Splunk) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "SplunkTarget", name)

	endpoint := config.GetAttr("endpoint").AsString()
	t.SetNestedField(endpoint, "spec", "endpoint")

	if v := config.GetAttr("index"); !v.IsNull() {
		index := v.AsString()
		t.SetNestedField(index, "spec", "index")
	}

	if config.GetAttr("skip_tls_verify").True() {
		t.SetNestedField(true, "spec", "skipTLSVerify")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	hecTokenSecretRef := secrets.SecretKeyRefsSplunkHEC(authSecretName)
	t.SetNestedMap(hecTokenSecretRef, "spec", "token", "valueFromSecret")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "SplunkTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Splunk) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "SplunkTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
