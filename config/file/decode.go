package file

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"bridgedl/config"
)

// decodeBridge performs a partial decoding of the Body of a Bridge Description
// File into a Bridge struct.
//
// Because a Bridge struct is only a representation of the code contained
// inside a Bridge Description File:
//   * fields of type hcl.Body are left to be decoded into their own type based
//     on the labels and attributes read in this partial decoding.
//   * fields of type hcl.Expression are left to be evaluated against a
//     hcl.EvalContext if necessary.
//
// This function leverages the low-level HCL APIs instead of calling
// gohcl.DecodeBody() in order to have better control over the decoding and
// validation of the configuration blocks, and over the contents of error
// diagnostics.
func decodeBridge(b hcl.Body, brg *config.Bridge) hcl.Diagnostics {
	var diags hcl.Diagnostics

	content, contentDiags := b.Content(config.BridgeSchema)
	diags = diags.Extend(contentDiags)

	for _, blk := range content.Blocks {
		switch t := blk.Type; t {
		case config.BlkChannel:
			ch, decodeDiags := decodeChannelBlock(blk)
			diags = diags.Extend(decodeDiags)

			if ch != nil {
				brg.Channels = append(brg.Channels, ch)
			}

		case config.BlkRouter:
			rtr, decodeDiags := decodeRouterBlock(blk)
			diags = diags.Extend(decodeDiags)

			if rtr != nil {
				brg.Routers = append(brg.Routers, rtr)
			}

		case config.BlkTransf:
			trsf, decodeDiags := decodeTransformerBlock(blk)
			diags = diags.Extend(decodeDiags)

			if trsf != nil {
				brg.Transformers = append(brg.Transformers, trsf)
			}

		case config.BlkSource:
			src, decodeDiags := decodeSourceBlock(blk)
			diags = diags.Extend(decodeDiags)

			if src != nil {
				brg.Sources = append(brg.Sources, src)
			}

		case config.BlkTarget:
			trg, decodeDiags := decodeTargetBlock(blk)
			diags = diags.Extend(decodeDiags)

			if trg != nil {
				brg.Targets = append(brg.Targets, trg)
			}

		case config.BlkFunc:
			fn, decodeDiags := decodeFunctionBlock(blk)
			diags = diags.Extend(decodeDiags)

			if fn != nil {
				brg.Functions = append(brg.Functions, fn)
			}

		default:
			// should never occur because the hcl.BodyContent was
			// validated against a hcl.BodySchema
			continue
		}
	}

	return diags
}

// decodeChannelBlock performs a partial decoding of the Body of a "channel"
// block into a Channel struct.
func decodeChannelBlock(blk *hcl.Block) (*config.Channel, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.ChannelBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrTo])
	diags = diags.Extend(decodeDiags)

	ch := &config.Channel{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		To:         to,
		Config:     remain,
	}

	return ch, diags
}

// decodeFunctionBlock performs a partial decoding of the Body of a "function"
// block into a Function struct.
func decodeFunctionBlock(blk *hcl.Block) (*config.Function, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.FunctionBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrReplyTo])
	diags = diags.Extend(decodeDiags)

	fn := &config.Function{
		Identifier: blk.Labels[0],
		ReplyTo:    to,
		Config:     remain,
	}

	return fn, diags
}

// decodeRouterBlock performs a partial decoding of the Body of a "router"
// block into a Router struct.
func decodeRouterBlock(blk *hcl.Block) (*config.Router, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	_, remain, contentDiags := blk.Body.PartialContent(config.RouterBlockSchema)
	diags = diags.Extend(contentDiags)

	rtr := &config.Router{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		Config:     remain,
	}

	return rtr, diags
}

// decodeSourceBlock performs a partial decoding of the Body of a "source"
// block into a Source struct.
func decodeSourceBlock(blk *hcl.Block) (*config.Source, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.SourceBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrTo])
	diags = diags.Extend(decodeDiags)

	src := &config.Source{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		To:         to,
		Config:     remain,
	}

	return src, diags
}

// decodeTargetBlock performs a partial decoding of the Body of a "target"
// block into a Target struct.
func decodeTargetBlock(blk *hcl.Block) (*config.Target, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.TargetBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrReplyTo])
	diags = diags.Extend(decodeDiags)

	trg := &config.Target{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		ReplyTo:    to,
		Config:     remain,
	}

	return trg, diags
}

// decodeTransformerBlock performs a partial decoding of the Body of a "transformer"
// block into a Transformer struct.
func decodeTransformerBlock(blk *hcl.Block) (*config.Transformer, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0].Ptr()))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1].Ptr()))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.TransformerBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrTo])
	diags = diags.Extend(decodeDiags)

	rtr := &config.Transformer{
		Type:       blk.Labels[0],
		Identifier: blk.Labels[1],
		To:         to,
		Config:     remain,
	}

	return rtr, diags
}

// decodeBlockRef decodes an expression attribute representing a reference to
// another block.
//
// References are expected to be in the format "block_type.identifier".
// Those attributes are ultimately used as static references to construct a
// graph of connections between messaging components within the Bridge.
// In the parsing/decoding phase, we are not interested in evaluating or
// validating their value (the actual hcl.Block that is referenced).
//
// Under the hood, those expressions are interpreted as a hcl.Traversal which,
// in terms of HCL language definition, is a variable name followed by zero or
// more attributes or index operators with constant operands (e.g. "foo",
// "foo.bar", "foo[0]").
func decodeBlockRef(attr *hcl.Attribute) (hcl.Traversal, hcl.Diagnostics) {
	if attr == nil {
		return nil, nil
	}

	return hcl.AbsTraversalForExpr(attr.Expr)
}
