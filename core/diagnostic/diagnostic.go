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

package diagnostic

import "github.com/hashicorp/hcl/v2"

// DedupDiagnostics allows the accumulation of hcl.Diagnostic while avoiding
// duplicate entries. It implements a subset of the methods of hcl.Diagnostics.
type DedupDiagnostics struct {
	diags hcl.Diagnostics
	// An index of the diagnostic entries present in the diags list.
	index map[diagIndex]struct{}
}

// diagIndex is a reduced version of a hcl.Diagnostic which uniquely identifies
// a given instance of a particular diagnostic.
// It is meant to be used as the key to an index of diagnostics.
type diagIndex struct {
	summary string
	detail  string
	subject hcl.Range
}

// indexableDiag takes a hcl.Diagnostic and returns a corresponding diagIndex.
func indexableDiag(diag *hcl.Diagnostic) diagIndex {
	dIdx := diagIndex{
		summary: diag.Summary,
		detail:  diag.Detail,
	}

	if subj := diag.Subject; subj != nil {
		dIdx.subject = *subj
	}

	return dIdx
}

// NewDedupDiagnostics returns an initialized DedupDiagnostics.
func NewDedupDiagnostics() *DedupDiagnostics {
	return &DedupDiagnostics{
		index: make(map[diagIndex]struct{}),
	}
}

// Append appends a new hcl.Diagnostic to the receiver and returns the whole
// DedupDiagnostics.
func (d *DedupDiagnostics) Append(diag *hcl.Diagnostic) *DedupDiagnostics {
	diagIdx := indexableDiag(diag)

	if _, exists := d.index[diagIdx]; exists {
		return d
	}

	d.index[diagIdx] = struct{}{}
	d.diags = append(d.diags, diag)

	return d
}

// Extend concatenates the given hcl.Diagnostics with the receiver after
// filtering the diagnostics that already exist, and returns the whole
// DedupDiagnostics.
func (d *DedupDiagnostics) Extend(diags hcl.Diagnostics) *DedupDiagnostics {
	filteredDiags := diags[:0]

	for _, diag := range diags {
		diagIdx := indexableDiag(diag)

		if _, exists := d.index[diagIdx]; exists {
			continue
		}

		d.index[diagIdx] = struct{}{}
		filteredDiags = append(filteredDiags, diag)
	}

	d.diags = append(d.diags, filteredDiags...)

	return d
}

// Diagnostics returns the hcl.Diagnostics accumulated in the receiver.
func (d *DedupDiagnostics) Diagnostics() hcl.Diagnostics {
	return d.diags
}
