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

// HCL blocks supported in global settings ("bridge" block).
const (
	BlkDelivery = "delivery"
)

// Block attributes that can appear in global settings (sub-blocks of "bridge" block).
const (
	AttrRetries        = "retries"
	AttrDeadLetterSink = "dead_letter_sink"
)

// BridgeBlockSchema is the shallow structure of a "bridge" block.
// There can be at most one such block declared inside a Bridge.
// Used for validation during decoding.
var BridgeBlockSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{{
		Type: BlkDelivery,
	}},
}

// DeliveryBlockSchema is the shallow structure of a "bridge.delivery" block.
// Used for validation during decoding.
var DeliveryBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrRetries,
		Required: false,
	}, {
		Name:     AttrDeadLetterSink,
		Required: false,
	}},
}

// Delivery represents the message delivery options which apply to the Bridge.
type Delivery struct {
	Retries        *int64
	DeadLetterSink hcl.Traversal
}
