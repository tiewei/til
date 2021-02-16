package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// sourceBlockSchema is the shallow structure of a "source" block.
// Used for validation during decoding.
var sourceBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     attrTo,
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
	To BlockRef

	// Configuration of the source.
	Config hcl.Body
}

// decodeSourceBlock performs a partial decoding of the Body of a "source"
// block into a Source struct.
func decodeSourceBlock(blk *hcl.Block) (*Source, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(sourceBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[attrTo])
	diags = diags.Extend(decodeDiags)

	src := &Source{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		To:         to,
		Config:     remain,
	}

	return src, diags
}
