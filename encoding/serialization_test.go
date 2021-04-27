package encoding_test

import . "bridgedl/encoding"

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
