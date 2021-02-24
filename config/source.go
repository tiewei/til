package config

import "github.com/hashicorp/hcl/v2"

// SourceBlockSchema is the shallow structure of a "source" block.
// Used for validation during decoding.
var SourceBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrTo,
		Required: true,
	}},
}

// Source represents a generic event source.
type Source struct {
	// Indicates which type of source is contained within the block.
	Type string
	// An identifier that is unique among all Sources within a Bridge.
	Identifier string

	// Destination of events.
	To hcl.Traversal

	// Configuration of the source.
	Config hcl.Body
}
