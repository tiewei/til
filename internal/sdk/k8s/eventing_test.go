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
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "bridgedl/internal/sdk/k8s"
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
