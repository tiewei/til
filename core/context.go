package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/graph"
	"bridgedl/translate"
)

// Context encapsulates everything that is required for performing operations
// on a Bridge.
type Context struct {
	Bridge      *config.Bridge
	Translators *translate.TranslatorProviders
}

// Graph builds a directed graph which represents event flows between messaging
// components of a Bridge.
func (c *Context) Graph() (*graph.DirectedGraph, hcl.Diagnostics) {
	b := &GraphBuilder{
		Bridge:      c.Bridge,
		Translators: c.Translators,
	}

	return b.Build()
}
