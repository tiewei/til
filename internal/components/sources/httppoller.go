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
	}
}

// Manifests implements translation.Translatable.
func (*HTTPPoller) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject(k8s.APISources, "HTTPPollerSource", k8s.RFC1123Name(id))

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

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
