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

package lang_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config/file"
	"bridgedl/fs"
	. "bridgedl/lang"
	"bridgedl/lang/k8s"
)

// List of available fixture files.
// Should be kept in sync with the content of the "fixtures/" directory.
const (
	bodyMultiBlkRefsOneType   = "multiple_blk_refs_one_type.hcl"
	bodyMultiBlkRefsMultiType = "multiple_blk_refs_multi_type.hcl"
	bodyOneBlkRefOneNonblkRef = "one_blk_ref_one_nonblk_ref.hcl"
)

func TestDecodeSafe(t *testing.T) {
	p := &file.Parser{
		Parser: hclparse.NewParser(),
		FS:     populatedFixtureFS(t),
	}

	s := fixtureHCLSpec()

	k8sDst := k8sDestination()

	t.Run("fully populated context with single root and valid values", func(t *testing.T) {
		b := mustParseHCLFile(t, p, bodyMultiBlkRefsOneType)

		ctx := newEvalContext()
		ctx.Variables["channel"] = cty.ObjectVal(map[string]cty.Value{
			"some_channel":  k8sDst,
			"other_channel": k8sDst,
		})

		_, complete, diags := DecodeSafe(b, s, ctx)
		if diags.HasErrors() {
			t.Fatal("Failed to decode test HCL:", diags)
		}

		if !complete {
			t.Error("Expected decoding to be complete, but the EvalContext was missing values")
		}
	})

	t.Run("fully populated context with multiple roots and valid values", func(t *testing.T) {
		b := mustParseHCLFile(t, p, bodyMultiBlkRefsMultiType)

		ctx := newEvalContext()
		ctx.Variables["channel"] = cty.ObjectVal(map[string]cty.Value{
			"some_channel": k8sDst,
		})
		ctx.Variables["router"] = cty.ObjectVal(map[string]cty.Value{
			"some_router": k8sDst,
		})

		_, complete, diags := DecodeSafe(b, s, ctx)
		if diags.HasErrors() {
			t.Fatal("Failed to decode test HCL:", diags)
		}

		if !complete {
			t.Error("Expected decoding to be complete, but the EvalContext was missing values")
		}
	})

	t.Run("fully populated context with non-block references", func(t *testing.T) {
		b := mustParseHCLFile(t, p, bodyOneBlkRefOneNonblkRef)

		ctx := newEvalContext()
		ctx.Variables["channel"] = cty.ObjectVal(map[string]cty.Value{
			"some_channel": k8sDst,
		})
		ctx.Variables["foo"] = cty.ObjectVal(map[string]cty.Value{
			"bar": cty.StringVal("baz"),
		})

		_, complete, diags := DecodeSafe(b, s, ctx)
		if diags.HasErrors() {
			t.Fatal("Failed to decode test HCL:", diags)
		}

		if !complete {
			t.Error("Expected decoding to be complete, but the EvalContext was missing values")
		}
	})

	t.Run("fully populated context with one invalid value", func(t *testing.T) {
		b := mustParseHCLFile(t, p, bodyMultiBlkRefsMultiType)

		ctx := newEvalContext()
		ctx.Variables["channel"] = cty.ObjectVal(map[string]cty.Value{
			"some_channel": cty.NumberIntVal(1234), // (!) number unexpected here
		})
		ctx.Variables["router"] = cty.ObjectVal(map[string]cty.Value{
			"some_router": k8sDst,
		})

		_, _, diags := DecodeSafe(b, s, ctx)
		if !diags.HasErrors() {
			t.Fatal("Expected failure decoding HCL config due to incorrect value type")
		}

		errDiags := diags.Errs()
		if nErrs := len(errDiags); nErrs != 1 {
			t.Error("Expected a single error diagnostic, got", nErrs)
		}
		if errDiags[0].(*hcl.Diagnostic).Summary != "Incorrect attribute value type" {
			t.Error("Unexpected type of error diagnostic:", errDiags[0])
		}
	})

	t.Run("partially populated context with only block references", func(t *testing.T) {
		b := mustParseHCLFile(t, p, bodyMultiBlkRefsMultiType)

		ctx := newEvalContext()
		ctx.Variables["channel"] = cty.ObjectVal(map[string]cty.Value{
			"some_channel": k8sDst,
		}) // (!) value of router reference is missing

		_, complete, diags := DecodeSafe(b, s, ctx)
		if diags.HasErrors() {
			t.Fatal("Failed to decode test HCL:", diags)
		}

		if complete {
			t.Error("Expected decoding to be incomplete due to missing values in the EvalContext")
		}
	})

	t.Run("partially populated context with non-block references", func(t *testing.T) {
		b := mustParseHCLFile(t, p, bodyOneBlkRefOneNonblkRef)

		ctx := newEvalContext()
		ctx.Variables["channel"] = cty.ObjectVal(map[string]cty.Value{
			"some_channel": k8sDst,
		}) // (!) value of foo var is missing

		_, _, diags := DecodeSafe(b, s, ctx)
		if !diags.HasErrors() {
			t.Fatal("Expected failure decoding HCL config due to missing non-block variable in the EvalContext")
		}

		errDiags := diags.Errs()
		if nErrs := len(errDiags); nErrs != 1 {
			t.Error("Expected a single error diagnostic, got", nErrs)
		}
		if errDiags[0].(*hcl.Diagnostic).Summary != "Unknown variable" {
			t.Error("Unexpected type of error diagnostic:", errDiags[0])
		}
	})
}

// fixtureHCLSpec returns a HCL decode spec that is used to decode
// configuration bodies from this package's fixture files.
func fixtureHCLSpec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"some_attr": &hcldec.AttrSpec{
			Name:     "some_attr",
			Type:     cty.String,
			Required: true,
		},
		"optional_attr": &hcldec.AttrSpec{
			Name:     "optional_attr",
			Type:     cty.Number,
			Required: false,
		},
		"some_ref": &hcldec.AttrSpec{
			Name:     "some_ref",
			Type:     k8s.DestinationCty,
			Required: true,
		},
		"errors": &hcldec.BlockSpec{
			TypeName: "errors",
			Nested: &hcldec.ObjectSpec{
				"nested_attr": &hcldec.AttrSpec{
					Name:     "nested_attr",
					Type:     cty.Bool,
					Required: false,
				},
				"nested_ref": &hcldec.AttrSpec{
					Name:     "nested_ref",
					Type:     k8s.DestinationCty,
					Required: true,
				},
			},
			Required: false,
		},
	}
}

func newEvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
	}
}

// k8sDestination returns a cty value which is a valid "duck" Destination object.
func k8sDestination() cty.Value {
	return k8s.NewDestination("v0", "Foo", "foo")
}

func mustParseHCLFile(t *testing.T, p *file.Parser, f string) hcl.Body {
	t.Helper()

	hclFile, diags := p.ParseHCLFile(f)
	if diags.HasErrors() {
		t.Fatalf("Failed to parse HCL file %q: %v", f, diags)
	}

	return hclFile.Body
}

// populatedFixtureFS returns a fs.FS populated with all *.hcl files from the
// "fixtures/" directory.
func populatedFixtureFS(t *testing.T) fs.FS {
	t.Helper()

	const hclExt = ".hcl"
	const fixturesDir = "fixtures/"

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatal("Error reading fixtures directory:", err)
	}

	mfs := fs.NewMemFS()

	for _, e := range entries {
		if e.IsDir() || e.Name()[len(e.Name())-len(hclExt):] != hclExt {
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
