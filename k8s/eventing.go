package k8s

import (
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	APIEventing  = "eventing.knative.dev/v1"
	APIMessaging = "messaging.knative.dev/v1"
)

// NewBroker returns a new Knative Broker.
func NewBroker(name string) *unstructured.Unstructured {
	b := &unstructured.Unstructured{}

	b.SetAPIVersion(APIEventing)
	b.SetKind("Broker")
	b.SetName(name)

	return b
}

// NewTrigger returns a new Knative Trigger.
func NewTrigger(name, broker string, dst cty.Value, filter map[string]interface{}) *unstructured.Unstructured {
	t := &unstructured.Unstructured{}

	t.SetAPIVersion(APIEventing)
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

// NewChannel returns a new Knative Channel.
func NewChannel(name string) *unstructured.Unstructured {
	s := &unstructured.Unstructured{}

	s.SetAPIVersion(APIMessaging)
	s.SetKind("Channel")
	s.SetName(name)

	return s
}

// NewSubscription returns a new Knative Subscription.
func NewSubscription(name, channel string, dst, replyDst cty.Value) *unstructured.Unstructured {
	s := &unstructured.Unstructured{}

	s.SetAPIVersion(APIMessaging)
	s.SetKind("Subscription")
	s.SetName(name)

	ch := map[string]interface{}{
		"apiVersion": APIMessaging,
		"kind":       "Channel",
		"name":       channel,
	}
	_ = unstructured.SetNestedMap(s.Object, ch, "spec", "channel")

	sinkRef := dst.GetAttr("ref")

	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "subscriber", "ref")

	replyRef := replyDst.GetAttr("ref")

	if !replyRef.IsNull() {
		reply := map[string]interface{}{
			"apiVersion": replyRef.GetAttr("apiVersion").AsString(),
			"kind":       replyRef.GetAttr("kind").AsString(),
			"name":       replyRef.GetAttr("name").AsString(),
		}
		_ = unstructured.SetNestedMap(s.Object, reply, "spec", "reply", "ref")
	}

	return s
}
