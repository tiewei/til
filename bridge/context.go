package bridge

import (
	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/graph"
	"bridgedl/graph/build"
)

// Context encapsulates everything that is required for performing operations
// on a Bridge.
type Context struct {
	Bridge *config.Bridge
}

// Graph builds a directed graph which represents event flows between messaging
// components of a Bridge.
func (c *Context) Graph() (*graph.DirectedGraph, hcl.Diagnostics) {
	b := &build.GraphBuilder{
		Bridge: c.Bridge,
	}

	return b.Build()
}
