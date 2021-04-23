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
		"not a number": {
			in:        cty.False,
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			diags := IsInt()(tc.in)

			if tc.expectErr && diags == nil {
				t.Error("Expected validation to fail")
			}
			if !tc.expectErr && diags != nil {
				t.Error("Expected validation to pass. Got diagnostic:", diags)
			}
		})
	}
}
