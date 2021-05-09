package k8s

import (
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	APIEventing  = "eventing.knative.dev/v1"
	APIMessaging = "messaging.knative.dev/v1"
)

// DecodeDestination returns a JSON-serializable representation of a the given
// k8s.DestinationCty, in a format that is compatible with the "unstructured"
// package from k8s.io/apimachinery.
// Panics if the receiver is not a non-null k8s.DestinationCty.
//
// This helper is intended to be used for populating "sink" attributes, or
// other attributes with similar semantics.
func DecodeDestination(dst cty.Value) map[string]interface{} {
	dstRef := dst.GetAttr("ref")
	out := map[string]interface{}{
		"apiVersion": dstRef.GetAttr("apiVersion").AsString(),
		"kind":       dstRef.GetAttr("kind").AsString(),
		"name":       dstRef.GetAttr("name").AsString(),
	}

	return out
}

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
	validateDNS1123Subdomain(name)

	t := &unstructured.Unstructured{}

	t.SetAPIVersion(APIEventing)
	t.SetKind("Trigger")
	t.SetName(name)

	_ = unstructured.SetNestedField(t.Object, broker, "spec", "broker")

	sink := DecodeDestination(dst)
	_ = unstructured.SetNestedMap(t.Object, sink, "spec", "subscriber", "ref")

	if len(filter) != 0 {
		_ = unstructured.SetNestedMap(t.Object, filter, "spec", "filter", "attributes")
	}

	return t
}

// NewChannel returns a new Knative Channel.
func NewChannel(name string) *unstructured.Unstructured {
	validateDNS1123Subdomain(name)

	c := &unstructured.Unstructured{}

	c.SetAPIVersion(APIMessaging)
	c.SetKind("Channel")
	c.SetName(name)

	return c
}

// NewSubscription returns a new Knative Subscription.
func NewSubscription(name, channel string, dst, replyDst cty.Value) *unstructured.Unstructured {
	validateDNS1123Subdomain(name)

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

	sink := DecodeDestination(dst)
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "subscriber", "ref")

	if !replyDst.IsNull() {
		reply := DecodeDestination(replyDst)
		_ = unstructured.SetNestedMap(s.Object, reply, "spec", "reply", "ref")
	}

	return s
}
