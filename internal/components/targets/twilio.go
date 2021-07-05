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
	"til/internal/sdk/k8s"
	"til/internal/sdk/secrets"
	"til/translation"
)

type Twilio struct{}

var (
	_ translation.Decodable    = (*Twilio)(nil)
	_ translation.Translatable = (*Twilio)(nil)
	_ translation.Addressable  = (*Twilio)(nil)
)

// Spec implements translation.Decodable.
func (*Twilio) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"default_phone_from": &hcldec.AttrSpec{
			Name:     "default_phone_from",
			Type:     cty.String,
			Required: false,
		},
		"default_phone_to": &hcldec.AttrSpec{
			Name:     "default_phone_to",
			Type:     cty.String,
			Required: false,
		},
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Twilio) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "TwilioTarget", name)

	if v := config.GetAttr("default_phone_from"); !v.IsNull() {
		defaultPhoneFrom := config.GetAttr("default_phone_from").AsString()
		t.SetNestedField(defaultPhoneFrom, "spec", "defaultPhoneFrom")
	}

	if v := config.GetAttr("default_phone_to"); !v.IsNull() {
		defaultPhoneTo := config.GetAttr("default_phone_to").AsString()
		t.SetNestedField(defaultPhoneTo, "spec", "defaultPhoneTo")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	sidSecretRef, tokenSecretRef := secrets.SecretKeyRefsTwilio(authSecretName)
	t.SetNestedMap(sidSecretRef, "spec", "sid", "secretKeyRef")
	t.SetNestedMap(tokenSecretRef, "spec", "token", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APITargets, "TwilioTarget", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Twilio) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "TwilioTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
