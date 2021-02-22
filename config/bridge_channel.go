package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// channelBlockSchema is the shallow structure of a "channel" block.
// Used for validation during decoding.
var channelBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     attrTo,
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
}

// decodeChannelBlock performs a partial decoding of the Body of a "channel"
// block into a Channel struct.
func decodeChannelBlock(blk *hcl.Block) (*Channel, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(channelBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[attrTo])
	diags = diags.Extend(decodeDiags)

	ch := &Channel{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		To:         to,
		Config:     remain,
	}

	return ch, diags
}
