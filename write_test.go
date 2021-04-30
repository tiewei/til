package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestWriteManifests(t *testing.T) {
	testCases := map[string] /*output*/ struct {
		writerFunc   manifestsWriterFunc
		expectOutput string
	}{
		"JSON List-manifest": {
			writerFunc:   writeManifestsJSON,
			expectOutput: expectOutputJSON,
		},
		"YAML documents": {
			writerFunc:   writeManifestsYAML,
			expectOutput: expectOutputYAML,
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
	`        "name": "object-1"` + "\n" +
	`      }` + "\n" +
	`    },` + "\n" +
	`    {` + "\n" +
	`      "apiVersion": "fake/v0",` + "\n" +
	`      "kind": "Bar",` + "\n" +
	`      "metadata": {` + "\n" +
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
	"  name: object-1\n" +
	"---\n" +
	"apiVersion: fake/v0\n" +
	"kind: Bar\n" +
	"metadata:\n" +
	"  name: object-2\n"

func newUnstructured(apiVersion, kind, name string) *unstructured.Unstructured {
	o := &unstructured.Unstructured{}
	o.SetAPIVersion(apiVersion)
	o.SetKind(kind)
	o.SetName(name)
	return o
}
