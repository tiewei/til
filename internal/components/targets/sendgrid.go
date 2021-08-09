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

type Sendgrid struct{}

var (
	_ translation.Decodable    = (*Sendgrid)(nil)
	_ translation.Translatable = (*Sendgrid)(nil)
	_ translation.Addressable  = (*Sendgrid)(nil)
)

// Spec implements translation.Decodable.
func (*Sendgrid) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"default_from_email": &hcldec.AttrSpec{
			Name:     "default_from_email",
			Type:     cty.String,
			Required: false,
		},
		"default_from_name": &hcldec.AttrSpec{
			Name:     "default_from_name",
			Type:     cty.String,
			Required: false,
		},
		"default_subject": &hcldec.AttrSpec{
			Name:     "default_subject",
			Type:     cty.String,
			Required: false,
		},
		"default_to_email": &hcldec.AttrSpec{
			Name:     "default_to_email",
			Type:     cty.String,
			Required: false,
		},
		"default_to_name": &hcldec.AttrSpec{
			Name:     "default_to_name",
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
func (*Sendgrid) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "SendgridTarget", name)

	if v := config.GetAttr("default_from_email"); !v.IsNull() {
		fromEmail := v.AsString()
		t.SetNestedField(fromEmail, "spec", "defaultFromEmail")
	}

	if v := config.GetAttr("default_from_name"); !v.IsNull() {
		fromName := v.AsString()
		t.SetNestedField(fromName, "spec", "defaultFromName")
	}

	if v := config.GetAttr("default_subject"); !v.IsNull() {
		subject := v.AsString()
		t.SetNestedField(subject, "spec", "defaultSubject")
	}

	if v := config.GetAttr("default_to_email"); !v.IsNull() {
		toEmail := v.AsString()
		t.SetNestedField(toEmail, "spec", "defaultToEmail")
	}

	if v := config.GetAttr("default_to_name"); !v.IsNull() {
		toName := v.AsString()
		t.SetNestedField(toName, "spec", "defaultToName")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	apiKeySecretRef := secrets.SecretKeyRefsSendgrid(authSecretName)
	t.SetNestedMap(apiKeySecretRef, "spec", "apiKey", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APITargets, "SendgridTarget", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Sendgrid) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "SendgridTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
