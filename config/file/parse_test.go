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
	bridgeMissingAttrs = "missing_attrs.brg.hcl"
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
		if n := len(brg.Functions); n != 1 {
			t.Error("Expected 1 function, got", n)
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
			len(brg.Functions) != 0 ||
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
			len(brg.Functions) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all except sources to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with malformed block headers", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeBadBlkHdrs)

		errDiags := diags.Errs()

		const expectNumErrDiags = 4
		if len(errDiags) != expectNumErrDiags {
			t.Fatalf("Expected %d error diagnostics:\n%s", expectNumErrDiags, errDiagsAsString(diags))
		}

		if errDiags[0].(*hcl.Diagnostic).Summary != "Extraneous label for source" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Summary != "Missing identifier for source" {
			t.Fatal("Unexpected type of error diagnostic:", errDiags[1])
		}
		for i := 2; i < expectNumErrDiags; i++ {
			if errDiags[i].(*hcl.Diagnostic).Summary != "Invalid identifier" {
				t.Fatal("Unexpected type of error diagnostic:", errDiags[i])
			}
		}

		if errDiags[0].(*hcl.Diagnostic).Subject.Start.Line != 4 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[0])
		}
		if errDiags[1].(*hcl.Diagnostic).Subject.Start.Line != 13 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[1])
		}
		for i := 2; i < expectNumErrDiags; i++ {
			if errDiags[i].(*hcl.Diagnostic).Subject.Start.Line != 22 {
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
			len(brg.Functions) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all except sources to be empty. Parsed bridge:\n%+v", brg)
		}
	})

	t.Run("with missing attributes", func(t *testing.T) {
		brg, diags := p.LoadBridge(bridgeMissingAttrs)

		errDiags := diags.Errs()

		const expectNumErrDiags = 3
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
		if errDiags[2].(*hcl.Diagnostic).Subject.Start.Line != 20 {
			t.Fatal("Unexpected location of error diagnostic:", errDiags[2])
		}

		if n := len(brg.Channels); n != 1 {
			t.Error("Expected 1 channel, got", n)
		}
		if n := len(brg.Transformers); n != 1 {
			t.Error("Expected 1 transformer, got", n)
		}
		if n := len(brg.Sources); n != 1 {
			t.Error("Expected 1 source, got", n)
		}
		if len(brg.Routers) != 0 ||
			len(brg.Functions) != 0 ||
			len(brg.Targets) != 0 {

			t.Errorf("Expected all other components to be empty. Parsed bridge:\n%+v", brg)
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
