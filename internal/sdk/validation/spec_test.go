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

package validation

import (
	"math"
	"strconv"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestIsInt(t *testing.T) {
	testCases := map[string]struct {
		in        cty.Value
		expectErr bool
	}{
		"zero": {
			in:        cty.Zero,
			expectErr: false,
		},
		"unsigned int": {
			in:        cty.NumberUIntVal(42),
			expectErr: false,
		},
		"negative int": {
			in:        cty.NumberIntVal(-42),
			expectErr: false,
		},
		"float with null fractional part": {
			in:        cty.NumberFloatVal(42.0),
			expectErr: false,
		},
		"float with non-null fractional part": {
			in:        cty.NumberFloatVal(42.1),
			expectErr: true,
		},
		"integer overflow": {
			in:        cty.MustParseNumberVal(strconv.FormatUint(math.MaxInt64+1, 10)),
			expectErr: true,
		},
		"null value": {
			in:        cty.NullVal(cty.Number),
			expectErr: false,
		},
		"not a number": {
			in:        cty.False,
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			diags := IsInt(tc.in)

			if tc.expectErr && diags == nil {
				t.Error("Expected validation to fail")
			}
			if !tc.expectErr && diags != nil {
				t.Error("Expected validation to pass. Got diagnostic:", diags)
			}
		})
	}
}

func TestIsCEContextAttribute(t *testing.T) {
	testCases := map[string]struct {
		in        cty.Value
		expectErr bool
	}{
		"only lowercase alphanum chars": {
			in:        cty.StringVal("abc123"),
			expectErr: false,
		},
		"contains uppercase chars": {
			in:        cty.StringVal("aBc123"),
			expectErr: true,
		},
		"contains international chars": {
			in:        cty.StringVal("àbc123"),
			expectErr: true,
		},
		"contains spaces": {
			in:        cty.StringVal("abc 123"),
			expectErr: true,
		},
		"not a string": {
			in:        cty.False,
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			diags := IsCEContextAttribute(tc.in)

			if tc.expectErr && diags == nil {
				t.Error("Expected validation to fail")
			}
			if !tc.expectErr && diags != nil {
				t.Error("Expected validation to pass. Got diagnostic:", diags)
			}
		})
	}
}

func TestContainsCEContextAttributes(t *testing.T) {
	testCases := map[string]struct {
		in        cty.Value
		expectErr bool
	}{
		"tuple type with valid attributes": {
			in: cty.TupleVal([]cty.Value{
				cty.StringVal("abc123"),
				cty.StringVal("def456"),
			}),
			expectErr: false,
		},
		"tuple type with invalid attributes": {
			in: cty.TupleVal([]cty.Value{
				cty.StringVal("abc123"),
				cty.StringVal("Oops"),
			}),
			expectErr: true,
		},
		"tuple type with non-string attributes": {
			in: cty.TupleVal([]cty.Value{
				cty.StringVal("abc123"),
				cty.NumberIntVal(1),
			}),
			expectErr: true,
		},
		"object type with valid attributes": {
			in: cty.ObjectVal(map[string]cty.Value{
				"abc123": cty.StringVal("val1"),
				"def456": cty.NumberIntVal(1),
			}),
			expectErr: false,
		},
		"object type with invalid attributes": {
			in: cty.ObjectVal(map[string]cty.Value{
				"abc123": cty.StringVal("val1"),
				"Oops":   cty.NumberIntVal(1),
			}),
			expectErr: true,
		},
		"null collection": {
			in:        cty.NullVal(cty.Map(cty.String)),
			expectErr: false,
		},
		"not a collection": {
			in:        cty.False,
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			diags := ContainsCEContextAttributes(tc.in)

			if tc.expectErr && diags == nil {
				t.Error("Expected validation to fail")
			}
			if !tc.expectErr && diags != nil {
				t.Error("Expected validation to pass. Got diagnostic:", diags)
			}
		})
	}
}
