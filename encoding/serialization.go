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

package encoding

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"sigs.k8s.io/yaml"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Serializer can serialize Kubernetes API objects into various formats.
type Serializer struct {
	// A HCL identifier which represents the Bridge that "owns" the
	// Kubernetes objects processed by this Serializer.
	BridgeIdentifier string
}

// NewSerializer returns a Serializer initialized with the given Bridge identifier.
func NewSerializer(brgID string) *Serializer {
	return &Serializer{
		BridgeIdentifier: brgID,
	}
}

// ManifestsWriterFunc is the signature of a function which serializes
// Kubernetes manifests and writes the result to the given io.Writer.
//
// NOTE(antoineco): We assume for the time being that all processed manifests
// are unstructured.Unstructured objects. This might change in the future. See
// translation.Translatable.
type ManifestsWriterFunc func(out io.Writer, manifests []interface{}) error

// WriteManifestsJSON marshals the given manifests to a Kubernetes object of
// kind "List" in JSON format, and writes the result to out.
func (s *Serializer) WriteManifestsJSON(out io.Writer, manifests []interface{}) error {
	list := &unstructured.UnstructuredList{}
	list.SetAPIVersion("v1")
	list.SetKind("List")

	for _, m := range manifests {
		m = injectBridgeLabels(m, s.BridgeIdentifier)
		list.Items = append(list.Items, *m.(*unstructured.Unstructured))
	}

	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling manifests to JSON: %w", err)
	}

	if _, err := fmt.Fprintln(out, string(b)); err != nil {
		return fmt.Errorf("writing generated manifests: %w", err)
	}

	return nil
}

// WriteManifestsYAML marshals the given manifests to a sequence of YAML
// documents separated by a "---" marker, and writes the result to out.
func (s *Serializer) WriteManifestsYAML(out io.Writer, manifests []interface{}) error {
	for i, m := range manifests {
		if i > 0 {
			fmt.Fprintln(out, "---")
		}

		m = injectBridgeLabels(m, s.BridgeIdentifier)

		b, err := yaml.Marshal(m)
		if err != nil {
			return fmt.Errorf("marshaling manifest to YAML: %w", err)
		}

		if _, err := out.Write(b); err != nil {
			return fmt.Errorf("writing generated manifest: %w", err)
		}
	}

	return nil
}

// WriteBridgeJSON marshals the given manifests to a Bridge API object in JSON
// format, and writes the result to out.
func (s *Serializer) WriteBridgeJSON(out io.Writer, manifests []interface{}) error {
	brg := newBridge(s.BridgeIdentifier, manifests)

	b, err := json.MarshalIndent(brg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling Bridge object to JSON: %w", err)
	}

	if _, err := fmt.Fprintln(out, string(b)); err != nil {
		return fmt.Errorf("writing generated Bridge object: %w", err)
	}

	return nil
}

// WriteBridgeYAML marshals the given manifests to a Bridge API object in YAML
// format, and writes the result to out.
func (s *Serializer) WriteBridgeYAML(out io.Writer, manifests []interface{}) error {
	brg := newBridge(s.BridgeIdentifier, manifests)

	b, err := yaml.Marshal(brg)
	if err != nil {
		return fmt.Errorf("marshaling Bridge object to JSON: %w", err)
	}

	if _, err := out.Write(b); err != nil {
		return fmt.Errorf("writing generated Bridge object: %w", err)
	}

	return nil
}

// newBridge returns a Bridge object which components correspond to the given manifests.
func newBridge(brgID string, manifests []interface{}) *unstructured.Unstructured {
	brg := &unstructured.Unstructured{}
	brg.SetAPIVersion("flow.triggermesh.io/v1alpha1")
	brg.SetKind("Bridge")
	brg.SetName(sanitizeBridgeIdentifier(brgID))

	var components []interface{}

	for _, m := range manifests {
		components = append(components, map[string]interface{}{
			"object": m.(*unstructured.Unstructured).Object,
		})
	}
	_ = unstructured.SetNestedSlice(brg.Object, components, "spec", "components")

	return brg
}

// injectBridgeLabels sets a standard set of labels in the metadata of the
// given manifest.
func injectBridgeLabels(manifest interface{}, brgID string) interface{} {
	const lblBridgeIdentifier = "bridges.triggermesh.io/id"

	m := manifest.(*unstructured.Unstructured)

	lbls := m.GetLabels()
	if lbls == nil {
		lbls = make(map[string]string, 1)
	}
	lbls[lblBridgeIdentifier] = brgID

	m.SetLabels(lbls)

	return m
}

// sanitizeBridgeIdentifier converts the given Bridge identifier into a valid
// Kubernetes object name.
//
// The implementation is purposedly naive because the input is assumed to be a
// valid HCL identifiers, which already contain a limited set of characters
// (see hclsyntax.ValidIdentifier).
func sanitizeBridgeIdentifier(id string) string {
	return strings.ToLower(strings.ReplaceAll(id, "_", "-"))
}
