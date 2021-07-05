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

	"til/config/globals"
	"til/internal/sdk"
	"til/internal/sdk/k8s"
	"til/translation"
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
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Function) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	f := k8s.NewObject(k8s.APIExt, "Function", name)

	runtime := config.GetAttr("runtime").AsString()
	f.SetNestedField(runtime, "spec", "runtime")

	code := config.GetAttr("code").AsString()
	f.SetNestedField(code, "spec", "code")

	entrypoint := "main"
	if v := config.GetAttr("entrypoint"); !v.IsNull() {
		entrypoint = v.AsString()
	}
	f.SetNestedField(entrypoint, "spec", "entrypoint")

	if config.GetAttr("public").True() {
		f.SetNestedField(true, "spec", "public")
	}

	if v := config.GetAttr("ce_context"); !v.IsNull() {
		ceCtx := sdk.DecodeStringMap(v)
		f.SetNestedMap(ceCtx, "spec", "ceOverrides", "extensions")
	}

	manifests = append(manifests, f.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APIExt, "Function", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Function) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APIExt, "Function", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
