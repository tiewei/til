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

package k8s

import "github.com/zclconf/go-cty/cty"

// DestinationCty is a non-primitive cty.Type that represents a Knative "duck" Destination.
var DestinationCty = cty.Object(map[string]cty.Type{
	"ref": cty.Object(map[string]cty.Type{
		"apiVersion": cty.String,
		"kind":       cty.String,
		"name":       cty.String,
	}),
	"uri": cty.String,
})

// NewDestination returns a new Knative "duck" Destination as a cty.Value which
// satisfies the DestinationCty type.
func NewDestination(apiVersion, kind, name string) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"ref": cty.ObjectVal(map[string]cty.Value{
			"apiVersion": cty.StringVal(apiVersion),
			"kind":       cty.StringVal(kind),
			"name":       cty.StringVal(name),
		}),
		"uri": cty.NullVal(cty.String),
	})
}

// IsDestination verifies that the given cty.Value conforms to the
// DestinationCty type.
func IsDestination(v cty.Value) bool {
	return v.Type().Equals(DestinationCty)
}
