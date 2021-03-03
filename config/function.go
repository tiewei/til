package config

import "github.com/hashicorp/hcl/v2"

// FunctionBlockSchema is the shallow structure of a "function" block.
// Used for validation during decoding.
var FunctionBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrReplyTo,
		Required: false,
	}},
}

// Function represents a function.
type Function struct {
	// An identifier that is unique among all Functions within a Bridge.
	Identifier string

	// Destination of event responses.
	ReplyTo hcl.Traversal

	// Configuration of the function.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
