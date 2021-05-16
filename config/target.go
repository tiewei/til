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

// TargetBlockSchema is the shallow structure of a "target" block.
// Used for validation during decoding.
var TargetBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrReplyTo,
		Required: false,
	}},
}

// Target represents a generic event target.
type Target struct {
	// Indicates which type of target is contained within the block.
	Type string
	// An identifier that is unique among all Targets within a Bridge.
	Identifier string

	// Destination of event responses.
	ReplyTo hcl.Traversal

	// Configuration of the target.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
