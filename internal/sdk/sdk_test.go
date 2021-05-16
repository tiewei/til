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

package sdk_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	. "bridgedl/internal/sdk"
)

func TestDecodeStringSlice(t *testing.T) {
	testCases := map[string]struct {
		input       cty.Value
		expect      []interface{}
		expectPanic bool
	}{
		"tuple type with string values": {
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("val1"),
				cty.StringVal("val2"),
			}),
			expect: []interface{}{"val1", "val2"},
		},
		"tuple type with non-string values": {
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("val1"),
				cty.NumberIntVal(1),
			}),
			expectPanic: true,
		},
		"object type with string values": {
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("val1"),
				"key2": cty.StringVal("val2"),
			}),
			expect: []interface{}{"val1", "val2"},
		},
		"object type with non-string values": {
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("val1"),
				"key2": cty.NumberIntVal(1),
			}),
			expectPanic: true,
		},
		"null collection type": {
			input:       cty.NullVal(cty.Set(cty.String)),
			expectPanic: true,
		},
		"non-collection type": {
			input:       cty.StringVal("val1"),
			expectPanic: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			defer handlePanic(t, tc.expectPanic)

			out := DecodeStringSlice(tc.input)

			if diff := cmp.Diff(tc.expect, out); diff != "" {
				t.Error("Unexpected diff: (-:expect, +:got)", diff)
			}
		})
	}
}

func TestDecodeStringMap(t *testing.T) {
	testCases := map[string]struct {
		input       cty.Value
		expect      map[string]interface{}
		expectPanic bool
	}{
		"object type with string values": {
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("val1"),
				"key2": cty.StringVal("val2"),
			}),
			expect: map[string]interface{}{"key1": "val1", "key2": "val2"},
		},
		"object type with null string values": {
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.NullVal(cty.String),
				"key2": cty.StringVal("val2"),
			}),
			expect: map[string]interface{}{"key2": "val2"},
		},
		"object type with non-string values": {
			input: cty.ObjectVal(map[string]cty.Value{
				"key1": cty.StringVal("val1"),
				"key2": cty.NumberIntVal(1),
			}),
			expectPanic: true,
		},
		"map type with string values": {
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.StringVal("val1"),
				"key2": cty.StringVal("val2"),
			}),
			expect: map[string]interface{}{"key1": "val1", "key2": "val2"},
		},
		"map type with null string values": {
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NullVal(cty.String),
				"key2": cty.StringVal("val2"),
			}),
			expect: map[string]interface{}{"key2": "val2"},
		},
		"map type with non-string values": {
			input: cty.MapVal(map[string]cty.Value{
				"key1": cty.NumberIntVal(1),
				"key2": cty.NumberIntVal(2),
			}),
			expectPanic: true,
		},
		"tuple type": {
			input: cty.TupleVal([]cty.Value{
				cty.StringVal("val1"),
				cty.StringVal("val2"),
			}),
			expectPanic: true,
		},
		"null collection type": {
			input:       cty.NullVal(cty.Set(cty.String)),
			expectPanic: true,
		},
		"non-collection type": {
			input:       cty.StringVal("val1"),
			expectPanic: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			defer handlePanic(t, tc.expectPanic)

			out := DecodeStringMap(tc.input)

			if diff := cmp.Diff(tc.expect, out); diff != "" {
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
