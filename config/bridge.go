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
	LblType = "type"
	LblID   = "identifier"
)

// Common block attributes.
const (
	AttrTo      = "to"
	AttrReplyTo = "reply_to"
)

// BridgeSchema is the shallow structure of a Bridge Description File.
// Used for validation during decoding.
var BridgeSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{{
		Type:       BlkChannel,
		LabelNames: []string{LblType, LblID},
	}, {
		Type:       BlkRouter,
		LabelNames: []string{LblType, LblID},
	}, {
		Type:       BlkTransf,
		LabelNames: []string{LblType, LblID},
	}, {
		Type:       BlkSource,
		LabelNames: []string{LblType, LblID},
	}, {
		Type:       BlkTarget,
		LabelNames: []string{LblType, LblID},
	}, {
		Type:       BlkFunc,
		LabelNames: []string{LblID},
	}},
}

// Bridge represents the body of a Bridge Description File.
type Bridge struct {
	// Messaging components
	Channels     []*Channel
	Routers      []*Router
	Transformers []*Transformer
	Sources      []*Source
	Targets      []*Target
	Functions    []*Function
}
