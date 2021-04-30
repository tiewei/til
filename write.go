package main

import (
	"encoding/json"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type manifestsWriterFunc func(out io.Writer, manifests []interface{}) error

// writeManifestsJSON marshals the given manifests to a Kubernetes object of
// kind "List" in JSON format, and writes the result to out.
func writeManifestsJSON(out io.Writer, manifests []interface{}) error {
	// NOTE(antoineco): We assume for the time being that all generated
	// manifests are unstructured.Unstructured objects. This might change
	// in the future. See translation.Translatable.
	list := &unstructured.UnstructuredList{}
	list.SetAPIVersion("v1")
	list.SetKind("List")

	for _, m := range manifests {
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

// writeManifestsYAML marshals the given manifests to a sequence of YAML
// documents separated by a "---" marker, and writes the result to out.
func writeManifestsYAML(out io.Writer, manifests []interface{}) error {
	for i, m := range manifests {
		if i > 0 {
			fmt.Fprintln(out, "---")
		}

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
