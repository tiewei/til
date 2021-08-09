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

package addr

import (
	"github.com/hashicorp/hcl/v2"

	"til/config"
)

// Channel is the address of a "channel" block within a Bridge description.
type Channel struct {
	Identifier string
}

var _ Referenceable = (*Channel)(nil)

// Addr implements Referenceable.
func (ch Channel) Addr() string {
	return config.CategoryChannels.String() + "." + ch.Identifier
}

// Router is the address of a "router" block within a Bridge description.
type Router struct {
	Identifier string
}

var _ Referenceable = (*Router)(nil)

// Addr implements Referenceable.
func (rtr Router) Addr() string {
	return config.CategoryRouters.String() + "." + rtr.Identifier
}

// Transformer is the address of a "transformer" block within a Bridge description.
type Transformer struct {
	Identifier string
}

var _ Referenceable = (*Transformer)(nil)

// Addr implements Referenceable.
func (trsf Transformer) Addr() string {
	return config.CategoryTransformers.String() + "." + trsf.Identifier
}

// Source is the address of a "source" block within a Bridge description.
type Source struct {
	Identifier string
}

// Target is the address of a "target" block within a Bridge description.
type Target struct {
	Identifier string
}

var _ Referenceable = (*Target)(nil)

// Addr implements Referenceable.
func (trg Target) Addr() string {
	return config.CategoryTargets.String() + "." + trg.Identifier
}

// MessagingComponent is an address that can represent any messaging component
// within a Bridge description.
type MessagingComponent struct {
	Category    config.ComponentCategory
	Type        string
	Identifier  string
	SourceRange hcl.Range
}
