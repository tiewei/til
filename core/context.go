package core

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/graph"
)

// Context encapsulates everything that is required for performing operations
// on a Bridge.
type Context struct {
	Bridge *config.Bridge
	Specs  *Specs
}

func NewContext(brg *config.Bridge) (*Context, hcl.Diagnostics) {
	cmpImpls, diags := initComponents(brg)
	if diags.HasErrors() {
		return nil, diags
	}

	return &Context{
		Bridge: brg,
		Specs:  initSpecs(cmpImpls),
	}, nil
}

// Graph builds a directed graph which represents event flows between messaging
// components of a Bridge.
func (c *Context) Graph() (*graph.DirectedGraph, hcl.Diagnostics) {
	b := &GraphBuilder{
		Bridge: c.Bridge,
		Specs:  c.Specs,
	}

	return b.Build()
}
