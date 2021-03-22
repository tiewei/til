package lang

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config"
	"bridgedl/k8s"
)

// DecodeSafe attempts to decode the given hcl.Body with the information
// contained in the provided hcl.EvalContext.
// Unlike hcldec.Decode, missing variables are tolerated and, in case of error,
// the decoding is retried with a placeholder injected in place of each missing
// variable which represents a block reference. When this occurs, the returned
// boolean value is false to indicate that the body was decoded with
// placeholders, and is therefore incomplete/inaccurate.
//
// This function exists to compensate for topology cycles, which can be
// legitimate in a Bridge description but require visiting a certain component
// multiple times: once to determine its event address for populating the
// evaluation context that will be consumed by predecessors in the cycle
// (without translation), once at a later point when the event addresses of all
// successors in the cycle have been determined (with translation).
//
// Passing an inaccurate "duck" destination value to an Addressable
// implementations is assumed to be acceptable and should have no influence on
// the returned event addresses, since only the block configuration should
// matter, not the actual events destinations. On the opposite, passing an
// incomplete configuration to a Translatable implementation would yield
// invalid results.
func DecodeSafe(b hcl.Body, s hcldec.Spec, ctx *hcl.EvalContext) (cty.Value, bool /*complete*/, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	config, decDiags := hcldec.Decode(b, s, ctx)
	if !decDiags.HasErrors() {
		return config, true, diags.Extend(decDiags)
	}

	config, decDiags = decodeIgnoreUnknownRefs(b, s, ctx)
	diags = diags.Extend(decDiags)

	return config, false, diags
}

// TraverseAbsSafe is similar to DecodeSafe but evaluates a single absolute
// hcl.Traversal.
func TraverseAbsSafe(t hcl.Traversal, ctx *hcl.EvalContext) (cty.Value, bool /*complete*/, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	val, evalDiags := t.TraverseAbs(ctx)
	if !evalDiags.HasErrors() {
		return val, true, diags.Extend(evalDiags)
	}

	defaultVal := cty.NullVal(k8s.DestinationCty)
	ctx = evalContextEnsureVars(ctx, defaultVal, filterBlockRefs(t)...)

	val, evalDiags = t.TraverseAbs(ctx)
	diags = diags.Extend(evalDiags)

	return val, false, diags
}

// decodeIgnoreUnknownRefs decodes a hcl.Body after replacing each unknown
// block reference in the given hcl.EvalContext with a null cty value.
func decodeIgnoreUnknownRefs(b hcl.Body, s hcldec.Spec, ctx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {
	blockRefs := filterBlockRefs(hcldec.Variables(b, s)...)

	defaultVal := cty.NullVal(k8s.DestinationCty)
	ctx = evalContextEnsureVars(ctx, defaultVal, blockRefs...)

	return hcldec.Decode(b, s, ctx)
}

// filterBlockRefs filters a list of hcl.Traversal, keeping only the elements
// which represent a block reference.
func filterBlockRefs(vars ...hcl.Traversal) []hcl.Traversal {
	if len(vars) == 0 {
		return nil
	}

	blockRefs := vars[:0]

	for _, t := range vars {
		cmpCat := config.AsComponentCategory(t.RootName())
		if !referenceableTypes().Has(cmpCat) {
			continue
		}

		blockRefs = append(blockRefs, t)
	}

	return blockRefs
}
