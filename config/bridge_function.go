package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// functionBlockSchema is the shallow structure of a "function" block.
// Used for validation during decoding.
var functionBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     attrReplyTo,
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
}

// decodeFunctionBlock performs a partial decoding of the Body of a "function"
// block into a Function struct.
func decodeFunctionBlock(blk *hcl.Block) (*Function, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(functionBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[attrReplyTo])
	diags = diags.Extend(decodeDiags)

	fn := &Function{
		Identifier: blk.Labels[0],
		ReplyTo:    to,
		Config:     remain,
	}

	return fn, diags
}
