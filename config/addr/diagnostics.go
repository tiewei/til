package addr

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// badRefFormatDiagnostic returns a hcl.Diagnostic which indicates that a block
// reference is not expressed in a correct format.
func badRefFormatDiagnostic(subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid block reference",
		Detail:   "A block reference is expressed as a block type separated from an identifier by a dot.",
		Subject:  subj.Ptr(),
	}
}

// badRefTypeDiagnostic returns a hcl.Diagnostic which indicates that the type
// indicated in a block reference is not supported.
func badRefTypeDiagnostic(typ string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid block reference",
		Detail: fmt.Sprintf("The expression refers to a block type that is unknown or doesn't support "+
			"references: %q. Valid values are %q", typ, referenceableTypes()),
		Subject: subj.Ptr(),
	}
}
