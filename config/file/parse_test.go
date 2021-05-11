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

package file_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	. "bridgedl/config/file"
	"bridgedl/fs"
)

// List of available fixture files.
// Should be kept in sync with the content of the "fixtures/" directory.
const (
	bridgeValid        = "valid.brg.hcl"
	bridgeUnknBlkType  = "unkn_blk_type.brg.hcl"
	bridgeBadBlkRefs   = "bad_blk_refs.brg.hcl"
	bridgeBadBlkHdrs   = "bad_blk_hdrs.brg.hcl"
	bridgeBadGlobals   = "bad_globals.brg.hcl"
	bridgeMissingAttrs = "missing_attrs.brg.hcl"
	bridgeDuplIDs      = "dupl_ids.brg.hcl"
	bridgeDuplGlobals  = "dupl_globals.brg.hcl"
)

func TestLoadBridge(t *testing.T) {
	p := &Parser{
		Parser: hclparse.NewParser(),
		FS:     populatedFixtureFS(t),
	}

	t.Run("valid description", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeValid)
		if diags.HasErrors() {
			t.Fatalf("Returned error diagnostics:\n%s", errDiagsAsString(diags))
		}

		if n := len(brg.Channels); n != 1 {
			t.Error("Expected 1 channel, got", n)
		}
		if n := len(brg.Routers); n != 1 {
			t.Error("Expected 1 router, got", n)
		}
		if n := len(brg.Transformers); n != 1 {
			t.Error("Expected 1 transformer, got", n)
		}
		if n := len(brg.Sources); n != 1 {
			t.Error("Expected 1 source, got", n)
		}
		if n := len(brg.Targets); n != 1 {
			t.Error("Expected 1 target, got", n)
		}
	})

	t.Run("with unknown block", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeUnknBlkType)

		errDiags := diags.Errs()

		const expectNumErrDiags = 1
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostic:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}
		if errDiags[0].(*hcl.Diagnostic).Summary != "Unsupported block type" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[0])
		}
		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 12 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}

		const expectNumSources = 1
		if n := len(brg.Sources); n != expectNumSources {
			t.Errorf("Expected %d source, got %d", expectNumSources, n)
		}
		if len(brg.Channels) != 0 ||
			len(brg.Routers) != 0 ||
			len(brg.Transformers) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all except sources to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with malformed block references", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeBadBlkRefs)

		errDiags := diags.Errs()

		const expectNumErrDiags = 3
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostics:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		for i := 0; i < expectNumErrDiags; i++ {
			if errDiags[i].(*hcl.Diagnostic).Summary != "Invalid expression" {
				t.Fatal("Unexpected type of error diagnostic:", errDiags[i])
			}
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 9 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Subject.Start.Line != 18 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[1])
		}
		if errDiags[2].(*hcl.Diagnostic).Subject.Start.Line != 27 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[2])
		}

		const expectNumSources = 3
		if n := len(brg.Sources); n != expectNumSources {
			t.Errorf("Expected %d sources, got %d", expectNumSources, n)
		}
		if len(brg.Channels) != 0 ||
			len(brg.Routers) != 0 ||
			len(brg.Transformers) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all except sources to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with malformed block headers", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeBadBlkHdrs)

		errDiags := diags.Errs()

		const expectNumErrDiags = 5
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostics:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		if errDiags[0].(*hcl.Diagnostic).Summary != "Missing identifier for bridge" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Summary != "Extraneous label for source" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[1])
		}
		if errDiags[2].(*hcl.Diagnostic).Summary != "Missing identifier for source" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[2])
		}
		for i := 3; i < expectNumErrDiags; i++ {
			if errDiags[i].(*hcl.Diagnostic).Summary != "Invalid identifier" {
				t.Fatal("Unexpected type of error diagnostic:", errDiags[i])
			}
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 4 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Subject.Start.Line != 7 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[1])
		}
		if errDiags[2].(*hcl.Diagnostic).Subject.Start.Line != 16 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[2])
		}
		for i := 3; i < expectNumErrDiags; i++ {
			if errDiags[i].(*hcl.Diagnostic).Subject.Start.Line != 25 {
				t.Fatal("Unexpected location of error diagnostic:", errDiags[i])
			}
		}

		const expectNumSources = 1
		if n := len(brg.Sources); n != expectNumSources {
			t.Errorf("Expected %d source, got %d", expectNumSources, n)
		}
		if len(brg.Channels) != 0 ||
			len(brg.Routers) != 0 ||
			len(brg.Transformers) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all except sources to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with malformed global settings", func(t *testing.T) {
		_, diags := p.LoadBridge(bridgeBadGlobals)

		errDiags := diags.Errs()

		const expectNumErrDiags = 2
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostics:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		if errDiags[0].(*hcl.Diagnostic).Summary != "Wrong attribute type" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Summary != "Invalid expression" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[1])
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 6 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Subject.Start.Line != 8 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[1])
		}
	})

	t.Run("with missing attributes", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeMissingAttrs)

		errDiags := diags.Errs()

		const expectNumErrDiags = 2
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostics:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		for i := 0; i < expectNumErrDiags; i++ {
			if errDiags[i].(*hcl.Diagnostic).Summary != "Missing required argument" {
				t.Fatal("Unexpected type of error diagnostic:", errDiags[i])
			}
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 4 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Subject.Start.Line != 12 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[1])
		}

		if n := len(brg.Transformers); n != 1 {
			t.Error("Expected 1 transformer, got", n)
		}
		if n := len(brg.Sources); n != 1 {
			t.Error("Expected 1 source, got", n)
		}
		if len(brg.Channels) != 0 ||
			len(brg.Routers) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all other components to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with duplicated identifiers", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeDuplIDs)

		errDiags := diags.Errs()

		const expectNumErrDiags = 1
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostic:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		if errDiags[0].(*hcl.Diagnostic).Summary != "Duplicate block" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[0])
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 12 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}

		if n := len(brg.Sources); n != 1 {
			t.Error("Expected 1 source, got", n)
		}
		if len(brg.Routers) != 0 ||
			len(brg.Channels) != 0 ||
			len(brg.Transformers) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all other components to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with duplicated global settings", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeDuplGlobals)

		errDiags := diags.Errs()

		const expectNumErrDiags = 1
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostic:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		if errDiags[0].(*hcl.Diagnostic).Summary != "Redefined global config" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[0])
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 6 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}

		if n := len(brg.Sources); n != 1 {
			t.Error("Expected 1 source, got", n)
		}
	})
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

// populatedFixtureFS returns a fs.FS populated with all *.brg.hcl files from
// the "fixtures/" directory.
func populatedFixtureFS(t *testing.T) fs.FS {
	t.Helper()

	const brgExt = ".brg.hcl"
	const fixturesDir = "fixtures/"

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatal("Error reading fixtures directory:", err)
	}

	mfs := fs.NewMemFS()

	for _, e := range entries {
		if e.IsDir() || e.Name()[len(e.Name())-len(brgExt):] != brgExt {
			continue
		}

		data, err := os.ReadFile(filepath.Join(fixturesDir, e.Name()))
		if err != nil {
			t.Fatal("Error reading fixture file:", err)
		}

		if err := mfs.CreateFile(e.Name(), data); err != nil {
			t.Fatal("Error populating FS:", err)
		}
	}

	return mfs
}
