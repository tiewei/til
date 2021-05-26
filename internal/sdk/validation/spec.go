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
	"math/big"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// ValidateSpecFunc is the signature of a validation function used in hcldec.ValidateSpec
// to validate a given HCL spec.
type ValidateSpecFunc func(cty.Value) hcl.Diagnostics

// IsInt is a ValidateSpecFunc which asserts that the given value is an integer.
func IsInt(v cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if v.IsNull() {
		return diags
	}
	if v.Type() != cty.Number || !isInt64(v.AsBigFloat()) {
		diags = diags.Append(wrongTypeDiagnostic(v, "integer"))
	}

	return diags
}

// isInt64 returns whether the given cty.Number value can be represented as an int64.
func isInt64(v *big.Float) bool {
	bigInt, accuracy := v.Int(nil)
	if accuracy != big.Exact {
		return false
	}

	return bigInt.IsInt64()
}

// IsCEContextAttribute is a ValidateSpecFunc which asserts that the given
// string follows the naming conventions for CloudEvent context attributes.
// https://github.com/cloudevents/spec/blob/v1.0.1/spec.md#attribute-naming-convention
func IsCEContextAttribute(v cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if v.Type() != cty.String {
		diags = diags.Append(wrongTypeDiagnostic(v, "string"))
		return diags
	}

	if s := v.AsString(); !isLowercaseAlphanum(s) {
		diags = diags.Append(invalidCEContextAttrDiagnostic(s))
	}

	return diags
}

// isLowercaseAlphanum returns whether the given string value contains only
// lower-case letters or digits.
func isLowercaseAlphanum(v string) bool {
	// operate on bytes instead of runes, since all alphanumeric characters
	// are represented in a single byte
	for i := 0; i < len(v); i++ {
		if ch := v[i]; !(('a' <= ch && ch <= 'z') || ('0' <= ch && ch <= '9')) {
			return false
		}
	}

	return true
}

// ContainsCEContextAttributes is a ValidateSpecFunc which asserts that all
// string keys in the given collection follow the naming conventions for
// CloudEvent context attributes.
// https://github.com/cloudevents/spec/blob/v1.0.1/spec.md#attribute-naming-convention
func ContainsCEContextAttributes(val cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if !val.CanIterateElements() {
		diags = diags.Append(wrongTypeDiagnostic(val, "collection"))
		return diags
	}
	if val.IsNull() {
		return diags
	}

	for iter := val.ElementIterator(); iter.Next(); {
		var v cty.Value

		switch typ := val.Type(); {
		case typ.IsListType(), typ.IsSetType(), typ.IsTupleType():
			_, v = iter.Element()
		case typ.IsMapType(), typ.IsObjectType():
			v, _ = iter.Element()
		}

		diags = diags.Extend(IsCEContextAttribute(v))
	}

	return diags
}
