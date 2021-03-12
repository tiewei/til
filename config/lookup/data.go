package lookup

import (
	"bridgedl/config"
	"bridgedl/config/addr"
)

// DataSource exposes lookup functions that provide access to the data
// contained in a Bridge configuration.
//
// It uses lookup keys which match the values returned by the `*Key` functions
// in this package.
type DataSource struct {
	config *config.Bridge
}

// NewDataSource returns a DataSource for the given Bridge.
func NewDataSource(brg *config.Bridge) *DataSource {
	return &DataSource{
		config: brg,
	}
}

// Channel returns the Channel with the given identifier if it exists.
func (s *DataSource) Channel(identifier string) *config.Channel {
	if len(s.config.Channels) == 0 {
		return nil
	}
	return s.config.Channels[addr.Channel{Identifier: identifier}]
}

// Router returns the Router with the given identifier if it exists.
func (s *DataSource) Router(identifier string) *config.Router {
	if len(s.config.Routers) == 0 {
		return nil
	}
	return s.config.Routers[addr.Router{Identifier: identifier}]
}

// Transformer returns the Transformer with the given identifier if it exists.
func (s *DataSource) Transformer(identifier string) *config.Transformer {
	if len(s.config.Transformers) == 0 {
		return nil
	}
	return s.config.Transformers[addr.Transformer{Identifier: identifier}]
}

// Source returns the Source with the given identifier if it exists.
func (s *DataSource) Source(identifier string) *config.Source {
	if len(s.config.Sources) == 0 {
		return nil
	}
	return s.config.Sources[addr.Source{Identifier: identifier}]
}

// Target returns the Target with the given identifier if it exists.
func (s *DataSource) Target(identifier string) *config.Target {
	if len(s.config.Targets) == 0 {
		return nil
	}
	return s.config.Targets[addr.Target{Identifier: identifier}]
}

// Function returns the Function with the given identifier if it exists.
func (s *DataSource) Function(identifier string) *config.Function {
	if len(s.config.Functions) == 0 {
		return nil
	}
	return s.config.Functions[addr.Function{Identifier: identifier}]
}
