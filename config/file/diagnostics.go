package file

import (
	"strconv"

	"github.com/hashicorp/hcl/v2"
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
func duplicateBlockDiagnostic(addr string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Duplicate block",
		Detail:   "Found a duplicate block for the identifier " + strconv.Quote(addr),
		Subject:  subj.Ptr(),
	}
}
