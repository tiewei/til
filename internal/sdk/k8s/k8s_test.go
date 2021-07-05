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

	. "til/internal/sdk/k8s"
)

func TestNewObject(t *testing.T) {
	const (
		apiVersion = "test/v0"
		kind       = "Test"
	)

	testCases := map[string]struct {
		name          string
		expectObjBody map[string]interface{}
		expectPanic   bool
	}{
		"valid Kubernetes object name": {
			name: "bridge.my-object-1234",
			expectObjBody: map[string]interface{}{
				"apiVersion": apiVersion,
				"kind":       kind,
				"metadata": map[string]interface{}{
					"name": "bridge.my-object-1234",
				},
			},
		},
		"invalid Kubernetes object name": {
			name:        "bridge_my-object-1234",
			expectPanic: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			defer handlePanic(t, tc.expectPanic)

			out := NewObject(apiVersion, kind, tc.name)

			got := out.Unstructured().Object
			if diff := cmp.Diff(tc.expectObjBody, got); diff != "" {
				t.Error("Unexpected diff: (-:expect, +:got)", diff)
			}
		})
	}
}

func handlePanic(t *testing.T, expectPanic bool) {
	t.Helper()

	if r := recover(); r != nil && !expectPanic {
		t.Fatal("Unexpected panic:", r)
	}
}
