package config

import "github.com/hashicorp/hcl/v2"

// badIdentifierDiagnostic returns a hcl.Diagnostic which indicates that the
// format of an identifier is not valid for a native syntax expression.
func badIdentifierDiagnostic(subj *hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid identifier",
		Detail: "An identifier must start with a letter or underscore and may contain only " +
			"letters, digits, underscores, and dashes.",
		Subject: subj,
	}
}

// badReferenceDiagnostic returns a hcl.Diagnostic which indicates that a block
// reference is not expressed correctly.
func badReferenceDiagnostic(expr hcl.Expression) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity:   hcl.DiagError,
		Summary:    "Invalid block reference",
		Detail:     "A block reference is expressed as a block type separated from an identifier by a dot.",
		Subject:    expr.Range().Ptr(),
		Expression: expr,
	}
}
