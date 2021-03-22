package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config"
	"bridgedl/graph"
	"bridgedl/k8s"
	"bridgedl/translation"
)

// BridgeTranslator translates a Bridge into a collection of Kubernetes API
// objects for deploying that Bridge.
type BridgeTranslator struct {
	Impls  *componentImpls
	Bridge *config.Bridge
}

// Translate performs the translation.
func (t *BridgeTranslator) Translate(g *graph.DirectedGraph) ([]interface{}, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	var bridgeManifests []interface{}

	evalCtx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
	}

	// The returned SCCs are topologically sorted as a byproduct of the
	// cycle detection algorithm.
	//
	// If the Bridge is acyclic, there is no strongly connected component
	// with a size greater than 1. The topological order guarantees that
	// each evaluated component has access to all the variables (event
	// addresses) it needs, because those were evaluated ahead of time.
	//
	// If the Bridge contains cycles, it exists strongly connected
	// components with a size greater than 1 within the graph. In this
	// case, the first component evaluated in each cycle will always be
	// missing the event address of its successor, and this gap will be
	// compensated for in the evaluation context with a placeholder value.
	sccs := g.StronglyConnectedComponents()

	for _, scc := range sccs {
		var translDiags hcl.Diagnostics
		var manifests []interface{}

		manifests, evalCtx, translDiags = translateComponents(evalCtx, scc)
		diags = diags.Extend(translDiags)

		bridgeManifests = append(bridgeManifests, manifests...)
	}

	return bridgeManifests, diags
}

// translateComponents translates all components from a list of graph vertices.
func translateComponents(evalCtx *hcl.EvalContext, vs []graph.Vertex) ([]interface{}, *hcl.EvalContext, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	var manifests []interface{}
	var incompleteDecodeQueue []MessagingComponentVertex

	// first pass: push addresses of visited components to the
	// evaluation context, translate immediately if the address of
	// all successors allowed the config body to be fully decoded
	for _, v := range vs {
		cmp, ok := v.(MessagingComponentVertex)
		if !ok {
			continue
		}

		ref, ok := v.(ReferenceableVertex)
		if ok {
			var evalCtxDiags hcl.Diagnostics
			evalCtx, evalCtxDiags = appendToEvalContext(evalCtx, ref)
			diags = diags.Extend(evalCtxDiags)
		}

		cfg, evDst, complete, cfgDiags := configAndDestination(evalCtx, cmp)
		diags = diags.Extend(cfgDiags)
		if cfgDiags.HasErrors() {
			continue
		}

		// We are in a cycle and not all successors have been
		// evaluated yet. Push the component to a queue and
		// delay its second evaluation to a point where the
		// evaluation context contains all necessary addresses.
		if !complete {
			incompleteDecodeQueue = append(incompleteDecodeQueue, cmp)
			continue
		}

		res, translDiags := translate(evalCtx, cmp, cfg, evDst)
		diags = diags.Extend(translDiags)

		manifests = append(manifests, res...)
	}

	// second pass: remaining components which evaluation was
	// delayed due to cycles (incomplete evaluation context)
	for _, cmp := range incompleteDecodeQueue {
		cfg, evDst, complete, cfgDiags := configAndDestination(evalCtx, cmp)
		diags = diags.Extend(cfgDiags)
		if cfgDiags.HasErrors() {
			continue
		}

		// It is not acceptable to end up with an incomplete
		// decoding on a second visit. After all vertices in
		// the strongly connected componenent have been walked,
		// the evaluation context should contain all addresses.
		if !complete {
			diags = diags.Append(undecodableDiagnostic(cmp.ComponentAddr()))
			continue
		}

		res, translDiags := translate(evalCtx, cmp, cfg, evDst)
		diags = diags.Extend(translDiags)

		manifests = append(manifests, res...)
	}

	return manifests, evalCtx, diags
}

// translate invokes the translator of the given component.
func translate(ctx *hcl.EvalContext, cmp MessagingComponentVertex, cfg, evDst cty.Value) (
	[]interface{}, hcl.Diagnostics) {

	var diags hcl.Diagnostics

	cmpAddr := cmp.ComponentAddr()

	transl, ok := cmp.Implementation().(translation.Translatable)
	if !ok {
		diags = diags.Append(noTranslatableDiagnostic(cmpAddr))
		return nil, diags
	}

	return transl.Manifests(cmpAddr.Identifier, cfg, evDst), diags
}

// appendToEvalContext appends the event address of the given referenceable
// Bridge component to an existing hcl.EvalContext. The hcl.EvalContext is
// modified in place.
//
// Because traversal expressions always represent references to Addressables in
// the Bridge Description Language, all values nested of this EvalContext
// represent event addresses.
func appendToEvalContext(ctx *hcl.EvalContext, cmp ReferenceableVertex) (*hcl.EvalContext, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	cmpAddr := cmp.ComponentAddr()

	vals, exists := ctx.Variables[cmpAddr.Category.String()]
	if exists && vals.Type().HasAttribute(cmpAddr.Identifier) {
		// evaluation context already contains a value for this component
		return ctx, diags
	}

	evAddr, _, evAddrDiags := cmp.EventAddress(ctx)
	diags = diags.Extend(evAddrDiags)
	if evAddrDiags.HasErrors() {
		return ctx, diags
	}

	var attrs map[string]cty.Value

	if exists {
		attrs = make(map[string]cty.Value, vals.LengthInt()+1)

		valsIter := vals.ElementIterator()
		for valsIter.Next() {
			attr, val := valsIter.Element()
			attrs[attr.AsString()] = val
		}
	} else {
		attrs = make(map[string]cty.Value, 1)
	}

	attrs[cmpAddr.Identifier] = evAddr

	ctx.Variables[cmpAddr.Category.String()] = cty.ObjectVal(attrs)

	return ctx, diags
}

// configAndDestination returns the decoded configuration of the given
// component as well as its event destination, if applicable.
func configAndDestination(ctx *hcl.EvalContext, v MessagingComponentVertex) (
	cfg, evDst cty.Value, complete bool, diags hcl.Diagnostics) {

	cfgComplete := true
	dstComplete := true

	cfg = cty.NullVal(cty.DynamicPseudoType)
	if dec, ok := v.(DecodableConfigVertex); ok {
		var cfgDiags hcl.Diagnostics
		cfg, cfgComplete, cfgDiags = dec.DecodedConfig(ctx)
		diags = diags.Extend(cfgDiags)
	}

	evDst = cty.NullVal(k8s.DestinationCty)
	if sdr, ok := v.(EventSenderVertex); ok {
		var dstDiags hcl.Diagnostics
		evDst, dstComplete, dstDiags = sdr.EventDestination(ctx)
		diags = diags.Extend(dstDiags)
	}

	complete = cfgComplete && dstComplete

	return cfg, evDst, complete, diags
}
