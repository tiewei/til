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

package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/fs"
	"bridgedl/graph"
	"bridgedl/lang/k8s"
	"bridgedl/translation"
)

// BridgeTranslator translates a Bridge into a collection of Kubernetes API
// objects for deploying that Bridge.
type BridgeTranslator struct {
	Impls *componentImpls

	// used by functions that access the file system
	BaseDir string
	FS      fs.FS
}

// Translate performs the translation.
func (t *BridgeTranslator) Translate(g *graph.DirectedGraph) ([]interface{}, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	var bridgeManifests []interface{}

	eval := NewEvaluator(t.BaseDir, t.FS)

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
	// missing the event address of its successor, and this gap is
	// compensated for by injecting a temporary placeholder value into the
	// evaluation context. The "incomplete" component is re-evaluated once
	// all addresses in the cycle have been determined.
	sccs := g.StronglyConnectedComponents()

	for _, scc := range sccs {
		manifests, translDiags := translateComponents(eval, scc)
		diags = diags.Extend(translDiags)

		bridgeManifests = append(bridgeManifests, manifests...)
	}

	return bridgeManifests, diags
}

// translateComponents translates all components from a list of graph vertices.
func translateComponents(e *Evaluator, vs []graph.Vertex) ([]interface{}, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	var manifests []interface{}
	var incompleteDecodeQueue []MessagingComponentVertex

	// first pass: push addresses of visited components to the Evaluator,
	// translate immediately if the addresses of all successors allowed the
	// config body to be fully decoded
	for _, v := range vs {
		cmp, ok := v.(MessagingComponentVertex)
		if !ok {
			continue
		}

		ref, ok := v.(ReferenceableVertex)
		if ok {
			evalDiags := appendToEvaluator(e, ref)
			diags = diags.Extend(evalDiags)
		}

		cfg, evDst, complete, cfgDiags := configAndDestination(e, cmp)
		diags = diags.Extend(cfgDiags)
		if cfgDiags.HasErrors() {
			continue
		}

		// We are in a cycle and not all successors have been evaluated
		// yet. Push the component to a queue and delay its second
		// evaluation to a point where the Evaluator contains all
		// necessary addresses.
		if !complete {
			incompleteDecodeQueue = append(incompleteDecodeQueue, cmp)
			continue
		}

		res, translDiags := translate(cmp, cfg, evDst)
		diags = diags.Extend(translDiags)

		manifests = append(manifests, res...)
	}

	// second pass: remaining components which evaluation was delayed due
	// to cycles (incomplete evaluation context)
	for _, cmp := range incompleteDecodeQueue {
		cfg, evDst, complete, cfgDiags := configAndDestination(e, cmp)
		diags = diags.Extend(cfgDiags)
		if cfgDiags.HasErrors() {
			continue
		}

		// It is not acceptable to end up with an incomplete decoding
		// on a second visit. After all vertices in the strongly
		// connected componenent have been evaluated, the Evaluator
		// should contain all addresses.
		if !complete {
			diags = diags.Append(undecodableDiagnostic(cmp.ComponentAddr()))
			continue
		}

		res, translDiags := translate(cmp, cfg, evDst)
		diags = diags.Extend(translDiags)

		manifests = append(manifests, res...)
	}

	return manifests, diags
}

// translate invokes the translator of the given component.
func translate(cmp MessagingComponentVertex, cfg, evDst cty.Value) (
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

// appendToEvaluator appends the event address of the given referenceable
// Bridge component to an Evaluator.
func appendToEvaluator(e *Evaluator, cmp ReferenceableVertex) hcl.Diagnostics {
	var diags hcl.Diagnostics

	cmpAddr := cmp.ComponentAddr()
	if e.HasVariable(cmpAddr.Category.String(), cmpAddr.Identifier) {
		return diags
	}

	evAddr, _, evAddrDiags := cmp.EventAddress(e)
	diags = diags.Extend(evAddrDiags)
	if evAddrDiags.HasErrors() {
		return diags
	}

	e.InsertVariable(cmpAddr.Category.String(), cmpAddr.Identifier, evAddr)

	return diags
}

// configAndDestination returns the decoded configuration of the given
// component as well as its event destination, if applicable.
func configAndDestination(e *Evaluator, v MessagingComponentVertex) (
	cfg, evDst cty.Value, complete bool, diags hcl.Diagnostics) {

	cfgComplete := true
	dstComplete := true

	cfg = cty.NullVal(cty.DynamicPseudoType)
	if dec, ok := v.(DecodableConfigVertex); ok {
		var cfgDiags hcl.Diagnostics
		cfg, cfgComplete, cfgDiags = dec.DecodedConfig(e)
		diags = diags.Extend(cfgDiags)
	}

	evDst = cty.NullVal(k8s.DestinationCty)
	if sdr, ok := v.(EventSenderVertex); ok {
		var dstDiags hcl.Diagnostics
		evDst, dstComplete, dstDiags = sdr.EventDestination(e)
		diags = diags.Extend(dstDiags)
	}

	complete = cfgComplete && dstComplete

	return cfg, evDst, complete, diags
}
