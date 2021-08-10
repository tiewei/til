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

type GCloudWorkflows struct{}

var (
	_ translation.Decodable    = (*GCloudWorkflows)(nil)
	_ translation.Translatable = (*GCloudWorkflows)(nil)
	_ translation.Addressable  = (*GCloudWorkflows)(nil)
)

// Spec implements translation.Decodable.
func (*GCloudWorkflows) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"service_account": &hcldec.AttrSpec{
			Name:     "service_account",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*GCloudWorkflows) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "GoogleCloudWorkflowsTarget", name)

	svcAccountSecretName := config.GetAttr("service_account").GetAttr("name").AsString()
	keySecretRef := secrets.SecretKeyRefsGCloudServiceAccount(svcAccountSecretName)
	t.SetNestedMap(keySecretRef, "spec", "credentialsJson", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APITargets, "GoogleCloudWorkflowsTarget", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*GCloudWorkflows) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "GoogleCloudWorkflowsTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
