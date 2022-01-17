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

package sources

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
		"delivery": &hcldec.BlockSpec{
			TypeName: "delivery",
			Nested: &hcldec.ObjectSpec{
				"dead_letter_queue_manager": &hcldec.AttrSpec{
					Name:     "dead_letter_queue_manager",
					Type:     cty.String,
					Required: false,
				},
				"dead_letter_queue": &hcldec.AttrSpec{
					Name:     "dead_letter_queue",
					Type:     cty.String,
					Required: false,
				},
				"retry": &hcldec.AttrSpec{
					Name:     "retry",
					Type:     cty.Number,
					Required: false,
				},
			},
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*IBMMQ) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	s := k8s.NewObject(k8s.APISources, "IBMMQSource", name)

	connectionName := config.GetAttr("connection_name").AsString()
	s.SetNestedField(connectionName, "spec", "connectionName")

	qMgr := config.GetAttr("queue_manager").AsString()
	s.SetNestedField(qMgr, "spec", "queueManager")

	qName := config.GetAttr("queue_name").AsString()
	s.SetNestedField(qName, "spec", "queueName")

	cName := config.GetAttr("channel_name").AsString()
	s.SetNestedField(cName, "spec", "channelName")

	secret := config.GetAttr("credentials").GetAttr("name").AsString()
	username, password := secrets.SecretKeyRefsBasicAuth(secret)
	s.SetNestedMap(username, "spec", "credentials", "username", "valueFromSecret")
	s.SetNestedMap(password, "spec", "credentials", "password", "valueFromSecret")

	if delivery := config.GetAttr("delivery"); !delivery.IsNull() {
		if dlqm := delivery.GetAttr("dead_letter_queue_manager"); !dlqm.IsNull() {
			s.SetNestedField(dlqm.AsString(), "spec", "delivery", "deadLetterQueueManager")
		}
		if dlq := delivery.GetAttr("dead_letter_queue"); !dlq.IsNull() {
			s.SetNestedField(dlq.AsString(), "spec", "delivery", "deadLetterQueue")
		}
		if retry := delivery.GetAttr("retry"); !retry.IsNull() {
			r, _ := retry.AsBigFloat().Int64()
			s.SetNestedField(r, "spec", "delivery", "retry")
		}
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
