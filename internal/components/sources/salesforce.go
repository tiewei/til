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
	"bridgedl/internal/sdk/validation"
	"bridgedl/translation"
)

type Salesforce struct{}

var (
	_ translation.Decodable    = (*Salesforce)(nil)
	_ translation.Translatable = (*Salesforce)(nil)
)

// Spec implements translation.Decodable.
func (*Salesforce) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"channel": &hcldec.AttrSpec{
			Name:     "channel",
			Type:     cty.String,
			Required: true,
		},
		"replay_id": &hcldec.ValidateSpec{
			Wrapped: &hcldec.AttrSpec{
				Name:     "replay_id",
				Type:     cty.Number,
				Required: false,
			},
			Func: validation.IsInt,
		},
		"client_id": &hcldec.AttrSpec{
			Name:     "client_id",
			Type:     cty.String,
			Required: true,
		},
		"server": &hcldec.AttrSpec{
			Name:     "server",
			Type:     cty.String,
			Required: true,
		},
		"user": &hcldec.AttrSpec{
			Name:     "user",
			Type:     cty.String,
			Required: true,
		},
		"secret_key": &hcldec.AttrSpec{
			Name:     "secret_key",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Salesforce) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "SalesforceSource", name)

	channel := config.GetAttr("channel").AsString()
	s.SetNestedField(channel, "spec", "subscription", "channel")

	if v := config.GetAttr("replay_id"); !v.IsNull() {
		replayID, _ := v.AsBigFloat().Int64()
		s.SetNestedField(replayID, "spec", "subscription", "replayID")
	}

	clientID := config.GetAttr("client_id").AsString()
	s.SetNestedField(clientID, "spec", "auth", "clientID")

	server := config.GetAttr("server").AsString()
	s.SetNestedField(server, "spec", "auth", "server")

	user := config.GetAttr("user").AsString()
	s.SetNestedField(user, "spec", "auth", "user")

	secrKeySecretName := config.GetAttr("secret_key").GetAttr("name").AsString()
	secrKeySecretRef := secrets.SecretKeyRefsSalesforceOAuthJWT(secrKeySecretName)
	s.SetNestedMap(secrKeySecretRef, "spec", "auth", "certKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
