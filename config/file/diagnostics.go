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

package file

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"til/config"
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

// tooManyGlobalBlocksDiagnostic returns a hcl.Diagnostic which indicates that
// more than one block of the given type was defined in the global configuration.
func tooManyGlobalBlocksDiagnostic(blkType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Redefined global config",
		Detail:   fmt.Sprintf("A Bridge description must contain at most one %q block.", blkType),
		Subject:  subj.Ptr(),
	}
}

// wrongTypeDiagnostic returns a validation diagnostic which indicates that the
// given attribute value doesn't have the expected type.
func wrongTypeDiagnostic(v cty.Value, expectType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Wrong attribute type",
		Detail:   "The attribute value is not a " + expectType + ". Type: " + v.Type().FriendlyNameForConstraint(),
		Subject:  subj.Ptr(),
	}
}
