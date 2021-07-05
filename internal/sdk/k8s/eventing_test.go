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

package k8s_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"til/config/globals"
	. "til/internal/sdk/k8s"
)

func TestDecodeDestination(t *testing.T) {
	const (
		apiVersion = "test/v0"
		kind       = "Test"
		name       = "test"
	)

	testCases := map[string]struct {
		input       cty.Value
		expect      map[string]interface{}
		expectPanic bool
	}{
		"valid destination": {
			input: NewDestination(apiVersion, kind, name),
			expect: map[string]interface{}{
				"apiVersion": apiVersion,
				"kind":       kind,
				"name":       name,
			},
		},
		"invalid attribute type": {
			input: cty.ObjectVal(map[string]cty.Value{
				"ref": cty.ObjectVal(map[string]cty.Value{
					"apiVersion": cty.Zero,
					"kind":       cty.StringVal(kind),
					"name":       cty.StringVal(name),
				}),
			}),
			expectPanic: true,
		},
		"incomplete destination": {
			input: cty.ObjectVal(map[string]cty.Value{
				"ref": cty.ObjectVal(map[string]cty.Value{
					"apiVersion": cty.StringVal(apiVersion),
					"kind":       cty.StringVal(kind),
				}),
			}),
			expectPanic: true,
		},
		"null": {
			input:       cty.NullVal(DestinationCty),
			expectPanic: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			defer handlePanic(t, tc.expectPanic)

			out := DecodeDestination(tc.input)

			if diff := cmp.Diff(tc.expect, out); diff != "" {
				t.Error("Unexpected diff: (-:expect, +:got)", diff)
			}
		})
	}
}

func TestNewTrigger(t *testing.T) {
	const (
		name   = "test"
		broker = "my-bridge"

		sbAPI  = "test/v0"
		sbKind = "Test"
		sbName = "myapp"

		eventtype       = "ticket.v1"
		ticketid  int64 = 42
		urgent          = true
	)

	subscriber := NewDestination(sbAPI, sbKind, sbName)

	trgg := NewTrigger(name, broker, subscriber,
		Filter(map[string]interface{}{
			"type":     eventtype,
			"ticketid": ticketid,
			"urgent":   urgent,
		}),
	)

	expectTrgg := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": APIEventing,
			"kind":       "Trigger",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"broker": broker,
				"filter": map[string]interface{}{
					"attributes": map[string]interface{}{
						"type":     eventtype,
						"ticketid": ticketid,
						"urgent":   urgent,
					},
				},
				"subscriber": map[string]interface{}{
					"ref": map[string]interface{}{
						"apiVersion": sbAPI,
						"kind":       sbKind,
						"name":       sbName,
					},
				},
			},
		},
	}

	if d := cmp.Diff(expectTrgg, trgg); d != "" {
		t.Errorf("Unexpected diff: (-:expect, +:got) %s", d)
	}
}

func TestNewSubscription(t *testing.T) {
	const (
		name    = "test"
		channel = "my-channel"

		sbAPI  = "test/v0"
		sbKind = "TestSubscriber"
		sbName = "event-dest"

		replAPI  = "test/v0"
		replKind = "TestReply"
		replName = "reply-dest"
	)

	subscriber := NewDestination(sbAPI, sbKind, sbName)
	reply := NewDestination(replAPI, replKind, replName)

	sbsp := NewSubscription(name, channel, subscriber,
		ReplyDest(reply),
	)

	expectSbsp := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": APIMessaging,
			"kind":       "Subscription",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"channel": map[string]interface{}{
					"apiVersion": APIMessaging,
					"kind":       "Channel",
					"name":       channel,
				},
				"subscriber": map[string]interface{}{
					"ref": map[string]interface{}{
						"apiVersion": sbAPI,
						"kind":       sbKind,
						"name":       sbName,
					},
				},
				"reply": map[string]interface{}{
					"ref": map[string]interface{}{
						"apiVersion": replAPI,
						"kind":       replKind,
						"name":       replName,
					},
				},
			},
		},
	}

	if d := cmp.Diff(expectSbsp, sbsp); d != "" {
		t.Errorf("Unexpected diff: (-:expect, +:got) %s", d)
	}
}

