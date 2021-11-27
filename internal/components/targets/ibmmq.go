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

type IBMMQ struct{}

var (
	_ translation.Decodable    = (*IBMMQ)(nil)
	_ translation.Translatable = (*IBMMQ)(nil)
	_ translation.Addressable  = (*IBMMQ)(nil)
)

// Spec implements translation.Decodable.
func (*IBMMQ) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"connection_name": &hcldec.AttrSpec{
			Name:     "connection_name",
			Type:     cty.String,
			Required: true,
		},
		"queue_manager": &hcldec.AttrSpec{
			Name:     "queue_manager",
			Type:     cty.String,
			Required: true,
		},
		"queue_name": &hcldec.AttrSpec{
			Name:     "queue_name",
			Type:     cty.String,
			Required: true,
		},
		"channel_name": &hcldec.AttrSpec{
			Name:     "channel_name",
			Type:     cty.String,
			Required: true,
		},
		"credentials": &hcldec.AttrSpec{
			Name:     "credentials",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*IBMMQ) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "IBMMQTarget", name)

	connectionName := config.GetAttr("connection_name").AsString()
	t.SetNestedField(connectionName, "spec", "connectionName")

	qMgr := config.GetAttr("queue_manager").AsString()
	t.SetNestedField(qMgr, "spec", "queueManager")

	qName := config.GetAttr("queue_name").AsString()
	t.SetNestedField(qName, "spec", "queueName")

	cName := config.GetAttr("channel_name").AsString()
	t.SetNestedField(cName, "spec", "channelName")

	secret := config.GetAttr("credentials").GetAttr("name").AsString()
	username, password := secrets.SecretKeyRefsBasicAuth(secret)
	t.SetNestedMap(username, "spec", "credentials", "username", "valueFromSecret")
	t.SetNestedMap(password, "spec", "credentials", "password", "valueFromSecret")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)

		subscriber := k8s.NewDestination(k8s.APITargets, "IBMMQTarget", name)

		sbOpts := []k8s.SubscriptionOption{k8s.ReplyDest(eventDst)}
		sbOpts = k8s.AppendDeliverySubscriptionOptions(sbOpts, glb)

		subs := k8s.NewSubscription(name, name, subscriber, sbOpts...)

		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*IBMMQ) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "IBMMQTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
