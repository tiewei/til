package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"bridgedl/fs"
	"bridgedl/lang"
)

// Evaluator can evaluate graph vertices by providing access to variables and
// functions that are required for decoding HCL configurations.
//
// Traversal expressions always represent references to Addressable blocks in
// the current version of the Bridge Description Language. Therefore, all
// variables values stored in this Evaluator represent event addresses. This
// may change in the future.
type Evaluator struct {
	variables variablesIndexedByRoot
	functions map[string]function.Function

	// a hcl.EvalContext matching the current state of the Evaluator can be
	// cached to avoid re-generating it if no new variable was inserted
	// since the last retrieval
	cachedEvalCtx *hcl.EvalContext
}

// NewEvaluator returns an initialized Evaluator.
func NewEvaluator(baseDir string, fs fs.FS) *Evaluator {
	return &Evaluator{
		variables: make(variablesIndexedByRoot),
		functions: lang.Functions(baseDir, fs),
	}
}

// InsertVariable inserts a variable in the current Evaluator.
func (e *Evaluator) InsertVariable(root, varname string, val cty.Value) {
	if e.HasVariable(root, varname) {
		return
	}

	// invalidate cached hcl.EvalContext to force its re-creation upon next
	// retrieval
	e.cachedEvalCtx = nil

	if _, exists := e.variables[root]; !exists {
		e.variables[root] = make(map[string]cty.Value, 1)
	}
	e.variables[root][varname] = val
}

// HasVariable returns whether the current Evaluator contains a value for the
// given variable.
func (e *Evaluator) HasVariable(root, varname string) bool {
	vars, hasRoot := e.variables[root]
	if !hasRoot {
		return false
	}

	_, hasVar := vars[varname]
	return hasVar
}

// DecodeBlock evaluates the value of a configuration block.
//
// The returned boolean value indicates whether all expressions from the
// configuration body could be decoded without injecting placeholders into the
// evaluation context.
func (e *Evaluator) DecodeBlock(b hcl.Body, s hcldec.Spec) (cty.Value, bool, hcl.Diagnostics) {
	return lang.DecodeSafe(b, s, e.EvalContext())
}

// DecodeTraversal evaluates the value of a single traversal.
//
// The returned boolean value indicates whether the traversal expression could
// be decoded without injecting a placeholder into the evaluation context.
func (e *Evaluator) DecodeTraversal(t hcl.Traversal) (cty.Value, bool, hcl.Diagnostics) {
	return lang.TraverseAbsSafe(t, e.EvalContext())
}

// EvalContext returns an hcl.EvalContext which contains all variables and
// functions from the current Evaluator, in a suitable format.
func (e *Evaluator) EvalContext() *hcl.EvalContext {
	if e.cachedEvalCtx != nil {
		return e.cachedEvalCtx
	}

	evalCtx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value, len(e.variables)),
		Functions: e.functions,
	}

	for root, vars := range e.variables {
		evalCtx.Variables[root] = cty.ObjectVal(vars)
	}

	return evalCtx
}

// variablesIndexedByRoot is a collection of maps of variables names to values
// indexed by traversal root. It is intended to be used as a temporary data
// store for assembling an hcl.EvalContext.
//
// Example:
//   "router": {
//     "my_router": <address value>
//   }
//   "channel": {
//     "my_channel": <address value>
//   }
type variablesIndexedByRoot map[string]map[string]cty.Value
