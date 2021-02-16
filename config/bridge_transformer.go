package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// transformerBlockSchema is the shallow structure of a "transformer" block.
// Used for validation during decoding.
var transformerBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     attrTo,
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
	To BlockRef

	// Configuration of the transformer.
	Config hcl.Body
}

// decodeTransformerBlock performs a partial decoding of the Body of a "transformer"
// block into a Transformer struct.
func decodeTransformerBlock(blk *hcl.Block) (*Transformer, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(transformerBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[attrTo])
	diags = diags.Extend(decodeDiags)

	rtr := &Transformer{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		To:         to,
		Config:     remain,
	}

	return rtr, diags
}
