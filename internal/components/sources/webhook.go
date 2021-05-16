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

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Webhook struct{}

var (
	_ translation.Decodable    = (*Webhook)(nil)
	_ translation.Translatable = (*Webhook)(nil)
)

// Spec implements translation.Decodable.
func (*Webhook) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"event_type": &hcldec.AttrSpec{
			Name:     "event_type",
			Type:     cty.String,
			Required: true,
		},
		"event_source": &hcldec.AttrSpec{
			Name:     "event_source",
			Type:     cty.String,
			Required: false,
		},
		"basic_auth_username": &hcldec.AttrSpec{
			Name:     "basic_auth_username",
			Type:     cty.String,
			Required: false,
		},
		"basic_auth_password": &hcldec.AttrSpec{
			Name:     "basic_auth_password",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Webhook) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject(k8s.APISources, "WebhookSource", k8s.RFC1123Name(id))

	eventType := config.GetAttr("event_type").AsString()
	s.SetNestedField(eventType, "spec", "eventType")

	if v := config.GetAttr("event_source"); !v.IsNull() {
		eventSource := v.AsString()
		s.SetNestedField(eventSource, "spec", "eventSource")
	}

	basicAuthUsername := config.GetAttr("basic_auth_username")
	if !basicAuthUsername.IsNull() {
		s.SetNestedField(basicAuthUsername.AsString(), "spec", "basicAuthUsername")
	}

	basicAuthPassword := config.GetAttr("basic_auth_password")
	if !basicAuthPassword.IsNull() {
		s.SetNestedField(basicAuthPassword.AsString(), "spec", "basicAuthPassword", "value")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
