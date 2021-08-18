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

package diagnostic_test

import (
	"strings"
	"testing"

	. "til/core/diagnostic"

	"github.com/hashicorp/hcl/v2"
)

func TestDedupDiagnostics(t *testing.T) {
	diag1 := newDiagnostic("Roof", "The roof is on fire", 10, 15)
	diag2 := newDiagnostic("House", "The house is on fire", 20, 25)
	diag3 := newDiagnostic("World", "The world is on fire", 30, 35)

	t.Run("Append non-equal diagnostics", func(t *testing.T) {
		d := NewDedupDiagnostics()

		d = d.Append(diag1)
		d = d.Append(diag2)
		d = d.Append(diag3)

		gotDiags := d.Diagnostics()

		const expectNumDiags = 3
		if gotNumDiags := len(gotDiags); gotNumDiags != expectNumDiags {
			t.Fatalf("Expected %d diagnostics, got %d:\n%s",
				expectNumDiags, gotNumDiags, errDiagsAsString(gotDiags))
		}
	})

	t.Run("Append equal diagnostics", func(t *testing.T) {
		d := NewDedupDiagnostics()

		d = d.Append(diag1)
		d = d.Append(diag1) // dupe
		d = d.Append(diag2)
		d = d.Append(diag1) // dupe

		gotDiags := d.Diagnostics()

		const expectNumDiags = 2
		if gotNumDiags := len(gotDiags); gotNumDiags != expectNumDiags {
			t.Fatalf("Expected %d diagnostics, got %d:\n%s",
				expectNumDiags, gotNumDiags, errDiagsAsString(gotDiags))
		}
	})

	t.Run("Extend with non-equal diagnostics", func(t *testing.T) {
		d := NewDedupDiagnostics()

		var diags hcl.Diagnostics
		diags = diags.Append(diag1)
		diags = diags.Append(diag2)
		diags = diags.Append(diag3)

		d = d.Extend(diags)

		gotDiags := d.Diagnostics()

		const expectNumDiags = 3
		if gotNumDiags := len(gotDiags); gotNumDiags != expectNumDiags {
			t.Fatalf("Expected %d diagnostics, got %d:\n%s",
				expectNumDiags, gotNumDiags, errDiagsAsString(gotDiags))
		}
	})

	t.Run("Extend with equal diagnostics", func(t *testing.T) {
		d := NewDedupDiagnostics()

		var diags hcl.Diagnostics
		diags = diags.Append(diag1)
		diags = diags.Append(diag1) // dupe
		diags = diags.Append(diag2)
		diags = diags.Append(diag1) // dupe

		d = d.Extend(diags)
		d = d.Extend(diags) // push the same list a second time

		gotDiags := d.Diagnostics()

		const expectNumDiags = 2
		if gotNumDiags := len(gotDiags); gotNumDiags != expectNumDiags {
			t.Fatalf("Expected %d diagnostics, got %d:\n%s",
				expectNumDiags, gotNumDiags, errDiagsAsString(gotDiags))
		}
	})
}

func newDiagnostic(summary, detail string, start, end int) *hcl.Diagnostic {
	rng := func(start, end int) hcl.Range {
		return hcl.Range{
			Filename: "mybridge.brg.hcl",
			Start:    hcl.Pos{Line: 0, Column: start, Byte: 0},
			End:      hcl.Pos{Line: 0, Column: end, Byte: 0},
		}
	}

	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  summary,
		Detail:   detail,
		Subject:  rng(start, end).Ptr(),
	}
}

// errDiagsAsString returns a string representation of all given error
// diagnostics as individual entries separated by a newline character.
func errDiagsAsString(diags hcl.Diagnostics) string {
	var dsb strings.Builder

	errDiags := diags.Errs()

	for i, d := range errDiags {
		dsb.WriteString(d.Error())
		if i < len(errDiags) {
			dsb.WriteRune('\n')
		}
	}

	return dsb.String()
}
