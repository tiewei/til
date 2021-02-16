package config

import "github.com/hashicorp/hcl/v2"

// HCL blocks supported in a Bridge Description File.
const (
	BlkChannel = "channel"
	BlkRouter  = "router"
	BlkTransf  = "transformer"
	BlkSource  = "source"
	BlkTarget  = "target"
	BlkFunc    = "function"
)

// Common identifiers for HCL block labels.
const (
	lblType = "type"
	lblID   = "identifier"
)

// Common block attributes.
const (
	attrTo      = "to"
	attrReplyTo = "reply_to"
)

// bridgeSchema is the shallow structure of a Bridge Description File.
// Used for validation during decoding.
var bridgeSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{{
		Type:       BlkChannel,
		LabelNames: []string{lblType, lblID},
	}, {
		Type:       BlkRouter,
		LabelNames: []string{lblType, lblID},
	}, {
		Type:       BlkTransf,
		LabelNames: []string{lblType, lblID},
	}, {
		Type:       BlkSource,
		LabelNames: []string{lblType, lblID},
	}, {
		Type:       BlkTarget,
		LabelNames: []string{lblType, lblID},
	}, {
		Type:       BlkFunc,
		LabelNames: []string{lblID},
	}},
}

// Bridge represents the body of a Bridge description file.
type Bridge struct {
	// Messaging components
	Channels     []*Channel
	Routers      []*Router
	Transformers []*Transformer
	Sources      []*Source
	Targets      []*Target
	Functions    []*Function
}

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
func decodeBridge(b hcl.Body, brg *Bridge) hcl.Diagnostics {
	var diags hcl.Diagnostics

	content, contentDiags := b.Content(bridgeSchema)
	diags = diags.Extend(contentDiags)

	for _, blk := range content.Blocks {
		switch t := blk.Type; t {
		case BlkChannel:
			ch, decodeDiags := decodeChannelBlock(blk)
			diags = diags.Extend(decodeDiags)

			if ch != nil {
				brg.Channels = append(brg.Channels, ch)
			}

		case BlkRouter:
			rtr, decodeDiags := decodeRouterBlock(blk)
			diags = diags.Extend(decodeDiags)

			if rtr != nil {
				brg.Routers = append(brg.Routers, rtr)
			}

		case BlkTransf:
			trsf, decodeDiags := decodeTransformerBlock(blk)
			diags = diags.Extend(decodeDiags)

			if trsf != nil {
				brg.Transformers = append(brg.Transformers, trsf)
			}

		case BlkSource:
			src, decodeDiags := decodeSourceBlock(blk)
			diags = diags.Extend(decodeDiags)

			if src != nil {
				brg.Sources = append(brg.Sources, src)
			}

		case BlkTarget:
			trg, decodeDiags := decodeTargetBlock(blk)
			diags = diags.Extend(decodeDiags)

			if trg != nil {
				brg.Targets = append(brg.Targets, trg)
			}

		case BlkFunc:
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
