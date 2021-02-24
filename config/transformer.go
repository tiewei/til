package config

import "github.com/hashicorp/hcl/v2"

// TransformerBlockSchema is the shallow structure of a "transformer" block.
// Used for validation during decoding.
var TransformerBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrTo,
		Required: true,
	}},
}

// Transformer represents a generic message transformer.
type Transformer struct {
	// Indicates which type of transformer is contained within the block.
	Type string
	// An identifier that is unique among all Transformers within a Bridge.
	Identifier string

	// Destination of events.
	To hcl.Traversal

	// Configuration of the transformer.
	Config hcl.Body
}
