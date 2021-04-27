package file

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
)

// badIdentifierDiagnostic returns a hcl.Diagnostic which indicates that the
// format of an identifier is not valid for a native syntax expression.
func badIdentifierDiagnostic(subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid identifier",
		Detail: "An identifier must start with a letter or underscore and may contain only " +
			"letters, digits, underscores, and dashes.",
		Subject: subj.Ptr(),
	}
}

// duplicateBlockDiagnostic returns a hcl.Diagnostic which indicates that a
// duplicate block definition was found.
func duplicateBlockDiagnostic(cat config.ComponentCategory, identifier string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Duplicate block",
		Detail:   fmt.Sprintf("Found a duplicate %q block with the identifier %q.", cat, identifier),
		Subject:  subj.Ptr(),
	}
}

// tooManyBridgeBlocksDiagnostic returns a hcl.Diagnostic which indicates that
// more than one "bridge" block was defined.
func tooManyBridgeBlocksDiagnostic(subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Redefined global config",
		Detail:   `A Bridge description must contain at most one "bridge" block.`,
		Subject:  subj.Ptr(),
	}
}
