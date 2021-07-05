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

	"til/config"
	"til/config/addr"

	"til/internal/components/channels"
	"til/internal/components/routers"
	"til/internal/components/sources"
	"til/internal/components/targets"
	"til/internal/components/transformers"
)

// componentImpls encapsulates the implementation of all known
// component types for each supported component category (block type).
type componentImpls struct {
	channels     implForComponentType
	routers      implForComponentType
	transformers implForComponentType
	sources      implForComponentType
	targets      implForComponentType
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
			diags = diags.Append(noComponentImplDiagnostic(addr.MessagingComponent{
				Category:    config.CategoryChannels,
				Type:        ch.Type,
				SourceRange: ch.SourceRange,
			}))
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
			diags = diags.Append(noComponentImplDiagnostic(addr.MessagingComponent{
				Category:    config.CategoryRouters,
				Type:        rtr.Type,
				SourceRange: rtr.SourceRange,
			}))
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
			diags = diags.Append(noComponentImplDiagnostic(addr.MessagingComponent{
				Category:    config.CategoryTransformers,
				Type:        trsf.Type,
				SourceRange: trsf.SourceRange,
			}))
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
			diags = diags.Append(noComponentImplDiagnostic(addr.MessagingComponent{
				Category:    config.CategorySources,
				Type:        src.Type,
				SourceRange: src.SourceRange,
			}))
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
			diags = diags.Append(noComponentImplDiagnostic(addr.MessagingComponent{
				Category:    config.CategoryTargets,
				Type:        trg.Type,
				SourceRange: trg.SourceRange,
			}))
			continue
		}
		cmps.targets[trg.Type] = impl
	}

	return cmps, diags
}
