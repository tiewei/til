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

package sdk

import "github.com/zclconf/go-cty/cty"

// DecodeStringSlice decodes a cty.Value into a slice of strings, in a format
// that is compatible with the "unstructured" package from k8s.io/apimachinery.
// Panics if the given value is not a non-null collection type containing
// exclusively string values.
func DecodeStringSlice(val cty.Value) []interface{} {
	out := make([]interface{}, 0, val.LengthInt())
	for iter := val.ElementIterator(); iter.Next(); {
		_, v := iter.Element()
		out = append(out, v.AsString())
	}

	return out
}

// DecodeStringMap decodes a cty.Value into a map of string elements, in a
// format that is compatible with the "unstructured" package from k8s.io/apimachinery.
// Panics if the given value is not a non-null collection type containing
// exclusively string values.
func DecodeStringMap(val cty.Value) map[string]interface{} {
	out := make(map[string]interface{}, val.LengthInt())
	for iter := val.ElementIterator(); iter.Next(); {
		if k, v := iter.Element(); !v.IsNull() {
			out[k.AsString()] = v.AsString()
		}
	}

	return out
}
