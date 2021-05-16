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
