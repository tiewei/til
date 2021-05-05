package transformers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
	"bridgedl/translation"
)

// TODO(antoineco): this is a mock implementation. It exists only to allow
// testing the generator manually using the config.brg.hcl file at the root of
// the repo while this code is still at the stage of prototype.
type Bumblebee struct{}

var (
	_ translation.Decodable    = (*Bumblebee)(nil)
	_ translation.Translatable = (*Bumblebee)(nil)
	_ translation.Addressable  = (*Bumblebee)(nil)
)

// Spec implements translation.Decodable.
func (*Bumblebee) Spec() hcldec.Spec {
	nestedSpec := &hcldec.BlockListSpec{
		TypeName: "operation",
		Nested: &hcldec.ObjectSpec{
			"operation": &hcldec.BlockLabelSpec{
				Index: 0,
				Name:  "operation",
			},
			"path": &hcldec.BlockListSpec{
				TypeName: "path",
				MinItems: 1,
				Nested: &hcldec.ObjectSpec{
					"key": &hcldec.AttrSpec{
						Name: "key",
						Type: cty.String,
					},
					"value": &hcldec.AttrSpec{
						Name: "value",
						Type: cty.String,
					},
				},
			},
		},
	}

	return &hcldec.ObjectSpec{
		"context": &hcldec.BlockSpec{
			TypeName: "context",
			Nested:   nestedSpec,
		},
		"data": &hcldec.BlockSpec{
			TypeName: "data",
			Nested:   nestedSpec,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Bumblebee) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("flow.triggermesh.io/v1alpha1")
	s.SetKind("Transformation")
	s.SetName(k8s.RFC1123Name(id))

	context := parseBumblebeeOperations(config.GetAttr("context").AsValueSlice())
	_ = unstructured.SetNestedSlice(s.Object, context, "spec", "context")

	data := parseBumblebeeOperations(config.GetAttr("data").AsValueSlice())
	_ = unstructured.SetNestedSlice(s.Object, data, "spec", "data")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}

func parseBumblebeeOperations(values []cty.Value) []interface{} {
	var operations []interface{}
	for _, operation := range values {
		var paths []interface{}
		for _, path := range operation.AsValueMap()["path"].AsValueSlice() {
			p := make(map[string]interface{})
			if key := path.AsValueMap()["key"]; !key.IsNull() {
				p["key"] = key.AsString()
			}
			if value := path.AsValueMap()["value"]; !value.IsNull() {
				p["value"] = value.AsString()
			}
			paths = append(paths, p)
		}
		operations = append(operations, map[string]interface{}{
			"operation": operation.AsValueMap()["operation"].AsString(),
			"paths":     paths,
		})
	}
	return operations
}

// Address implements translation.Addressable.
func (*Bumblebee) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("flow.triggermesh.io/v1alpha1", "Transformation", k8s.RFC1123Name(id))
}
