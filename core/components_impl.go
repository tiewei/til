package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

	"bridgedl/config"
	"bridgedl/translation"

	"bridgedl/internal/components/channels"
	"bridgedl/internal/components/routers"
	"bridgedl/internal/components/sources"
	"bridgedl/internal/components/targets"
	"bridgedl/internal/components/transformers"
)

// componentImplementations encapsulates the implementation of all known
// component types for each supported component category (block type).
type componentImplementations struct {
	channels     implForComponentType
	routers      implForComponentType
	transformers implForComponentType
	sources      implForComponentType
	targets      implForComponentType
	functions    implForComponentType
}

type implForComponentType map[string]interface{}

// initComponents populates the component implementations associated with each
// component type present in the Bridge.
func initComponents(brg *config.Bridge) (*componentImplementations, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	cmps := &componentImplementations{
		channels:     make(implForComponentType),
		routers:      make(implForComponentType),
		transformers: make(implForComponentType),
		sources:      make(implForComponentType),
		targets:      make(implForComponentType),
		functions:    make(implForComponentType),
	}

	for _, ch := range brg.Channels {
		if _, ok := cmps.channels[ch.Type]; ok {
			continue
		}

		impl, ok := channels.All[ch.Type]
		if !ok {
			diags = diags.Append(noComponentImplDiagnostic(config.BlkChannel, ch.Type, ch.SourceRange))
			continue
		}
		cmps.channels[ch.Type] = impl
	}

	for _, rtr := range brg.Routers {
		if _, ok := cmps.routers[rtr.Type]; ok {
			continue
		}

		impl, ok := routers.All[rtr.Type]
		if !ok {
			diags = diags.Append(noComponentImplDiagnostic(config.BlkRouter, rtr.Type, rtr.SourceRange))
			continue
		}
		cmps.routers[rtr.Type] = impl
	}

	for _, trsf := range brg.Transformers {
		if _, ok := cmps.transformers[trsf.Type]; ok {
			continue
		}

		impl, ok := transformers.All[trsf.Type]
		if !ok {
			diags = diags.Append(noComponentImplDiagnostic(config.BlkTransf, trsf.Type, trsf.SourceRange))
			continue
		}
		cmps.transformers[trsf.Type] = impl
	}

	for _, src := range brg.Sources {
		if _, ok := cmps.sources[src.Type]; ok {
			continue
		}

		impl, ok := sources.All[src.Type]
		if !ok {
			diags = diags.Append(noComponentImplDiagnostic(config.BlkSource, src.Type, src.SourceRange))
			continue
		}
		cmps.sources[src.Type] = impl
	}

	for _, trg := range brg.Targets {
		if _, ok := cmps.targets[trg.Type]; ok {
			continue
		}

		impl, ok := targets.All[trg.Type]
		if !ok {
			diags = diags.Append(noComponentImplDiagnostic(config.BlkTarget, trg.Type, trg.SourceRange))
			continue
		}
		cmps.targets[trg.Type] = impl
	}

	// "function" types are not supported yet

	return cmps, diags
}

// Specs encapsulates the hcldec.Spec of known component types for each
// supported component category (block type).
type Specs struct {
	channels     specForComponentType
	routers      specForComponentType
	transformers specForComponentType
	sources      specForComponentType
	targets      specForComponentType
	function     specForComponentType
}

type specForComponentType map[string]hcldec.Spec

// SpecFor returns a decode spec for the given component type, if it exists.
func (s *Specs) SpecFor(cmpCat componentCategory, cmpType string) hcldec.Spec {
	switch cmpCat {
	case componentChannels:
		return s.channels[cmpType]
	case categoryRouters:
		return s.routers[cmpType]
	case categoryTransformers:
		return s.transformers[cmpType]
	case categorySources:
		return s.sources[cmpType]
	case categoryTargets:
		return s.targets[cmpType]
	case categoryFunctions:
		// "function" types are not supported yet
		return nil
	default:
		// should not happen, the list of categories is exhaustive
		return nil
	}
}

// initSpecs populates the Spec of each Decodable component type present in the Bridge.
func initSpecs(impls *componentImplementations) *Specs {
	specs := &Specs{
		channels:     make(specForComponentType),
		routers:      make(specForComponentType),
		transformers: make(specForComponentType),
		sources:      make(specForComponentType),
		targets:      make(specForComponentType),
	}

	for cmpType, impl := range impls.channels {
		if dec, ok := impl.(translation.Decodable); ok {
			specs.channels[cmpType] = dec.Spec()
		}
	}

	for cmpType, impl := range impls.routers {
		if dec, ok := impl.(translation.Decodable); ok {
			specs.routers[cmpType] = dec.Spec()
		}
	}

	for cmpType, impl := range impls.transformers {
		if dec, ok := impl.(translation.Decodable); ok {
			specs.transformers[cmpType] = dec.Spec()
		}
	}

	for cmpType, impl := range impls.sources {
		if dec, ok := impl.(translation.Decodable); ok {
			specs.sources[cmpType] = dec.Spec()
		}
	}

	for cmpType, impl := range impls.targets {
		if dec, ok := impl.(translation.Decodable); ok {
			specs.targets[cmpType] = dec.Spec()
		}
	}

	// "function" types are not supported yet

	return specs
}

type componentCategory uint8

const (
	componentChannels componentCategory = iota
	categoryRouters
	categoryTransformers
	categorySources
	categoryTargets
	categoryFunctions
)
