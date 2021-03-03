package config

import "github.com/hashicorp/hcl/v2"

// RouterBlockSchema is the shallow structure of a "router" block.
// Used for validation during decoding.
var RouterBlockSchema = &hcl.BodySchema{}

// Router represents a generic message router.
type Router struct {
	// Indicates which type of router is contained within the block.
	Type string
	// An identifier that is unique among all Routers within a Bridge.
	Identifier string

	// Configuration of the router.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
