package core

import (
	"path/filepath"

	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/fs"
	"bridgedl/graph"
)

// Context encapsulates everything that is required for performing operations
// on a Bridge.
type Context struct {
	Bridge *config.Bridge
	Impls  *componentImpls

	// interface used by functions that access the file system
	FS fs.FS
}

func NewContext(brg *config.Bridge) (*Context, hcl.Diagnostics) {
	cmpImpls, diags := initComponents(brg)
	if diags.HasErrors() {
		return nil, diags
	}

	return &Context{
		Bridge: brg,
		Impls:  cmpImpls,

		FS: (*fs.OSFS)(nil),
	}, nil
}

// Graph builds a directed graph which represents event flows between messaging
// components of a Bridge.
func (c *Context) Graph() (*graph.DirectedGraph, hcl.Diagnostics) {
	b := &GraphBuilder{
		Bridge: c.Bridge,
		Impls:  c.Impls,
	}

	return b.Build()
}

// Generate generates the deployment manifests for a Bridge.
func (c *Context) Generate() ([]interface{}, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	g, graphDiags := c.Graph()
	diags = diags.Extend(graphDiags)
	if diags.HasErrors() {
		return nil, diags
	}

	t := &BridgeTranslator{
		Impls: c.Impls,

		BaseDir: filepath.Dir(c.Bridge.Path),
		FS:      c.FS,
	}

	return t.Translate(g)
}
