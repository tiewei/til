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

type GCloudFirestore struct{}

var (
	_ translation.Decodable    = (*GCloudFirestore)(nil)
	_ translation.Translatable = (*GCloudFirestore)(nil)
	_ translation.Addressable  = (*GCloudFirestore)(nil)
)

// Spec implements translation.Decodable.
func (*GCloudFirestore) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"default_collection": &hcldec.AttrSpec{
			Name:     "default_collection",
			Type:     cty.String,
			Required: true,
		},
		"project_id": &hcldec.AttrSpec{
			Name:     "project_id",
			Type:     cty.String,
			Required: true,
		},
		"service_account": &hcldec.AttrSpec{
			Name:     "service_account",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*GCloudFirestore) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "GoogleCloudFirestoreTarget", name)

	defaultCollection := config.GetAttr("default_collection").AsString()
	t.SetNestedField(defaultCollection, "spec", "defaultCollection")

	projectID := config.GetAttr("project_id").AsString()
	t.SetNestedField(projectID, "spec", "projectID")

	svcAccountSecretName := config.GetAttr("service_account").GetAttr("name").AsString()
	keySecretRef := secrets.SecretKeyRefsGCloudServiceAccount(svcAccountSecretName)
	t.SetNestedMap(keySecretRef, "spec", "credentialsJson", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APITargets, "GoogleCloudFirestoreTarget", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*GCloudFirestore) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "GoogleCloudFirestoreTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
