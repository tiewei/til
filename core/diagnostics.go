package core

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// noComponentImplDiagnostic returns a hcl.Diagnostic which indicates that no
// implementation is available for a given component type.
func noComponentImplDiagnostic(cmpCat, cmpType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not implemented",
		Detail:   fmt.Sprintf("No implementation found for a %s of type %q", cmpCat, cmpType),
		Subject:  subj.Ptr(),
	}
}

// noDecodeSpecDiagnostic returns a hcl.Diagnostic which indicates that a spec
// for decoding a HCL body can not be acquired for a given component type.
func noDecodeSpecDiagnostic(cmpCat, cmpType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "No decode spec",
		Detail:   fmt.Sprintf("Could not find a decode spec for a %s of type %q", cmpCat, cmpType),
		Subject:  subj.Ptr(),
	}
}

// unknownReferenceDiagnostic returns a hcl.Diagnostic which indicates that a block
// reference refers to an unknown block.
func unknownReferenceDiagnostic(refAddr string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Reference to unknown block",
		Detail:   fmt.Sprintf("The expression %q doesn't match any known configuration block", refAddr),
		Subject:  subj.Ptr(),
	}
}
