package core

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/config"
	"bridgedl/config/lookup"
	"bridgedl/graph"
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

	evalCtx := evalContext(g)

	for _, v := range g.Vertices() {
		cmp, ok := v.(BridgeComponentVertex)
		if !ok {
			continue
		}

		cmpCat, cmpType, cmpID := cmp.Category(), cmp.Type(), cmp.Identifier()

		impl := t.Impls.ImplementationFor(cmpCat, cmpType)
		if impl == nil {
			// this would indicate a bug because initComponents
			// should validate that all components declared in the
			// bridge are at least partially implemented, so we
			// want such error to be loud
			panic(fmt.Errorf("no implementation for a %s of type %q. This should have been caught "+
				"earlier during validation.", cmpCat, cmpType))
		}

		transl, ok := impl.(translation.Translatable)
		if !ok {
			diags = diags.Append(noTranslatableDiagnostic(cmpCat, cmpType, cmp.SourceRange()))
			continue
		}

		var spec hcldec.Spec

		sp, ok := v.(AttachableSpecVertex)
		if ok {
			spec = sp.GetSpec()
		}

		dataSrc := lookup.NewDataSource(t.Bridge)

		config, cfgDiags := componentConfig(dataSrc, cmpCat, cmpID, spec, evalCtx)
		diags = diags.Extend(cfgDiags)

		eventDst, dstDiags := eventDestination(dataSrc, cmpCat, cmpID, evalCtx)
		diags = diags.Extend(dstDiags)

		manifests = append(manifests, transl.Manifests(cmpID, config, eventDst)...)
	}

	return manifests, diags
}

// evalContext returns an HCL evaluation context that allows evaluating
// traversal expressions within HCL blocks.
//
// Because traversal expressions always represent references to Addressables in
// the Bridge Description Language, all values nested of this EvalContext
// represent event addresses.
func evalContext(g *graph.DirectedGraph) *hcl.EvalContext {
	addrs := make(componentAddressesByCategory)

	// the keys of the graph's "up edges" are all vertices that are
	// referenced by at least one other vertex
	for v := range g.UpEdges() {
		addrCmp, ok := v.(AttachableAddressVertex)
		if !ok {
			// this would indicate a bug in the graph building
			// logic, so we want such error to be loud
			panic(errors.New("all referenceable vertices should also be addressable. This should have " +
				"been caught earlier during validation."))
		}

		cmpCat := addrCmp.Category()

		if _, ok := addrs[cmpCat]; !ok {
			addrs[cmpCat] = make(map[string]cty.Value)
		}

		addrs[cmpCat][addrCmp.Identifier()] = addrCmp.GetAddress()
	}

	evalCtx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value, len(addrs)),
	}

	for cmpCat, addrByCmpID := range addrs {
		evalCtx.Variables[cmpCat.String()] = cty.ObjectVal(addrByCmpID)
	}

	return evalCtx
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

// componentConfig returns the value decoded from the hcl.Config of the given
// component.
//
// It is acceptable for some component types to return a null cty.Value, if
// this type's configuration contains only optional attributes.
func componentConfig(ds *lookup.DataSource, cmpCat config.ComponentCategory, cmpID string,
	spec hcldec.Spec, evalCtx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {

	// some component types may not have any configuration body to decode
	// at all, in which case no hcldec.Spec will be passed as argument
	if spec == nil {
		return cty.NullVal(cty.DynamicPseudoType), nil
	}

	var bodyToDecode hcl.Body

	switch cmpCat {
	case config.CategoryChannels:
		bodyToDecode = ds.Channel(cmpID).Config
	case config.CategoryRouters:
		bodyToDecode = ds.Router(cmpID).Config
	case config.CategoryTransformers:
		bodyToDecode = ds.Transformer(cmpID).Config
	case config.CategorySources:
		bodyToDecode = ds.Source(cmpID).Config
	case config.CategoryTargets:
		bodyToDecode = ds.Target(cmpID).Config
	case config.CategoryFunctions:
		// "function" types are not supported yet, so there is
		// no HCL body to decode
		return cty.NullVal(cty.DynamicPseudoType), nil
	default:
		// should not happen, the list of categories is exhaustive
		panic(fmt.Errorf("encountered unknown component category %q", cmpCat))
	}

	return hcldec.Decode(bodyToDecode, spec, evalCtx)
}

// eventDestination returns the destination corresponding to the top-level "to"
// attribute of the given component, for component categories that include such
// attribute.
func eventDestination(ds *lookup.DataSource, cmpCat config.ComponentCategory, cmpID string,
	evalCtx *hcl.EvalContext) (cty.Value, hcl.Diagnostics) {

	var toExpr hcl.Traversal

	switch cmpCat {
	case config.CategoryChannels:
		toExpr = ds.Channel(cmpID).To

	case config.CategoryRouters:
		// "router" types only have nested block references
		return cty.NullVal(destinationCty), nil

	case config.CategoryTransformers:
		toExpr = ds.Transformer(cmpID).To

	case config.CategorySources:
		toExpr = ds.Source(cmpID).To

	case config.CategoryTargets:
		toExpr = ds.Target(cmpID).ReplyTo
		if toExpr == nil {
			return cty.NullVal(destinationCty), nil
		}

	case config.CategoryFunctions:
		toExpr = ds.Function(cmpID).ReplyTo
		if toExpr == nil {
			return cty.NullVal(destinationCty), nil
		}

	default:
		// should not happen, the list of categories is exhaustive
		panic(fmt.Errorf("encountered unknown component category %q", cmpCat))
	}

	return toExpr.TraverseAbs(evalCtx)
}
