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

package core

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"til/config/addr"
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

// wrongAddressTypeDiagnostic returns a hcl.Diagnostic which indicates that a
// component type implementation returned an unexpected type of event address.
func wrongAddressTypeDiagnostic(cmp addr.MessagingComponent) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Wrong address type",
		Detail: fmt.Sprintf("The event address computed for the %s %q is not a destination type",
			cmp.Category, cmp.Identifier),
		Subject: cmp.SourceRange.Ptr(),
	}
}

// undecodableDiagnostic returns a hcl.Diagnostic which indicates that the
// topology of the Bridge resulted in a component that could not be decoded.
func undecodableDiagnostic(cmp addr.MessagingComponent) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Undecodable",
		Detail: fmt.Sprintf("Attempts to decode the %s %q resulted in errors due to expressions that "+
			"could not be resolved at runtime. Please ensure the Bridge topology is valid.",
			cmp.Category, cmp.Type),
		Subject: cmp.SourceRange.Ptr(),
	}
}
