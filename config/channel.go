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

// ChannelBlockSchema is the shallow structure of a "channel" block.
// Used for validation during decoding.
var ChannelBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrTo,
		Required: true,
	}},
}

// Channel represents a generic messaging channel.
type Channel struct {
	// Indicates which type of channel is contained within the block.
	Type string
	// An identifier that is unique among all Channels within a Bridge.
	Identifier string

	// Destination of events.
	To hcl.Traversal

	// Configuration of the channel.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
