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

	var manifests []interface{}

	evalCtx, evalDiags := evalContext(g)
	diags = diags.Extend(evalDiags)

	for _, v := range g.Vertices() {
		cmp, ok := v.(MessagingComponentVertex)
		if !ok {
			continue
		}

		var loopDiags hcl.Diagnostics

		config := cty.NullVal(cty.DynamicPseudoType)
		if dec, ok := v.(DecodableConfigVertex); ok {
			var cfgDiags hcl.Diagnostics
			config, cfgDiags = dec.DecodedConfig(evalCtx)
			loopDiags = loopDiags.Extend(cfgDiags)
		}

		eventDst := cty.NullVal(k8s.DestinationCty)
		if sdr, ok := v.(EventSenderVertex); ok {
			var dstDiags hcl.Diagnostics
			eventDst, dstDiags = sdr.EventDestination(evalCtx)
			loopDiags = loopDiags.Extend(dstDiags)
		}

		diags = diags.Extend(loopDiags)
		if loopDiags.HasErrors() {
			continue
		}

		cmpAddr := cmp.ComponentAddr()

		transl, ok := cmp.Implementation().(translation.Translatable)
		if !ok {
			diags = diags.Append(noTranslatableDiagnostic(cmpAddr))
			continue
		}

		manifests = append(manifests, transl.Manifests(cmpAddr.Identifier, config, eventDst)...)
	}

	return manifests, diags
}

// evalContext returns an HCL evaluation context that allows evaluating
// traversal expressions within HCL blocks.
//
// Because traversal expressions always represent references to Addressables in
// the Bridge Description Language, all values nested of this EvalContext
// represent event addresses.
func evalContext(g *graph.DirectedGraph) (*hcl.EvalContext, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	addrs := make(componentAddressesByCategory)

	// the keys of the graph's "up edges" are all vertices that are
	// referenced by at least one other vertex
	for v := range g.UpEdges() {
		addr, ok := v.(AddressableVertex)
		if !ok {
			continue
		}

		cmpAddr := addr.ComponentAddr()

		evAddr, evAddrDiags := addr.EventAddress()
		diags = diags.Extend(evAddrDiags)

		if _, ok := addrs[cmpAddr.Category]; !ok {
			addrs[cmpAddr.Category] = make(map[string]cty.Value)
		}

		addrs[cmpAddr.Category][cmpAddr.Identifier] = evAddr
	}

	evalCtx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value, len(addrs)),
	}

	for cmpCat, addrByCmpID := range addrs {
		evalCtx.Variables[cmpCat.String()] = cty.ObjectVal(addrByCmpID)
	}

	return evalCtx, diags
}

// componentAddressesByCategory is a collection of maps of component ID to
// event address, indexed by component category.
//
// It is used as a temporary data store for assembling an hcl.EventContext.
//
// Example:
//   "router": {
//     "my_router": <address value>
//   }
//   "channel": {
//     "my_channel": <address value>
//   }
type componentAddressesByCategory map[config.ComponentCategory]map[string]cty.Value
