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

package core_test

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	. "bridgedl/core"
	"bridgedl/fs"
)

func TestEvaluator(t *testing.T) {
	const baseDir = "/fake/bridge/dir"

	t.Run("insert variables", func(t *testing.T) {
		const root1, root2 = "foo", "bar"
		const var1, var2 = "var1", "var2"

		// nil FS is safe b/c test suite does not invoke functions
		e := NewEvaluator(baseDir, nil)

		if e.HasVariable(root1, var1) {
			t.Fatal("Expected variable to not exist in the empty Evaluator")
		}
		if nVars := len(e.EvalContext().Variables); nVars != 0 {
			t.Fatal("Expected new Evaluator to have an empty list of variables, got", nVars)
		}

		testVal := testVariableValue()

		// initial insertion

		e.InsertVariable(root1, var1, testVal)

		if !e.HasVariable(root1, var1) {
			t.Fatal("Expected variable to exist in the Evaluator")
		}

		ctxVars := e.EvalContext().Variables
		if nVars := len(ctxVars); nVars != 1 {
			t.Fatal("Expected single root in Evaluator, got", nVars)
		}
		varVal := ctxVars[root1].GetAttr(var1)
		if varVal.Equals(testVal).False() {
			t.Fatal("Retrieved unexpected variable from Evaluator:", varVal.GoString())
		}

		// re-insertion with different value

		e.InsertVariable(root1, var1, cty.StringVal("test"))

		ctxVars = e.EvalContext().Variables
		if nVars := len(ctxVars); nVars != 1 {
			t.Fatal("Expected re-insertion to not create a new variable in Evaluator, got", nVars)
		}
		varVal = ctxVars[root1].GetAttr(var1)
		if varVal.Equals(testVal).False() {
			t.Fatal("Expected re-insertion to not update existing variable in Evaluator, got", varVal.GoString())
		}

		// inserton of a new variable at the same root

		e.InsertVariable(root1, var2, testVal)

		if !e.HasVariable(root1, var2) {
			t.Fatal("Expected variable to exist in the Evaluator")
		}

		ctxVars = e.EvalContext().Variables
		if nVars := len(ctxVars); nVars != 1 {
			t.Fatal("Expected single root in Evaluator, got", nVars)
		}
		varVal = ctxVars[root1].GetAttr(var2)
		if varVal.Equals(testVal).False() {
			t.Fatal("Retrieved unexpected variable from Evaluator:", varVal.GoString())
		}

		// insertion of new variable at a different root

		e.InsertVariable(root2, var1, testVal)

		if !e.HasVariable(root2, var1) {
			t.Fatal("Expected variable to exist in the Evaluator")
		}

		ctxVars = e.EvalContext().Variables
		if nVars := len(ctxVars); nVars != 2 {
			t.Fatal("Expected 2 roots in Evaluator, got", nVars)
		}
		varVal = ctxVars[root2].GetAttr(var1)
		if varVal.Equals(testVal).False() {
			t.Fatal("Retrieved unexpected variable from Evaluator:", varVal.GoString())
		}
	})

	t.Run("evaluate configurations", func(t *testing.T) {
		const fakeFileRelPath = "functions/myfunction.js"
		const fakeFileContents = "fake file contents"

		testFS := fs.NewMemFS()
		_ = testFS.CreateFile(filepath.Join(baseDir, fakeFileRelPath),
			[]byte(fakeFileContents),
		)

		// Those tests target basic use cases, such as the capacity of the
		// Evaluator to properly propagate variables and functions to an
		// EvalContext.
		// The "lang" package is better suited for testing the decoding
		// of more complex HCL bodies.
		testBody := hclBody(t, ``+
			`attr_var = my.var`+"\n"+
			`attr_func = file("`+fakeFileRelPath+`")`,
		)
		testSpec := hcldec.ObjectSpec(map[string]hcldec.Spec{
			"attr_var":  &hcldec.AttrSpec{Name: "attr_var", Type: cty.Bool},
			"attr_func": &hcldec.AttrSpec{Name: "attr_func", Type: cty.String},
		})

		e := NewEvaluator(baseDir, testFS)
		e.InsertVariable("my", "var", cty.True)

		// decode block

		v, _, diags := e.DecodeBlock(testBody, testSpec)
		if diags.HasErrors() {
			t.Fatal("Failed to decode HCL block:", diags)
		}

		if v := v.GetAttr("attr_var"); !v.True() {
			t.Error(`Unexpected value for "attr_var":`, v.GoString())
		}
		if v := v.GetAttr("attr_func"); v.AsString() != fakeFileContents {
			t.Errorf(`Unexpected value for "attr_func": %q`, v.GoString())
		}

		// decode expr

		attrs, diags := testBody.JustAttributes()
		if diags.HasErrors() {
			t.Fatal("Failed to interpret HCL attributes:", diags)
		}

		tr, _ := hcl.AbsTraversalForExpr(attrs["attr_var"].Expr) // we know it's a valid Traversal

		v, _, diags = e.DecodeTraversal(tr)
		if diags.HasErrors() {
			t.Fatal("Failed to decode HCL traversal:", diags)
		}

		if !v.True() {
			t.Error(`Unexpected value for "attr_var":`, v.GoString())
		}
	})
}

// testVariableValue returns a structured cty value that can be used as
// variable value in tests.
func testVariableValue() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"some_attr":  cty.True,
		"other_attr": cty.Zero,
	})
}

// hclBody parses the given HCL code and returns it as a hcl.Body.
func hclBody(t *testing.T, code string) hcl.Body {
	t.Helper()

	p := hclparse.NewParser()

	f, diags := p.ParseHCL([]byte(code), "irrelevant_filename.hcl")
	if diags.HasErrors() {
		t.Fatal("Failed to parse HCL:", diags)
	}

	return f.Body
}
