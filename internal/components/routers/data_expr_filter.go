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

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type DataExprFilter struct{}

var (
	_ translation.Decodable    = (*DataExprFilter)(nil)
	_ translation.Translatable = (*DataExprFilter)(nil)
	_ translation.Addressable  = (*DataExprFilter)(nil)
)

// Spec implements translation.Decodable.
func (*DataExprFilter) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"condition": &hcldec.AttrSpec{
			Name:     "condition",
			Type:     cty.String,
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
func (*DataExprFilter) Manifests(id string, config, _ cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	eventDst := config.GetAttr("to")
	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	f := k8s.NewObject(k8s.APIFlow, "Filter", name)

	expr := config.GetAttr("condition").AsString()
	f.SetNestedField(expr, "spec", "expression")

	sink := k8s.DecodeDestination(eventDst)
	f.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, f.Unstructured())
}

// Address implements translation.Addressable.
func (*DataExprFilter) Address(id string, _, _ cty.Value, _ globals.Accessor) cty.Value {
	return k8s.NewDestination(k8s.APIFlow, "Filter", k8s.RFC1123Name(id))
}
