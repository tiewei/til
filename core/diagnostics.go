package core

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// noTranslatorDiagnostic returns a hcl.Diagnostic which indicates that a block
// translator can not be acquired for a given component type.
func noTranslatorDiagnostic(blockType, componentType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "No translator",
		Detail:   fmt.Sprintf("Could not find a block translator for a %s of type %q", blockType, componentType),
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
