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

package k8s

import (
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/config/globals"
	"bridgedl/lang/k8s"
)

const (
	APIEventing         = "eventing.knative.dev/v1"
	APIEventingV1Alpha1 = "eventing.knative.dev/v1alpha1"
	APIMessaging        = "messaging.knative.dev/v1"
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
func NewTrigger(name, broker string, dst cty.Value, opts ...TriggerOption) *unstructured.Unstructured {
	validateDNS1123Subdomain(name)

	t := &unstructured.Unstructured{}

	t.SetAPIVersion(APIEventing)
	t.SetKind("Trigger")
	t.SetName(name)

	_ = unstructured.SetNestedField(t.Object, broker, "spec", "broker")

	sink := DecodeDestination(dst)
	_ = unstructured.SetNestedMap(t.Object, sink, "spec", "subscriber", "ref")

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// TriggerOption is a functional option of a Knative Trigger.
type TriggerOption func(*unstructured.Unstructured)

// Filter sets the context attributes to filter on.
func Filter(filter map[string]interface{}) TriggerOption {
	return func(o *unstructured.Unstructured) {
		if len(filter) != 0 {
			_ = unstructured.SetNestedMap(o.Object, filter, "spec", "filter", "attributes")
		}
	}
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
func NewSubscription(name, channel string, dst cty.Value, opts ...SubscriptionOption) *unstructured.Unstructured {
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

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SubscriptionOption is a functional option of a Knative Subscription.
type SubscriptionOption func(*unstructured.Unstructured)

// ReplyDest sets the destination of event replies.
func ReplyDest(replyDst cty.Value) SubscriptionOption {
	return func(o *unstructured.Unstructured) {
		if !replyDst.IsNull() {
			reply := DecodeDestination(replyDst)
			_ = unstructured.SetNestedMap(o.Object, reply, "spec", "reply", "ref")
		}
	}
}

// Retries sets the number of delivery retries.
func Retries(n int64) SubscriptionOption {
	return func(o *unstructured.Unstructured) {
		if n >= 0 {
			_ = unstructured.SetNestedField(o.Object, n, "spec", "delivery", "retry")
		}
	}
}

// DeadLetterSink sets the dead-letter sink.
func DeadLetterSink(deadletterDst cty.Value) SubscriptionOption {
	return func(o *unstructured.Unstructured) {
		if !deadletterDst.IsNull() {
			dls := DecodeDestination(deadletterDst)
			_ = unstructured.SetNestedMap(o.Object, dls, "spec", "delivery", "deadLetterSink", "ref")
		}
	}
}

// MaybeAppendChannel conditionally appends a Channel and Subscription to a
// list of manifests, based on the global settings provided by the given
// globals.Accessor.
// If a Channel and Subscription are indeed appended, the returned eventDst is
// a duck Destination matching the Channel.
//
// The purpose of this helper is to ease the implementation of global delivery
// settings across component implementations.
func MaybeAppendChannel(name string, manifests []interface{}, eventDst cty.Value,
	glb globals.Accessor) (newManifests []interface{}, newEventDst cty.Value) {

	d := glb.Delivery()

	// do not interpolate a Channel+Subscription if global delivery
	// settings were omitted
	if d == nil {
		return manifests, eventDst
	}

	name = name + "-" + eventDst.GetAttr("ref").GetAttr("name").AsString()

	ch := NewChannel(name)
	manifests = append(manifests, ch)

	var sbOpts []SubscriptionOption
	sbOpts = AppendDeliverySubscriptionOptions(sbOpts, glb)

	sb := NewSubscription(name, name, eventDst, sbOpts...)
	manifests = append(manifests, sb)

	// a Channel+Subscription were interpolated, the channel becomes the
	// new event destination
	eventDst = k8s.NewDestination(APIMessaging, "Channel", name)

	return manifests, eventDst
}

// AppendDeliverySubscriptionOptions conditionally appends delivery-related
// options to the given list of SubscriptionOption, based on the global
// settings provided by the given globals.Accessor.
func AppendDeliverySubscriptionOptions(sbOpts []SubscriptionOption, glb globals.Accessor) []SubscriptionOption {
	d := glb.Delivery()

	if d == nil {
		return sbOpts
	}

	if r := d.Retries; r != nil {
		sbOpts = append(sbOpts, Retries(*r))
	}
	if dls := d.DeadLetterSink; !dls.IsNull() && dls.IsKnown() {
		sbOpts = append(sbOpts, DeadLetterSink(dls))
	}

	return sbOpts
}
