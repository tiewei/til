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

package config

import "github.com/hashicorp/hcl/v2"

// HCL blocks supported in a Bridge Description File.
const (
	BlkBridge  = "bridge"
	BlkChannel = "channel"
	BlkRouter  = "router"
	BlkTransf  = "transformer"
	BlkSource  = "source"
	BlkTarget  = "target"
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
		Type:       BlkBridge,
		LabelNames: []string{LblID},
	}, {
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
	}},
}

// BridgeBlockSchema is the shallow structure of the "bridge" block, which is
// unique per Bridge.
// Used for validation during decoding.
var BridgeBlockSchema = &hcl.BodySchema{}

// Bridge represents the body of a Bridge Description File.
type Bridge struct {
	// Absolute path of the file this configuration was loaded from.
	Path string

	// Bridge globals.
	Identifier string

	// Indexed lists of messaging components.
	// Parsers should index each component with a key that uniquely identifies a block.
	Channels     map[interface{}]*Channel
	Routers      map[interface{}]*Router
	Transformers map[interface{}]*Transformer
	Sources      map[interface{}]*Source
	Targets      map[interface{}]*Target
}