func TestMaybeAppendChannel(t *testing.T) {
	const (
		name = "test"

		dstAPI  = "test/v0"
		dstKind = "TestSubscriber"
		dstName = "event-dest"

		dlsAPI  = "test/v0"
		dlsKind = "TestDeadLetter"
		dlsName = "deadletter-dest"
	)

	eventDst := NewDestination(dstAPI, dstKind, dstName)
	deadletterDst := NewDestination(dlsAPI, dlsKind, dlsName)

	commonAssertions := func(t *testing.T, manifests []interface{}, eventDst cty.Value) {
		if l := len(manifests); l != 2 {
			t.Fatalf("Expected 2 manifests: a Channel and a Subscription. Got %d: %s", l, prettyPrintManifests(manifests))
		}

		ch := manifests[0].(*unstructured.Unstructured)
		sbs := manifests[1].(*unstructured.Unstructured)

		if k := ch.GetKind(); k != "Channel" {
			t.Error("Expected Channel kind, got", k)
		}
		if k := sbs.GetKind(); k != "Subscription" {
			t.Error("Expected Subscription kind, got", k)
		}

		// concatenation of 'name' and 'dstName'
		const expectConcatName = "test-event-dest"

		for _, n := range []string{ch.GetName(), sbs.GetName()} {
			if n != expectConcatName {
				t.Error("Expected name to be the concatenation of two names (<sender>-<dest>), got", n)
			}
		}

		sbsName, found, err := unstructured.NestedString(sbs.Object, "spec", "channel", "name")
		if !found || err != nil {
			t.Error("Channel incorrectly set in Subscription spec:", sbs)
		} else if sbsName != expectConcatName {
			t.Error("Incorrect Channel configured in Subscription spec:", sbsName)
		}

		if eventDstKind := eventDst.GetAttr("ref").GetAttr("kind").AsString(); eventDstKind != "Channel" {
			t.Error("Expected new event destination to match the Channel, got", eventDstKind)
		}
	}

	t.Run("delivery settings not set", func(t *testing.T) {
		var manifests []interface{}

		glb := globalsAccessorFunc(func() *globals.Delivery {
			return nil
		})

		manifests, newEventDst := MaybeAppendChannel(name, manifests, eventDst, glb)

		if l := len(manifests); l != 0 {
			t.Errorf("Expected no manifest, got %d: %s", l, prettyPrintManifests(manifests))
		}
		if !newEventDst.RawEquals(eventDst) {
			t.Error("Expected event destination to be unchanged, got", newEventDst)
		}
	})

	t.Run("with empty delivery settings", func(t *testing.T) {
		var manifests []interface{}

		glb := globalsAccessorFunc(func() *globals.Delivery {
			return &globals.Delivery{}
		})

		manifests, eventDst := MaybeAppendChannel(name, manifests, eventDst, glb)

		commonAssertions(t, manifests, eventDst)

		sbs := manifests[1].(*unstructured.Unstructured)

		delivery, found, err := unstructured.NestedMap(sbs.Object, "spec", "delivery")
		if err != nil {
			t.Fatal("Error reading Subscription spec:", err)
		} else if found {
			t.Fatal("Expected no delivery attribute in the spec, got:", delivery)
		}
	})

	t.Run("with non-empty delivery settings", func(t *testing.T) {
		var manifests []interface{}

		glb := globalsAccessorFunc(func() *globals.Delivery {
			three := int64(3)
			return &globals.Delivery{
				Retries:        &three,
				DeadLetterSink: deadletterDst,
			}
		})

		manifests, eventDst := MaybeAppendChannel(name, manifests, eventDst, glb)

		commonAssertions(t, manifests, eventDst)

		sbs := manifests[1].(*unstructured.Unstructured)

		delivery, found, err := unstructured.NestedMap(sbs.Object, "spec", "delivery")
		if err != nil {
			t.Fatal("Error reading Subscription spec:", err)
		} else if !found {
			t.Fatal("Expected a delivery attribute in the spec")
		}

		retries, found, err := unstructured.NestedInt64(delivery, "retry")
		if err != nil {
			t.Fatal("Error reading Subscription spec:", err)
		} else if !found {
			t.Error("Expected a retry attribute in the delivery spec")
		} else if retries != 3 {
			t.Error("Expected 3 retries, got", retries)
		}

		dls, found, err := unstructured.NestedMap(delivery, "deadLetterSink", "ref")
		if err != nil {
			t.Fatal("Error reading Subscription spec:", err)
		} else if !found {
			t.Error("Expected a deadLetterSink attribute in the delivery spec")
		} else if name := dls["name"]; name != dlsName {
			t.Error("Unexpected deadLetterSink name:", name)
		}
	})
}

type globalsAccessorFunc func() *globals.Delivery

func (a globalsAccessorFunc) Delivery() *globals.Delivery {
	return a()
}

func prettyPrintManifests(manifests []interface{}) string {
	var b strings.Builder

	for _, m := range manifests {
		b.WriteByte('\n')
		b.WriteString(fmt.Sprint(m))
	}

	return b.String()
}
