package config

import "github.com/hashicorp/hcl/v2"

// TargetBlockSchema is the shallow structure of a "target" block.
// Used for validation during decoding.
var TargetBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrReplyTo,
		Required: false,
	}},
}

// Target represents a generic event target.
type Target struct {
	// Indicates which type of target is contained within the block.
	Type string
	// An identifier that is unique among all Targets within a Bridge.
	Identifier string

	// Destination of event responses.
	ReplyTo hcl.Traversal

	// Configuration of the target.
	Config hcl.Body
}
