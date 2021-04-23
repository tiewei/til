package sources

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
	}
}

// Manifests implements translation.Translatable.
func (*Kafka) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.knative.dev/v1beta1")
	s.SetKind("KafkaSource")
	s.SetName(k8s.RFC1123Name(id))

	if v := config.GetAttr("consumer_group"); !v.IsNull() {
		consumerGroup := v.AsString()
		_ = unstructured.SetNestedField(s.Object, consumerGroup, "spec", "consumerGroup")
	}

	var bootstrapServers []interface{}
	bSrvsIter := config.GetAttr("bootstrap_servers").ElementIterator()
	for bSrvsIter.Next() {
		_, srv := bSrvsIter.Element()
		bootstrapServers = append(bootstrapServers, srv.AsString())
	}
	_ = unstructured.SetNestedSlice(s.Object, bootstrapServers, "spec", "bootstrapServers")

	var topics []interface{}
	topicsIter := config.GetAttr("topics").ElementIterator()
	for topicsIter.Next() {
		_, topic := topicsIter.Element()
		topics = append(topics, topic.AsString())
	}
	_ = unstructured.SetNestedSlice(s.Object, topics, "spec", "topics")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
