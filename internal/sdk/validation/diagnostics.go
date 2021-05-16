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
	"strconv"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

const diagSummaryValidation = "Failed validation"

// wrongTypeDiagnostic returns a validation diagnostic which indicates that the
// given value doesn't have the expected type.
func wrongTypeDiagnostic(v cty.Value, expectType string) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  diagSummaryValidation,
		Detail:   "The provided value is not a " + expectType + ". Type: " + v.Type().GoString(),
	}
}

// invalidCEContextAttrDiagnostic returns a validation diagnostic which
// indicates that the given value is not a valid CloudEvent context attribute.
func invalidCEContextAttrDiagnostic(v string) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  diagSummaryValidation,
		Detail: "The provided value " + strconv.Quote(v) + " is not a valid CloudEvent context attribute. " +
			"CloudEvents attribute names must consist of lower-case letters ('a' to 'z') or digits " +
			"('0' to '9') from the ASCII character set.",
	}
}
