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

package routers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/validation"
	"bridgedl/translation"
)

type Splitter struct{}

var (
	_ translation.Decodable    = (*Splitter)(nil)
	_ translation.Translatable = (*Splitter)(nil)
	_ translation.Addressable  = (*Splitter)(nil)
)

// Spec implements translation.Decodable.
func (*Splitter) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"path": &hcldec.AttrSpec{
			Name:     "path",
			Type:     cty.String,
			Required: true,
		},
		"ce_context": &hcldec.BlockSpec{
			TypeName: "ce_context",
			Nested: &hcldec.ObjectSpec{
				"type": &hcldec.AttrSpec{
					Name:     "type",
					Type:     cty.String,
					Required: true,
				},
				"source": &hcldec.AttrSpec{
					Name:     "source",
					Type:     cty.String,
					Required: true,
				},
				"extensions": &hcldec.ValidateSpec{
					Wrapped: &hcldec.AttrSpec{
						Name:     "extensions",
						Type:     cty.Map(cty.String),
						Required: false,
					},
					Func: validation.ContainsCEContextAttributes,
				},
			},
			Required: true,
		},
		"to": &hcldec.AttrSpec{
			Name:     "to",
			Type:     k8s.DestinationCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Splitter) Manifests(id string, config, _ cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("routing.triggermesh.io/v1alpha1", "Splitter", k8s.RFC1123Name(id))

	path := config.GetAttr("path").AsString()
	s.SetNestedField(path, "spec", "path")

	typ := config.GetAttr("ce_context").GetAttr("type").AsString()
	s.SetNestedField(typ, "spec", "ceContext", "type")

	src := config.GetAttr("ce_context").GetAttr("source").AsString()
	s.SetNestedField(src, "spec", "ceContext", "source")

	if v := config.GetAttr("ce_context").GetAttr("extensions"); !v.IsNull() {
		ceExts := sdk.DecodeStringMap(v)
		s.SetNestedMap(ceExts, "spec", "ceContext", "extensions")
	}

	sink := k8s.DecodeDestination(config.GetAttr("to"))
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}

// Address implements translation.Addressable.
func (*Splitter) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("routing.triggermesh.io/v1alpha1", "Splitter", k8s.RFC1123Name(id))
}
