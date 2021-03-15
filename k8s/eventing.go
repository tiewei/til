package k8s

import (
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// NewBroker returns a new Knative Broker.
func NewBroker(name string) *unstructured.Unstructured {
	b := &unstructured.Unstructured{}

	b.SetAPIVersion("eventing.knative.dev/v1")
	b.SetKind("Broker")
	b.SetName(name)

	return b
}

// NewTrigger returns a new Knative Trigger.
func NewTrigger(name, broker string, dst cty.Value, filter map[string]interface{}) *unstructured.Unstructured {
	t := &unstructured.Unstructured{}

	t.SetAPIVersion("eventing.knative.dev/v1")
	t.SetKind("Trigger")
	t.SetName(name)

	_ = unstructured.SetNestedField(t.Object, broker, "spec", "broker")

	sinkRef := dst.GetAttr("ref")

	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(t.Object, sink, "spec", "subscriber", "ref")

	if len(filter) != 0 {
		_ = unstructured.SetNestedMap(t.Object, filter, "spec", "filter", "attributes")
	}

	return t
}
