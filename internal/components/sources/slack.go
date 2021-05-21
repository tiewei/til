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

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Slack struct{}

var (
	_ translation.Decodable    = (*Slack)(nil)
	_ translation.Translatable = (*Slack)(nil)
)

// Spec implements translation.Decodable.
func (*Slack) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"signing_secret": &hcldec.AttrSpec{
			Name:     "signing_secret",
			Type:     k8s.ObjectReferenceCty,
			Required: false,
		},
		"app_id": &hcldec.AttrSpec{
			Name:     "app_id",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Slack) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "SlackSource", name)

	if v := config.GetAttr("signing_secret"); !v.IsNull() {
		signingSecretSecretName := v.GetAttr("name").AsString()
		signingSecret := secrets.SecretKeyRefsSlackApp(signingSecretSecretName)
		s.SetNestedMap(signingSecret, "spec", "signingSecret", "valueFromSecret")
	}

	appID := config.GetAttr("app_id")
	if !appID.IsNull() {
		s.SetNestedField(appID.AsString(), "spec", "appID")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
