package config

import "github.com/hashicorp/hcl/v2"

// ChannelBlockSchema is the shallow structure of a "channel" block.
// Used for validation during decoding.
var ChannelBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrTo,
		Required: true,
	}},
}

// Channel represents a generic messaging channel.
type Channel struct {
	// Indicates which type of channel is contained within the block.
	Type string
	// An identifier that is unique among all Channels within a Bridge.
	Identifier string

	// Destination of events.
	To hcl.Traversal

	// Configuration of the channel.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
