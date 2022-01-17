/*
Copyright 2022 TriggerMesh Inc.

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
		"discard_ce_context": &hcldec.AttrSpec{
			Name:     "discard_ce_context",
			Type:     cty.Bool,
			Required: false,
		},
		"synchronous": &hcldec.BlockSpec{
			TypeName: "synchronous",
			Nested: &hcldec.ObjectSpec{
				"request_key": &hcldec.AttrSpec{
					Name:     "request_key",
					Type:     cty.String,
					Required: true,
				},
				"response_correlation_key": &hcldec.AttrSpec{
					Name:     "response_correlation_key",
					Type:     cty.String,
					Required: true,
				},
				"response_wait_timeout": &hcldec.AttrSpec{
					Name:     "response_wait_timeout",
					Type:     cty.Number,
					Required: false,
				},
			},
		},
		"reply_to": &hcldec.BlockSpec{
			TypeName: "reply_to",
			Nested: &hcldec.ObjectSpec{
				"manager": &hcldec.AttrSpec{
					Name:     "manager",
					Type:     cty.String,
					Required: false,
				},
				"queue": &hcldec.AttrSpec{
					Name:     "queue",
					Type:     cty.String,
					Required: false,
				},
			},
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

	if v := config.GetAttr("reply_to"); !v.IsNull() {
		if replyToQMgr := v.GetAttr("manager"); !replyToQMgr.IsNull() {
			t.SetNestedField(replyToQMgr.AsString(), "spec", "replyTo", "queueManager")
		}

		if replyToQueue := v.GetAttr("queue"); !replyToQueue.IsNull() {
			t.SetNestedField(replyToQueue.AsString(), "spec", "replyTo", "queueName")
		}
	}

	if config.GetAttr("discard_ce_context").True() {
		t.SetNestedField(true, "spec", "discardCloudEventContext")
	}

	secret := config.GetAttr("credentials").GetAttr("name").AsString()
	username, password := secrets.SecretKeyRefsBasicAuth(secret)
	t.SetNestedMap(username, "spec", "credentials", "username", "valueFromSecret")
	t.SetNestedMap(password, "spec", "credentials", "password", "valueFromSecret")

	manifests = append(manifests, t.Unstructured())

	if sync := config.GetAttr("synchronous"); !sync.IsNull() {
		parent := k8s.NewDestination(k8s.APITargets, "IBMMQTarget", name)
		s := synchronizer(name, sync, parent)
		manifests = append(manifests, s.Unstructured())
	}

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
func (m *IBMMQ) Address(id string, config, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if sync := config.GetAttr("synchronous"); !sync.IsNull() {
		return k8s.NewDestination(k8s.APIFlow, "Synchronizer", name)
	}

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "IBMMQTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}

func synchronizer(name string, config cty.Value, eventDst cty.Value) *k8s.Object {
	s := k8s.NewObject(k8s.APIFlow, "Synchronizer", name)

	requestCorrelationKey := "extensions.correlationid[0:23]"
	if requestKey := config.GetAttr("request_key"); !requestKey.IsNull() {
		requestCorrelationKey = requestKey.AsString()
	}
	s.SetNestedField(requestCorrelationKey, "spec", "requestKey")

	responseCorrelationKey := "extensions.correlationid"
	if responseKey := config.GetAttr("response_correlation_key"); !responseKey.IsNull() {
		responseCorrelationKey = responseKey.AsString()
	}
	s.SetNestedField(responseCorrelationKey, "spec", "responseCorrelationKey")

	s.SetNestedField(int64(10), "spec", "responseWaitTimeout")
	if timeout := config.GetAttr("response_wait_timeout"); !timeout.IsNull() {
		t, _ := timeout.AsBigFloat().Int64()
		s.SetNestedField(t, "spec", "responseWaitTimeout")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return s
}
