package core

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"bridgedl/config/addr"
)

// noComponentImplDiagnostic returns a hcl.Diagnostic which indicates that no
// implementation is available for a given component type.
func noComponentImplDiagnostic(cmp addr.MessagingComponent) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not implemented",
		Detail:   fmt.Sprintf("No implementation found for a %s of type %q", cmp.Category, cmp.Type),
		Subject:  cmp.SourceRange.Ptr(),
	}
}

// noTranslatableDiagnostic returns a hcl.Diagnostic which indicates that a
// Translatable interface can not be acquired for a given component type.
func noTranslatableDiagnostic(cmp addr.MessagingComponent) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not translatable",
		Detail:   fmt.Sprintf("Cannot find a translator for a %s of type %q", cmp.Category, cmp.Type),
		Subject:  cmp.SourceRange.Ptr(),
	}
}

// unknownReferenceDiagnostic returns a hcl.Diagnostic which indicates that a block
// reference refers to an unknown block.
func unknownReferenceDiagnostic(refAddr addr.Referenceable, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Reference to unknown block",
		Detail:   fmt.Sprintf("The expression %q doesn't match any known configuration block", refAddr.Addr()),
		Subject:  subj.Ptr(),
	}
}

// noAddressableDiagnostic returns a hcl.Diagnostic which indicates that an
// Addressable interface can not be acquired for a given component type.
func noAddressableDiagnostic(cmp addr.MessagingComponent) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not addressable",
		Detail: fmt.Sprintf("Cannot determine the address for sending events to a %s of type %q",
			cmp.Category, cmp.Type),
		Subject: cmp.SourceRange.Ptr(),
	}
}
