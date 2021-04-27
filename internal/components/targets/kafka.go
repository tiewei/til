package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
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

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("eventing.knative.dev/v1alpha1")
	s.SetKind("KafkaSink")
	s.SetName(k8s.RFC1123Name(id))

	topic := config.GetAttr("topic").AsString()
	_ = unstructured.SetNestedField(s.Object, topic, "spec", "topic")

	var bootstrapServers []interface{}
	bSrvsIter := config.GetAttr("bootstrap_servers").ElementIterator()
	for bSrvsIter.Next() {
		_, srv := bSrvsIter.Element()
		bootstrapServers = append(bootstrapServers, srv.AsString())
	}
	_ = unstructured.SetNestedSlice(s.Object, bootstrapServers, "spec", "bootstrapServers")

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	_ = unstructured.SetNestedField(s.Object, authSecretName, "spec", "auth", "secret", "ref", "name")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Kafka) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("eventing.knative.dev/v1alpha1", "KafkaSink", k8s.RFC1123Name(id))
}
