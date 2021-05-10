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

	name := k8s.RFC1123Name(id)

	s := k8s.NewObject(k8s.APIEventingV1Alpha1, "KafkaSink", name)

	topic := config.GetAttr("topic").AsString()
	s.SetNestedField(topic, "spec", "topic")

	bootstrapServers := sdk.DecodeStringSlice(config.GetAttr("bootstrap_servers"))
	s.SetNestedSlice(bootstrapServers, "spec", "bootstrapServers")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	s.SetNestedField(authSecretName, "spec", "auth", "secret", "ref", "name")

	manifests = append(manifests, s.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APIEventingV1Alpha1, "KafkaSink", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Kafka) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APIEventingV1Alpha1, "KafkaSink", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
