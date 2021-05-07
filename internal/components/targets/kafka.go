package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Kafka struct{}

var (
	_ translation.Decodable    = (*Kafka)(nil)
	_ translation.Translatable = (*Kafka)(nil)
	_ translation.Addressable  = (*Kafka)(nil)
)

// Spec implements translation.Decodable.
func (*Kafka) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"topic": &hcldec.AttrSpec{
			Name:     "topic",
			Type:     cty.String,
			Required: true,
		},
		"bootstrap_servers": &hcldec.AttrSpec{
			Name:     "bootstrap_servers",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Kafka) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("eventing.knative.dev/v1alpha1", "KafkaSink", id)

	topic := config.GetAttr("topic").AsString()
	s.SetNestedField(topic, "spec", "topic")

	bootstrapServers := sdk.DecodeStringSlice(config.GetAttr("bootstrap_servers"))
	s.SetNestedSlice(bootstrapServers, "spec", "bootstrapServers")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	s.SetNestedField(authSecretName, "spec", "auth", "secret", "ref", "name")

	return append(manifests, s.Unstructured())
}

// Address implements translation.Addressable.
func (*Kafka) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("eventing.knative.dev/v1alpha1", "KafkaSink", k8s.RFC1123Name(id))
}
