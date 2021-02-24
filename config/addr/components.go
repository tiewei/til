package addr

import "bridgedl/config"

// Channel is the address of a "channel" block within a Bridge description.
type Channel struct {
	Identifier string
}

var _ Referenceable = (*Channel)(nil)

// Addr implements Referenceable.
func (ch Channel) Addr() string {
	return config.BlkChannel + "." + ch.Identifier
}

// Router is the address of a "router" block within a Bridge description.
type Router struct {
	Identifier string
}

var _ Referenceable = (*Router)(nil)

// Addr implements Referenceable.
func (rtr Router) Addr() string {
	return config.BlkRouter + "." + rtr.Identifier
}

// Transformer is the address of a "transformer" block within a Bridge description.
type Transformer struct {
	Identifier string
}

var _ Referenceable = (*Transformer)(nil)

// Addr implements Referenceable.
func (trsf Transformer) Addr() string {
	return config.BlkTransf + "." + trsf.Identifier
}

// Source is the address of a "source" block within a Bridge description.
type Source struct {
	Identifier string
}

var _ Referenceable = (*Source)(nil)

// Addr implements Referenceable.
//
// Like other Referenceable types, Source implements the interface for the
// purpose of being identifiable as a configuration block. However, it is
// invalid to reference a Source as the destination of events within a Bridge.
func (src Source) Addr() string {
	return config.BlkSource + "." + src.Identifier
}

// Target is the address of a "target" block within a Bridge description.
type Target struct {
	Identifier string
}

var _ Referenceable = (*Target)(nil)

// Addr implements Referenceable.
func (trg Target) Addr() string {
	return config.BlkTarget + "." + trg.Identifier
}

// Function is the address of a "function" block within a Bridge description.
type Function struct {
	Identifier string
}

var _ Referenceable = (*Function)(nil)

// Addr implements Referenceable.
func (ch Function) Addr() string {
	return config.BlkChannel + "." + ch.Identifier
}
