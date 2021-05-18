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

package transformers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Function struct{}

var (
	_ translation.Decodable    = (*Function)(nil)
	_ translation.Translatable = (*Function)(nil)
	_ translation.Addressable  = (*Function)(nil)
)

// Spec implements translation.Decodable.
func (*Function) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"runtime": &hcldec.AttrSpec{
			Name:     "runtime",
			Type:     cty.String,
			Required: true,
		},
		"code": &hcldec.AttrSpec{
			Name:     "code",
			Type:     cty.String,
			Required: true,
		},
		"entrypoint": &hcldec.AttrSpec{
			Name:     "entrypoint",
			Type:     cty.String,
			Required: false,
		},
		"public": &hcldec.AttrSpec{
			Name:     "public",
			Type:     cty.Bool,
			Required: false,
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
					Required: false,
				},
				"subject": &hcldec.AttrSpec{
					Name:     "subject",
					Type:     cty.String,
					Required: false,
				},
			},
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Function) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	code := config.GetAttr("code").AsString()

	switch runtime := config.GetAttr("runtime").AsString(); runtime {
	case "js":
		t := k8s.NewObject(k8s.APITargets, "InfraTarget", name)

		t.SetNestedField(code, "spec", "script", "code")

		// route responses via a channel subscription
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "InfraTarget", name), eventDst)

		manifests = append(manifests, t.Unstructured(), ch, subs)

	default:
		f := k8s.NewObject(k8s.APIExt, "Function", name)

		f.SetNestedField(runtime, "spec", "runtime")
		f.SetNestedField(code, "spec", "code")

		entrypoint := "main"
		if v := config.GetAttr("entrypoint"); !v.IsNull() {
			entrypoint = v.AsString()
		}
		f.SetNestedField(entrypoint, "spec", "entrypoint")

		if config.GetAttr("public").True() {
			f.SetNestedField(true, "spec", "public")
		}

		ceCtx := sdk.DecodeStringMap(config.GetAttr("ce_context"))
		f.SetNestedMap(ceCtx, "spec", "ceOverrides", "extensions")

		sink := k8s.DecodeDestination(eventDst)
		f.SetNestedMap(sink, "spec", "sink", "ref")

		manifests = append(manifests, f.Unstructured())
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Function) Address(id string, config, _ cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if config.GetAttr("runtime").AsString() == "js" {
		return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
	}
	return k8s.NewDestination(k8s.APIExt, "Function", name)
}
