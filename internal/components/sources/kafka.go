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

package sources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config/globals"
	"bridgedl/internal/sdk"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Kafka struct{}

var (
	_ translation.Decodable    = (*Kafka)(nil)
	_ translation.Translatable = (*Kafka)(nil)
)

// Spec implements translation.Decodable.
func (*Kafka) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"consumer_group": &hcldec.AttrSpec{
			Name:     "consumer_group",
			Type:     cty.String,
			Required: false,
		},
		"bootstrap_servers": &hcldec.AttrSpec{
			Name:     "bootstrap_servers",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"topics": &hcldec.AttrSpec{
			Name:     "topics",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"sasl_auth": &hcldec.AttrSpec{
			Name:     "sasl_auth",
			Type:     k8s.ObjectReferenceCty,
			Required: false,
		},
		"tls": &hcldec.ValidateSpec{
			Wrapped: &hcldec.AttrSpec{
				Name:     "tls",
				Type:     cty.DynamicPseudoType,
				Required: false,
			},
			Func: validateKafkaAttrTLS,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Kafka) Manifests(id string, config, eventDst cty.Value, _ globals.Accessor) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.knative.dev/v1beta1", "KafkaSource", k8s.RFC1123Name(id))

	if v := config.GetAttr("consumer_group"); !v.IsNull() {
		consumerGroup := v.AsString()
		s.SetNestedField(consumerGroup, "spec", "consumerGroup")
	}

	bootstrapServers := sdk.DecodeStringSlice(config.GetAttr("bootstrap_servers"))
	s.SetNestedSlice(bootstrapServers, "spec", "bootstrapServers")

	topics := sdk.DecodeStringSlice(config.GetAttr("topics"))
	s.SetNestedSlice(topics, "spec", "topics")

	if v := config.GetAttr("sasl_auth"); !v.IsNull() {
		saslAuthSecretName := v.GetAttr("name").AsString()
		saslMech, saslUser, saslPasswd, _, _, _ := secrets.SecretKeyRefsKafka(saslAuthSecretName)
		s.SetNestedField(true, "spec", "net", "sasl", "enable")
		s.SetNestedMap(saslMech, "spec", "net", "sasl", "type", "secretKeyRef")
		s.SetNestedMap(saslUser, "spec", "net", "sasl", "user", "secretKeyRef")
		s.SetNestedMap(saslPasswd, "spec", "net", "sasl", "password", "secretKeyRef")
	}

	if v := config.GetAttr("tls"); !v.IsNull() {
		if k8s.IsObjectReference(v) {
			tlsSecretName := v.GetAttr("name").AsString()
			_, _, _, caCert, cert, key := secrets.SecretKeyRefsKafka(tlsSecretName)
			s.SetNestedField(true, "spec", "net", "tls", "enable")
			s.SetNestedMap(caCert, "spec", "net", "tls", "caCert", "secretKeyRef")
			s.SetNestedMap(cert, "spec", "net", "tls", "cert", "secretKeyRef")
			s.SetNestedMap(key, "spec", "net", "tls", "key", "secretKeyRef")
			// The protocol selection happens at runtime, based on the
			// presence or not of the above keys in the referenced Secret.
			// By marking each of these keys as optional, we attempt to
			// provide configuration parity with the "kafka" target, which
			// uses this same approach.
			s.SetNestedField(true, "spec", "net", "tls", "caCert", "secretKeyRef", "optional")
			s.SetNestedField(true, "spec", "net", "tls", "cert", "secretKeyRef", "optional")
			s.SetNestedField(true, "spec", "net", "tls", "key", "secretKeyRef", "optional")
		} else if v.True() {
			s.SetNestedField(true, "spec", "net", "tls", "enable")
		}
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}

func validateKafkaAttrTLS(val cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if !(k8s.IsObjectReference(val) || val.Type() == cty.Bool) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid attribute type",
			Detail:   `The "tls" attribute accepts either a secret reference or a boolean.`,
		})
	}

	return diags
}
