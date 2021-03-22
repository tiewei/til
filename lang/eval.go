package lang

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// evalContextEnsureVars returns a hcl.EvalContext where each given variable is
// replaced with the provided default value in case this variable is missing
// from the original hcl.EvalContext.
func evalContextEnsureVars(ctx *hcl.EvalContext, defaultVal cty.Value, vars ...hcl.Traversal) *hcl.EvalContext {
	missing := missingVars(ctx, vars...)
	if len(missing) == 0 {
		return ctx
	}

	ctxCpy := &hcl.EvalContext{}
	*ctxCpy = *ctx
	ctxCpy.Variables = make(map[string]cty.Value, len(ctx.Variables))

	for root, val := range ctx.Variables {
		missingAttrs, hasMissingAttrs := missing[root]
		delete(missing, root)

		// no variable missing for the given root, just reuse the
		// cty.Value from the existing hcl.EvalContext
		if !hasMissingAttrs {
			ctxCpy.Variables[root] = val
			continue
		}

		// expand the existing cty.Value into a map (assuming that
		// value is a cty.Object type), populate it with placeholders,
		// and recreate a cty.Object from those merged attributes
		numAttrs := val.LengthInt() + len(missingAttrs)
		attrs := make(map[string]cty.Value, numAttrs)

		valIter := val.ElementIterator()
		for valIter.Next() {
			attr, v := valIter.Element()
			attrs[attr.AsString()] = v
		}

		for _, attr := range missingAttrs {
			// if the missing attribute is already in the list, it
			// means we have a serious bug in the missingVars logic
			// and therefore make the failure very loud
			if _, exists := attrs[attr]; exists {
				panic("conflict: variable flagged as missing but present in the evaluation context")
			}
			attrs[attr] = defaultVal
		}

		ctxCpy.Variables[root] = cty.ObjectVal(attrs)
	}

	// if missing still contains elements, it means some roots don't exist
	// at all in the original hcl.EvalContext, so we add them without merge
	for root, missingAttrs := range missing {
		attrs := make(map[string]cty.Value, len(missingAttrs))
		for _, attr := range missingAttrs {
			attrs[attr] = defaultVal
		}
		ctxCpy.Variables[root] = cty.ObjectVal(attrs)
	}

	return ctxCpy
}

// missingVars finds missing variables in the given hcl.EvalContext.
func missingVars(ctx *hcl.EvalContext, vars ...hcl.Traversal) map[string][]string {
	if len(vars) == 0 {
		return nil
	}

	missing := make(map[string][]string)

	for _, v := range vars {
		vs := v.SimpleSplit()
		root := vs.RootName()

		if len(vs.Rel) != 1 {
			continue
		}
		attr := vs.Rel[0].(hcl.TraverseAttr).Name

		val, ok := ctx.Variables[root]
		if !ok || !val.Type().HasAttribute(attr) {
			missing[root] = append(missing[root], attr)
		}
	}

	return missing
}
