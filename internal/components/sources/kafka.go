package sources

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

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
func (*Kafka) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.knative.dev/v1beta1", "KafkaSource", id)

	if v := config.GetAttr("consumer_group"); !v.IsNull() {
		consumerGroup := v.AsString()
		s.SetNestedField(consumerGroup, "spec", "consumerGroup")
	}

	bootstrapServersVals := config.GetAttr("bootstrap_servers").AsValueSlice()
	bootstrapServers := make([]interface{}, 0, len(bootstrapServersVals))
	for _, v := range bootstrapServersVals {
		bootstrapServers = append(bootstrapServers, v.AsString())
	}
	s.SetNestedSlice(bootstrapServers, "spec", "bootstrapServers")

	topicsVals := config.GetAttr("topics").AsValueSlice()
	topics := make([]interface{}, 0, len(topicsVals))
	for _, v := range topicsVals {
		topics = append(topics, v.AsString())
	}
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

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}

func validateKafkaAttrTLS(val cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if !(k8s.IsObjectReference(val) || val.Type() == cty.Bool) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid attributes type",
			Detail:   `The "tls" attribute accepts either a secret reference or a boolean.`,
		})
	}

	return diags
}
