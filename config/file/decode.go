/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package file

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"til/config"
	"til/config/addr"
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

	visitedBridgeGlobals := false

	for _, blk := range content.Blocks {
		switch t := blk.Type; t {
		case config.BlkBridge:
			if visitedBridgeGlobals {
				diags = diags.Append(tooManyGlobalBlocksDiagnostic(t, blk.DefRange))
			}
			visitedBridgeGlobals = true

			setDiags := setGlobals(brg, blk)
			diags = diags.Extend(setDiags)

		case config.BlkChannel:
			addDiags := addChannelBlock(brg, blk)
			diags = diags.Extend(addDiags)

		case config.BlkRouter:
			addDiags := addRouterBlock(brg, blk)
			diags = diags.Extend(addDiags)

		case config.BlkTransf:
			addDiags := addTransformerBlock(brg, blk)
			diags = diags.Extend(addDiags)

		case config.BlkSource:
			addDiags := addSourceBlock(brg, blk)
			diags = diags.Extend(addDiags)

		case config.BlkTarget:
			addDiags := addTargetBlock(brg, blk)
			diags = diags.Extend(addDiags)

		default:
			// should never occur because the hcl.BodyContent was
			// validated against a hcl.BodySchema during parsing
			panic(fmt.Sprintf("found unexpected block type %q. The HCL body schema is outdated.", t))
		}
	}

	return diags
}

// setGlobals sets the global properties and options of a Bridge.
func setGlobals(brg *config.Bridge, blk *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0]))
	}

	content, contentDiags := blk.Body.Content(config.BridgeBlockSchema)
	diags = diags.Extend(contentDiags)

	var delivery *config.Delivery
	visitedDelivery := false

	for _, blk := range content.Blocks {
		switch t := blk.Type; t {
		case config.BlkDelivery:
			if visitedDelivery {
				diags = diags.Append(tooManyGlobalBlocksDiagnostic(t, blk.DefRange))
			}
			visitedDelivery = true

			var decodeDiags hcl.Diagnostics
			delivery, decodeDiags = decodeBridgeDeliveryBlock(blk)
			diags = diags.Extend(decodeDiags)
		}
	}

	brg.Identifier = blk.Labels[0]
	brg.Delivery = delivery

	return diags
}

// addChannelBlock adds a Channel component to a Bridge.
func addChannelBlock(brg *config.Bridge, blk *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics

	ch, decodeDiags := decodeChannelBlock(blk)
	diags = diags.Extend(decodeDiags)

	if ch == nil {
		return diags
	}

	if brg.Channels == nil {
		brg.Channels = make(map[interface{}]*config.Channel)
	}

	key := addr.Channel{Identifier: ch.Identifier}

	if _, exists := brg.Channels[key]; exists {
		diags = diags.Append(duplicateBlockDiagnostic(config.CategoryChannels, ch.Identifier, blk.DefRange))
	} else {
		brg.Channels[key] = ch
	}

	return diags
}

// addRouterBlock adds a Router component to a Bridge.
func addRouterBlock(brg *config.Bridge, blk *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics

	rtr, decodeDiags := decodeRouterBlock(blk)
	diags = diags.Extend(decodeDiags)

	if rtr == nil {
		return nil
	}

	if brg.Routers == nil {
		brg.Routers = make(map[interface{}]*config.Router)
	}

	key := addr.Router{Identifier: rtr.Identifier}

	if _, exists := brg.Routers[key]; exists {
		diags = diags.Append(duplicateBlockDiagnostic(config.CategoryRouters, rtr.Identifier, blk.DefRange))
	} else {
		brg.Routers[key] = rtr
	}

	return diags
}

// addTransformerBlock adds a Transformer component to a Bridge.
func addTransformerBlock(brg *config.Bridge, blk *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics

	trsf, decodeDiags := decodeTransformerBlock(blk)
	diags = diags.Extend(decodeDiags)

	if trsf == nil {
		return nil
	}

	if brg.Transformers == nil {
		brg.Transformers = make(map[interface{}]*config.Transformer)
	}

	key := addr.Transformer{Identifier: trsf.Identifier}

	if _, exists := brg.Transformers[key]; exists {
		diags = diags.Append(duplicateBlockDiagnostic(config.CategoryTransformers, trsf.Identifier, blk.DefRange))
	} else {
		brg.Transformers[key] = trsf
	}

	return diags
}

// addSourceBlock adds a Source component to a Bridge.
func addSourceBlock(brg *config.Bridge, blk *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics

	src, decodeDiags := decodeSourceBlock(blk)
	diags = diags.Extend(decodeDiags)

	if src == nil {
		return nil
	}

	if brg.Sources == nil {
		brg.Sources = make(map[interface{}]*config.Source)
	}

	key := addr.Source{Identifier: src.Identifier}

	if _, exists := brg.Sources[key]; exists {
		diags = diags.Append(duplicateBlockDiagnostic(config.CategorySources, src.Identifier, blk.DefRange))
	} else {
		brg.Sources[key] = src
	}

	return diags
}

