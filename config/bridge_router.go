package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// routerBlockSchema is the shallow structure of a "router" block.
// Used for validation during decoding.
var routerBlockSchema = &hcl.BodySchema{}

// Router represents a generic message router.
type Router struct {
	// Indicates which type of router is contained within the block.
	Type string
	// An identifier that is unique among all Routers within a Bridge.
	Identifier string

	// Configuration of the router.
	Config hcl.Body
}

// decodeRouterBlock performs a partial decoding of the Body of a "router"
// block into a Router struct.
func decodeRouterBlock(blk *hcl.Block) (*Router, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	_, remain, contentDiags := blk.Body.PartialContent(routerBlockSchema)
	diags = diags.Extend(contentDiags)

	rtr := &Router{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		Config:     remain,
	}

	return rtr, diags
}
