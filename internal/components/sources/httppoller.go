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

type HTTPPoller struct{}

var (
	_ translation.Decodable    = (*HTTPPoller)(nil)
	_ translation.Translatable = (*HTTPPoller)(nil)
)

// Spec implements translation.Decodable.
func (*HTTPPoller) Spec() hcldec.Spec {
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
		"endpoint": &hcldec.AttrSpec{
			Name:     "endpoint",
			Type:     cty.String,
			Required: true,
		},
		"method": &hcldec.AttrSpec{
			Name:     "method",
			Type:     cty.String,
			Required: true,
		},
		"interval": &hcldec.AttrSpec{
			Name:     "interval",
			Type:     cty.String,
			Required: true,
		},
		"skip_verify": &hcldec.AttrSpec{
			Name:     "skip_verify",
			Type:     cty.Bool,
			Required: false,
		},
		"ca_certificate": &hcldec.AttrSpec{
			Name:     "ca_certificate",
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
			Type:     k8s.ObjectReferenceCty,
			Required: false,
		},
		"headers": &hcldec.AttrSpec{
			Name:     "headers",
			Type:     cty.Map(cty.String),
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*HTTPPoller) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "HTTPPollerSource", name)

	eventType := config.GetAttr("event_type").AsString()
	s.SetNestedField(eventType, "spec", "eventType")

	if v := config.GetAttr("event_source"); !v.IsNull() {
		eventSource := v.AsString()
		s.SetNestedField(eventSource, "spec", "eventSource")
	}

	endpoint := config.GetAttr("endpoint").AsString()
	s.SetNestedField(endpoint, "spec", "endpoint")

	method := config.GetAttr("method").AsString()
	s.SetNestedField(method, "spec", "method")

	interval := config.GetAttr("interval").AsString()
	s.SetNestedField(interval, "spec", "interval")

	if config.GetAttr("skip_verify").True() {
		s.SetNestedField(true, "spec", "skipVerify")
	}

	if v := config.GetAttr("ca_certificate"); !v.IsNull() {
		s.SetNestedField(v.AsString(), "spec", "caCertificate")
	}

	if v := config.GetAttr("basic_auth_username"); !v.IsNull() {
		s.SetNestedField(v.AsString(), "spec", "basicAuthUsername")
	}

	if v := config.GetAttr("basic_auth_password"); !v.IsNull() {
		secretName := v.GetAttr("name").AsString()
		_, secretKey := secrets.SecretKeyRefsBasicAuth(secretName)
		s.SetNestedMap(secretKey, "spec", "basicAuthPassword", "valueFromSecret")
	}

	if v := config.GetAttr("headers"); !v.IsNull() {
		headers := sdk.DecodeStringMap(v)
		s.SetNestedMap(headers, "spec", "headers")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
