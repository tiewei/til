package core

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
)

// noComponentImplDiagnostic returns a hcl.Diagnostic which indicates that no
// implementation is available for a given component type.
func noComponentImplDiagnostic(cmpCat config.ComponentCategory, cmpType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not implemented",
		Detail:   fmt.Sprintf("No implementation found for a %s of type %q", cmpCat, cmpType),
		Subject:  subj.Ptr(),
	}
}

// noTranslatableDiagnostic returns a hcl.Diagnostic which indicates that a
// Translatable interface can not be acquired for a given component type.
func noTranslatableDiagnostic(cmpCat config.ComponentCategory, cmpType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not translatable",
		Detail:   fmt.Sprintf("Cannot find a translator for a %s of type %q", cmpCat, cmpType),
		Subject:  subj.Ptr(),
	}
}

// noDecodeSpecDiagnostic returns a hcl.Diagnostic which indicates that a spec
// for decoding a HCL body can not be acquired for a given component type.
func noDecodeSpecDiagnostic(cmpCat config.ComponentCategory, cmpType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "No decode spec",
		Detail:   fmt.Sprintf("Cannot find a decode spec for a %s of type %q", cmpCat, cmpType),
		Subject:  subj.Ptr(),
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
func noAddressableDiagnostic(cmpCat config.ComponentCategory, cmpType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Not addressable",
		Detail: fmt.Sprintf("Cannot determine the address for sending events to a %s of type %q",
			cmpCat, cmpType),
		Subject: subj.Ptr(),
	}
}

// wrongAddressTypeDiagnostic returns a hcl.Diagnostic which indicates that a
// component type implementation returned an unexpected type of event address.
func wrongAddressTypeDiagnostic(cmpCat config.ComponentCategory, cmpId string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Wrong address type",
		Detail: fmt.Sprintf("The event address computed for the %s %q is not a destination type",
			cmpCat, cmpId),
		Subject: subj.Ptr(),
	}
}
