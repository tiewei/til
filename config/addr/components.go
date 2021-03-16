package addr

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
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

// Function is the address of a "function" block within a Bridge description.
type Function struct {
	Identifier string
}

var _ Referenceable = (*Function)(nil)

// Addr implements Referenceable.
func (ch Function) Addr() string {
	return config.CategoryFunctions.String() + "." + ch.Identifier
}

// MessagingComponent is an address that can represent any messaging component
// within a Bridge description.
type MessagingComponent struct {
	Category    config.ComponentCategory
	Type        string
	Identifier  string
	SourceRange hcl.Range
}
