package validation

import (
	"math/big"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// ValidateSpecFunc is the signature of a validation function used in hcldec.ValidateSpec
// to validate a given HCL spec.
type ValidateSpecFunc func(cty.Value) hcl.Diagnostics

// IsInt returns a validation function which asserts that the given value is an integer.
func IsInt() ValidateSpecFunc {
	return func(v cty.Value) hcl.Diagnostics {
		var diags hcl.Diagnostics

		if !isInt64(v) {
			diags = diags.Append(wrongTypeDiagnostic(v, "integer"))
		}

		return diags
	}
}

// isInt64 returns whether the given cty.Number value can be represented as an int64.
func isInt64(v cty.Value) bool {
	if v.Type() != cty.Number {
		return false
	}

	bigInt, accuracy := v.AsBigFloat().Int(nil)
	if accuracy != big.Exact {
		return false
	}

	return bigInt.IsInt64()
}
