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

package channels

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/validation"
	"bridgedl/translation"
)

type PointToPoint struct{}

var (
	_ translation.Decodable    = (*PointToPoint)(nil)
	_ translation.Translatable = (*PointToPoint)(nil)
	_ translation.Addressable  = (*PointToPoint)(nil)
)

// Spec implements translation.Decodable.
func (*PointToPoint) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"delivery": &hcldec.BlockSpec{
			TypeName: "delivery",
			Nested: &hcldec.ObjectSpec{
				"retries": &hcldec.ValidateSpec{
					Wrapped: &hcldec.AttrSpec{
						Name:     "retries",
						Type:     cty.Number,
						Required: false,
					},
					Func: validation.IsInt,
				},
				"dead_letter_sink": &hcldec.AttrSpec{
					Name:     "dead_letter_sink",
					Type:     k8s.DestinationCty,
					Required: false,
				},
			},
			Required: false,
		},
		"to": &hcldec.AttrSpec{
			Name:     "to",
			Type:     k8s.DestinationCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*PointToPoint) Manifests(id string, config, _ cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	ch := k8s.NewChannel(name)
	manifests = append(manifests, ch)

	subscriber := config.GetAttr("to")

	var sbOpts []k8s.SubscriptionOption

	if delivery := config.GetAttr("delivery"); !delivery.IsNull() {
		if v := delivery.GetAttr("retries"); !v.IsNull() {
			retries, _ := v.AsBigFloat().Int64()
			sbOpts = append(sbOpts, k8s.Retries(retries))
		}
		dls := delivery.GetAttr("dead_letter_sink")
		sbOpts = append(sbOpts, k8s.DeadLetterSink(dls))
	}

	subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)
	manifests = append(manifests, subs)

	return manifests
}

// Address implements translation.Addressable.
func (*PointToPoint) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination(k8s.APIMessaging, "Channel", k8s.RFC1123Name(id))
}
