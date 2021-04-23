package validation

import (
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
