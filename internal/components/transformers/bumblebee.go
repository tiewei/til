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

package transformers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config/globals"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

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

	/*
		Example of value decoded from the spec above, for a "context"
		and a "data" block, both containing two nested "operation" blocks
		with one or more "path" sub-blocks:

		v: (map[string]interface {}) (len=2) {
		 "context": ([]interface {}) (len=2) {
		  (map[string]interface {}) (len=2) {
		   "operation": (string) "store",
		   "path": ([]interface {}) (len=2) {
		    (map[string]interface {}) (len=2) {
		     "key": (string) "$id",
		     "value": (string) "id"
		    },
		    (map[string]interface {}) (len=2) {
		     "key": (string) "$type",
		     "value": (string) "type"
		    }
		   }
		  },
		  (map[string]interface {}) (len=2) {
		   "operation": (string) "add",
		   "path": ([]interface {}) (len=1) {
		    (map[string]interface {}) (len=2) {
		     "key": (string) "idtype",
		     "value": (string) "${id}-${type}"
		    }
		   }
		  }
		 },
		 "data": ([]interface {}) (len=2) { ... }
		}
	*/
}

// Manifests implements translation.Translatable.
func (*Bumblebee) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	t := k8s.NewObject(k8s.APIFlow, "Transformation", name)

	context := parseBumblebeeOperations(config.GetAttr("context").AsValueSlice())
	t.SetNestedSlice(context, "spec", "context")

	data := parseBumblebeeOperations(config.GetAttr("data").AsValueSlice())
	t.SetNestedSlice(data, "spec", "data")

	sink := k8s.DecodeDestination(eventDst)
	t.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, t.Unstructured())
}

// Address implements translation.Addressable.
func (*Bumblebee) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination(k8s.APIFlow, "Transformation", k8s.RFC1123Name(id))
}

func parseBumblebeeOperations(operationVals []cty.Value) []interface{} {
	operations := make([]interface{}, 0, len(operationVals))

	for _, operationVal := range operationVals {
		operationValMap := operationVal.AsValueMap()

		operation := operationValMap["operation"].AsString()

		pathVals := operationValMap["path"].AsValueSlice()
		paths := make([]interface{}, 0, len(pathVals))
		for _, pathVal := range pathVals {
			path := make(map[string]interface{})

			pathValMap := pathVal.AsValueMap()

			if key := pathValMap["key"]; !key.IsNull() {
				path["key"] = key.AsString()
			}
			if value := pathValMap["value"]; !value.IsNull() {
				path["value"] = value.AsString()
			}

			paths = append(paths, path)
		}

		operations = append(operations, map[string]interface{}{
			"operation": operation,
			"paths":     paths,
		})
	}

	return operations
}
