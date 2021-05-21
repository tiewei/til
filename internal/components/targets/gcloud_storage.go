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

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type GCloudStorage struct{}

var (
	_ translation.Decodable    = (*GCloudStorage)(nil)
	_ translation.Translatable = (*GCloudStorage)(nil)
	_ translation.Addressable  = (*GCloudStorage)(nil)
)

// Spec implements translation.Decodable.
func (*GCloudStorage) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"bucket_name": &hcldec.AttrSpec{
			Name:     "bucket_name",
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
func (*GCloudStorage) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "GoogleCloudStorageTarget", name)

	bucketName := config.GetAttr("bucket_name").AsString()
	t.SetNestedField(bucketName, "spec", "bucketName")

	svcAccountSecretName := config.GetAttr("service_account").GetAttr("name").AsString()
	keySecretRef := secrets.SecretKeyRefsGCloudServiceAccount(svcAccountSecretName)
	t.SetNestedMap(keySecretRef, "spec", "credentialsJson", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APITargets, "GoogleCloudStorageTarget", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*GCloudStorage) Address(id string, _, eventDst cty.Value, _ globals.Accessor) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "GoogleCloudStorageTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
