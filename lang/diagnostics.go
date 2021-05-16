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

package lang

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
func badRefTypeDiagnostic(blockType string, subj hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid block reference",
		Detail: fmt.Sprintf("The expression refers to a block type %q that is unknown or doesn't support "+
			"references. Valid values are %v", blockType, referenceableTypes()),
		Subject: subj.Ptr(),
	}
}
