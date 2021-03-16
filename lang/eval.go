package lang

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/k8s"
)

// DecodeIgnoreVars decodes a hcl.Body but evaluates each encoutered variable
// as a null cty value.
//
// This is used in situations where a consumer needs access to a decoded
// configuration body while the values of variables contained in this body are
// still unknown. For instance, in calls to (translation.Addressable).Address.
//
// NOTE(antoineco): we could avoid this hack with a different evaluation model.
// By applying a topological sort to the graph representation of the Bridge
// instead of visiting its vertices in a random order, it is possible to ensure
// that the value of a given variable is known prior to translating the
// vertices that depend on that variable.
func DecodeIgnoreVars(b hcl.Body, s hcldec.Spec) (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if b == nil || s == nil {
		return cty.NullVal(cty.DynamicPseudoType), diags
	}

	evalCtx := evalContextNullVars(b, s)
	config, decDiags := hcldec.Decode(b, s, evalCtx)
	diags = diags.Extend(decDiags)

	return config, diags
}

// evalContextNullVars generates an hcl.EvalContext for the given hcl.Body
// where every variable has a null cty value.
func evalContextNullVars(b hcl.Body, s hcldec.Spec) *hcl.EvalContext {
	traversals := hcldec.Variables(b, s)

	vars := make(map[string]map[string]cty.Value)

	for _, t := range traversals {
		ts := t.SimpleSplit()

		root := ts.RootName()

		if vars[root] == nil {
			vars[root] = make(map[string]cty.Value)
		}
		vars[root][ts.Rel[0].(hcl.TraverseAttr).Name] = cty.NullVal(k8s.DestinationCty)
	}

	evalCtx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
	}

	for root, nulls := range vars {
		evalCtx.Variables[root] = cty.ObjectVal(nulls)
	}

	return evalCtx
}
