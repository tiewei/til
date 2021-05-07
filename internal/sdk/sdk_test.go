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

func handlePanic(t *testing.T, expectPanic bool) {
	t.Helper()

	if r := recover(); r != nil && !expectPanic {
		t.Fatal("Unexpected panic:", r)
	}
}
