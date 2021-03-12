package core

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

	"bridgedl/config"
	"bridgedl/translation"

	"bridgedl/internal/components/channels"
	"bridgedl/internal/components/functions"
	"bridgedl/internal/components/routers"
	"bridgedl/internal/components/sources"
	"bridgedl/internal/components/targets"
	"bridgedl/internal/components/transformers"
)

// componentImpls encapsulates the implementation of all known
// component types for each supported component category (block type).
type componentImpls struct {
	channels     implForComponentType
	routers      implForComponentType
	transformers implForComponentType
	sources      implForComponentType
	targets      implForComponentType
	// "function" types are not supported yet
}

type implForComponentType map[string]interface{}

// ImplementationFor returns an implementation interface for the given
// component type, if it exists.
func (i *componentImpls) ImplementationFor(cmpCat config.ComponentCategory, cmpType string) interface{} {
	switch cmpCat {
	case config.CategoryChannels:
		return i.channels[cmpType]
	case config.CategoryRouters:
		return i.routers[cmpType]
	case config.CategoryTransformers:
		return i.transformers[cmpType]
	case config.CategorySources:
		return i.sources[cmpType]
	case config.CategoryTargets:
		return i.targets[cmpType]
	case config.CategoryFunctions:
		return (*functions.Function)(nil)
	default:
		// should not happen, the list of categories is exhaustive
		return nil
	}
}

// initComponents populates the component implementations associated with each
// component type present in the Bridge.
func initComponents(brg *config.Bridge) (*componentImpls, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	cmps := &componentImpls{
		channels:     make(implForComponentType),
		routers:      make(implForComponentType),
		transformers: make(implForComponentType),
		sources:      make(implForComponentType),
		targets:      make(implForComponentType),
	}

	for _, ch := range brg.Channels {
		if _, ok := cmps.channels[ch.Type]; ok {
			continue
		}

		impl, ok := channels.All[ch.Type]
		if !ok {
			diags = diags.Append(noComponentImplDiagnostic(config.CategoryChannels, ch.Type, ch.SourceRange))
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
			diags = diags.Append(noComponentImplDiagnostic(config.CategoryRouters, rtr.Type, rtr.SourceRange))
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
			diags = diags.Append(noComponentImplDiagnostic(config.CategoryTransformers, trsf.Type, trsf.SourceRange))
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
			diags = diags.Append(noComponentImplDiagnostic(config.CategorySources, src.Type, src.SourceRange))
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
			diags = diags.Append(noComponentImplDiagnostic(config.CategoryTargets, trg.Type, trg.SourceRange))
			continue
		}
		cmps.targets[trg.Type] = impl
	}

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
	// "function" types are not supported yet
}

type specForComponentType map[string]hcldec.Spec

// SpecFor returns a decode spec for the given component type, if it exists.
func (s *Specs) SpecFor(cmpCat config.ComponentCategory, cmpType string) hcldec.Spec {
	switch cmpCat {
	case config.CategoryChannels:
		return s.channels[cmpType]
	case config.CategoryRouters:
		return s.routers[cmpType]
	case config.CategoryTransformers:
		return s.transformers[cmpType]
	case config.CategorySources:
		return s.sources[cmpType]
	case config.CategoryTargets:
		return s.targets[cmpType]
	case config.CategoryFunctions:
		// "function" types are not supported yet, so we should never
		// have to decode a function config
		panic("not implemented")
	default:
		// should not happen, the list of categories is exhaustive
		return nil
	}
}

// initSpecs populates the Spec of each Decodable component type present in the Bridge.
func initSpecs(impls *componentImpls) *Specs {
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

	return specs
}

// Addressables encapsulates known Addressable component types for each
// supported component category (block type).
type Addressables struct {
	channels     addressableForComponentType
	routers      addressableForComponentType
	transformers addressableForComponentType
	sources      addressableForComponentType
	targets      addressableForComponentType
}

type addressableForComponentType map[string]translation.Addressable

// AddressableFor returns an Addressable interface for the given component
// type, if it exists.
func (a *Addressables) AddressableFor(cmpCat config.ComponentCategory, cmpType string) translation.Addressable {
	switch cmpCat {
	case config.CategoryChannels:
		return a.channels[cmpType]
	case config.CategoryRouters:
		return a.routers[cmpType]
	case config.CategoryTransformers:
		return a.transformers[cmpType]
	case config.CategorySources:
		return a.sources[cmpType]
	case config.CategoryTargets:
		return a.targets[cmpType]
	case config.CategoryFunctions:
		return (*functions.Function)(nil)
	default:
		// should not happen, the list of categories is exhaustive
		panic("unknown component category")
	}
}

// initAddressables populates an instance of Addressables with all known
// Addressable component types present in the Bridge.
func initAddressables(impls *componentImpls) *Addressables {
	specs := &Addressables{
		channels:     make(addressableForComponentType),
		routers:      make(addressableForComponentType),
		transformers: make(addressableForComponentType),
		sources:      make(addressableForComponentType),
		targets:      make(addressableForComponentType),
	}

	for cmpType, impl := range impls.channels {
		if addr, ok := impl.(translation.Addressable); ok {
			specs.channels[cmpType] = addr
		}
	}

	for cmpType, impl := range impls.routers {
		if addr, ok := impl.(translation.Addressable); ok {
			specs.routers[cmpType] = addr
		}
	}

	for cmpType, impl := range impls.transformers {
		if addr, ok := impl.(translation.Addressable); ok {
			specs.transformers[cmpType] = addr
		}
	}

	for cmpType, impl := range impls.sources {
		if addr, ok := impl.(translation.Addressable); ok {
			specs.sources[cmpType] = addr
		}
	}

	for cmpType, impl := range impls.targets {
		if addr, ok := impl.(translation.Addressable); ok {
			specs.targets[cmpType] = addr
		}
	}

	return specs
}
