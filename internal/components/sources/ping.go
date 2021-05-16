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
	"encoding/json"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Ping struct{}

var (
	_ translation.Decodable    = (*Ping)(nil)
	_ translation.Translatable = (*Ping)(nil)
)

// Spec implements translation.Decodable.
func (*Ping) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"schedule": &hcldec.AttrSpec{
			Name:     "schedule",
			Type:     cty.String,
			Required: false,
		},
		"data": &hcldec.AttrSpec{
			Name:     "data",
			Type:     cty.String,
			Required: true,
		},
		"content_type": &hcldec.AttrSpec{
			Name:     "content_type",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Ping) Manifests(id string, config, eventDst cty.Value) []interface{} {
	const defaultSchedule = "* * * * *" // every minute

	var manifests []interface{}

	s := k8s.NewObject("sources.knative.dev/v1beta2", "PingSource", k8s.RFC1123Name(id))

	schedule := defaultSchedule
	if v := config.GetAttr("schedule"); !v.IsNull() {
		schedule = v.AsString()
	}
	s.SetNestedField(schedule, "spec", "schedule")

	data := config.GetAttr("data").AsString()
	s.SetNestedField(data, "spec", "data")

	if v := config.GetAttr("content_type"); !v.IsNull() {
		contentType := v.AsString()
		s.SetNestedField(contentType, "spec", "contentType")
	} else if json.Valid([]byte(data)) {
		contentType := "application/json"
		s.SetNestedField(contentType, "spec", "contentType")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
