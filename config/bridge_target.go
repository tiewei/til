package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// targetBlockSchema is the shallow structure of a "target" block.
// Used for validation during decoding.
var targetBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     attrReplyTo,
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

// decodeTargetBlock performs a partial decoding of the Body of a "target"
// block into a Target struct.
func decodeTargetBlock(blk *hcl.Block) (*Target, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(targetBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[attrReplyTo])
	diags = diags.Extend(decodeDiags)

	trg := &Target{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		ReplyTo:    to,
		Config:     remain,
	}

	return trg, diags
}