// addTargetBlock adds a Target component to a Bridge.
func addTargetBlock(brg *config.Bridge, blk *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics

	trg, decodeDiags := decodeTargetBlock(blk)
	diags = diags.Extend(decodeDiags)

	if trg == nil {
		return nil
	}

	if brg.Targets == nil {
		brg.Targets = make(map[interface{}]*config.Target)
	}

	key := addr.Target{Identifier: trg.Identifier}

	if _, exists := brg.Targets[key]; exists {
		diags = diags.Append(duplicateBlockDiagnostic(config.CategoryTargets, trg.Identifier, blk.DefRange))
	} else {
		brg.Targets[key] = trg
	}

	return diags
}

// decodeBridgeDeliveryBlock performs a decoding of the Body of a
// "bridge.delivery" block into a Delivery struct.
func decodeBridgeDeliveryBlock(blk *hcl.Block) (*config.Delivery, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	content, contentDiags := blk.Body.Content(config.DeliveryBlockSchema)
	diags = diags.Extend(contentDiags)

	retries, decodeDiags := decodeInt64Val(content.Attributes[config.AttrRetries])
	diags = diags.Extend(decodeDiags)

	dls, decodeDiags := decodeBlockRef(content.Attributes[config.AttrDeadLetterSink])
	diags = diags.Extend(decodeDiags)

	d := &config.Delivery{
		Retries:        retries,
		DeadLetterSink: dls,
	}

	return d, diags
}

// decodeChannelBlock performs a partial decoding of the Body of a "channel"
// block into a Channel struct.
func decodeChannelBlock(blk *hcl.Block) (*config.Channel, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0]))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1]))
	}

	_, remain, contentDiags := blk.Body.PartialContent(config.ChannelBlockSchema)
	diags = diags.Extend(contentDiags)

	ch := &config.Channel{
		Type:        blk.Labels[0],
		Identifier:  blk.Labels[1],
		Config:      remain,
		SourceRange: blk.DefRange,
	}

	return ch, diags
}

// decodeRouterBlock performs a partial decoding of the Body of a "router"
// block into a Router struct.
func decodeRouterBlock(blk *hcl.Block) (*config.Router, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0]))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1]))
	}

	_, remain, contentDiags := blk.Body.PartialContent(config.RouterBlockSchema)
	diags = diags.Extend(contentDiags)

	rtr := &config.Router{
		Type:        blk.Labels[0],
		Identifier:  blk.Labels[1],
		Config:      remain,
		SourceRange: blk.DefRange,
	}

	return rtr, diags
}

// decodeTransformerBlock performs a partial decoding of the Body of a "transformer"
// block into a Transformer struct.
func decodeTransformerBlock(blk *hcl.Block) (*config.Transformer, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0]))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1]))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.TransformerBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrTo])
	diags = diags.Extend(decodeDiags)

	rtr := &config.Transformer{
		Type:        blk.Labels[0],
		Identifier:  blk.Labels[1],
		To:          to,
		Config:      remain,
		SourceRange: blk.DefRange,
	}

	return rtr, diags
}

// decodeSourceBlock performs a partial decoding of the Body of a "source"
// block into a Source struct.
func decodeSourceBlock(blk *hcl.Block) (*config.Source, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0]))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1]))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.SourceBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrTo])
	diags = diags.Extend(decodeDiags)

	src := &config.Source{
		Type:        blk.Labels[0],
		Identifier:  blk.Labels[1],
		To:          to,
		Config:      remain,
		SourceRange: blk.DefRange,
	}

	return src, diags
}

// decodeTargetBlock performs a partial decoding of the Body of a "target"
// block into a Target struct.
func decodeTargetBlock(blk *hcl.Block) (*config.Target, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	if !hclsyntax.ValidIdentifier(blk.Labels[0]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[0]))
	}
	if !hclsyntax.ValidIdentifier(blk.Labels[1]) {
		diags = diags.Append(badIdentifierDiagnostic(blk.LabelRanges[1]))
	}

	content, remain, contentDiags := blk.Body.PartialContent(config.TargetBlockSchema)
	diags = diags.Extend(contentDiags)

	to, decodeDiags := decodeBlockRef(content.Attributes[config.AttrReplyTo])
	diags = diags.Extend(decodeDiags)

	trg := &config.Target{
		Type:        blk.Labels[0],
		Identifier:  blk.Labels[1],
		ReplyTo:     to,
		Config:      remain,
		SourceRange: blk.DefRange,
	}

	return trg, diags
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

// decodeInt64Val decodes an integer attribute.
func decodeInt64Val(attr *hcl.Attribute) (*int64, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	var out *int64

	if attr == nil {
		return out, diags
	}

	val, evalDiags := attr.Expr.Value(nil)
	diags = diags.Extend(evalDiags)

	if isInt64(val) {
		out = new(int64)
		*out, _ = val.AsBigFloat().Int64()
	} else {
		diags = diags.Append(wrongTypeDiagnostic(val, "integer", attr.Expr.Range()))
	}

	return out, diags
}
