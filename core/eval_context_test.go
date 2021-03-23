package core_test

import (
	"testing"

	"github.com/zclconf/go-cty/cty"

	. "bridgedl/core"
)

func TestEvaluator(t *testing.T) {
	t.Run("insert variables", func(t *testing.T) {
		const root1, root2 = "foo", "bar"
		const var1, var2 = "var1", "var2"

		e := NewEvaluator()

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
}

// testVariableValue returns a structured cty value that can be used as
// variable value in tests.
func testVariableValue() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"some_attr":  cty.True,
		"other_attr": cty.Zero,
	})
}
