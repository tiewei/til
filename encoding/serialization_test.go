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

package encoding_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "til/encoding"
)

func TestWriteManifests(t *testing.T) {
	const testBrgID = "Test_Bridge"

	s := NewSerializer(testBrgID)

	testCases := map[string] /*output*/ struct {
		writerFunc   ManifestsWriterFunc
		expectOutput string
	}{
		"JSON List-manifest": {
			writerFunc:   s.WriteManifestsJSON,
			expectOutput: expectOutputJSON,
		},
		"YAML documents": {
			writerFunc:   s.WriteManifestsYAML,
			expectOutput: expectOutputYAML,
		},
		"JSON Bridge": {
			writerFunc:   s.WriteBridgeJSON,
			expectOutput: expectOutputBridgeJSON,
		},
		"YAML Bridge": {
			writerFunc:   s.WriteBridgeYAML,
			expectOutput: expectOutputBridgeYAML,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			var stdout strings.Builder

			err := tc.writerFunc(&stdout, newTestManifests())
			if err != nil {
				t.Fatal("Returned an error:", err)
			}

			output := stdout.String()

			if diff := cmp.Diff(tc.expectOutput, output); diff != "" {
				t.Error("Unexpected diff: (-:expect, +:got)", diff)
			}
		})
	}
}

func newTestManifests() []interface{} {
	return []interface{}{
		newUnstructured("fake/v0", "Foo", "object-1"),
		newUnstructured("fake/v0", "Bar", "object-2"),
	}
}

const expectOutputJSON = "" +
	`{` + "\n" +
	`  "apiVersion": "v1",` + "\n" +
	`  "items": [` + "\n" +
	`    {` + "\n" +
	`      "apiVersion": "fake/v0",` + "\n" +
	`      "kind": "Foo",` + "\n" +
	`      "metadata": {` + "\n" +
	`        "labels": {` + "\n" +
	`          "bridges.triggermesh.io/id": "Test_Bridge"` + "\n" +
	`        },` + "\n" +
	`        "name": "object-1"` + "\n" +
	`      }` + "\n" +
	`    },` + "\n" +
	`    {` + "\n" +
	`      "apiVersion": "fake/v0",` + "\n" +
	`      "kind": "Bar",` + "\n" +
	`      "metadata": {` + "\n" +
	`        "labels": {` + "\n" +
	`          "bridges.triggermesh.io/id": "Test_Bridge"` + "\n" +
	`        },` + "\n" +
	`        "name": "object-2"` + "\n" +
	`      }` + "\n" +
	`    }` + "\n" +
	`  ],` + "\n" +
	`  "kind": "List"` + "\n" +
	`}` + "\n"

const expectOutputYAML = "" +
	"apiVersion: fake/v0\n" +
	"kind: Foo\n" +
	"metadata:\n" +
	"  labels:\n" +
	"    bridges.triggermesh.io/id: Test_Bridge\n" +
	"  name: object-1\n" +
	"---\n" +
	"apiVersion: fake/v0\n" +
	"kind: Bar\n" +
	"metadata:\n" +
	"  labels:\n" +
	"    bridges.triggermesh.io/id: Test_Bridge\n" +
	"  name: object-2\n"

const expectOutputBridgeJSON = "" +
	`{` + "\n" +
	`  "apiVersion": "flow.triggermesh.io/v1alpha1",` + "\n" +
	`  "kind": "Bridge",` + "\n" +
	`  "metadata": {` + "\n" +
	`    "name": "test-bridge"` + "\n" +
	`  },` + "\n" +
	`  "spec": {` + "\n" +
	`    "components": [` + "\n" +
	`      {` + "\n" +
	`        "object": {` + "\n" +
	`          "apiVersion": "fake/v0",` + "\n" +
	`          "kind": "Foo",` + "\n" +
	`          "metadata": {` + "\n" +
	`            "name": "object-1"` + "\n" +
	`          }` + "\n" +
	`        }` + "\n" +
	`      },` + "\n" +
	`      {` + "\n" +
	`        "object": {` + "\n" +
	`          "apiVersion": "fake/v0",` + "\n" +
	`          "kind": "Bar",` + "\n" +
	`          "metadata": {` + "\n" +
	`            "name": "object-2"` + "\n" +
	`          }` + "\n" +
	`        }` + "\n" +
	`      }` + "\n" +
	`    ]` + "\n" +
	`  }` + "\n" +
	`}` + "\n"

const expectOutputBridgeYAML = "" +
	"apiVersion: flow.triggermesh.io/v1alpha1\n" +
	"kind: Bridge\n" +
	"metadata:\n" +
	"  name: test-bridge\n" +
	"spec:\n" +
	"  components:\n" +
	"  - object:\n" +
	"      apiVersion: fake/v0\n" +
	"      kind: Foo\n" +
	"      metadata:\n" +
	"        name: object-1\n" +
	"  - object:\n" +
	"      apiVersion: fake/v0\n" +
	"      kind: Bar\n" +
	"      metadata:\n" +
	"        name: object-2\n"

func newUnstructured(apiVersion, kind, name string) *unstructured.Unstructured {
	o := &unstructured.Unstructured{}
	o.SetAPIVersion(apiVersion)
	o.SetKind(kind)
	o.SetName(name)
	return o
}
