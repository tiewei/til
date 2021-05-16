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

// TransformerBlockSchema is the shallow structure of a "transformer" block.
// Used for validation during decoding.
var TransformerBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{{
		Name:     AttrTo,
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
	To hcl.Traversal

	// Configuration of the transformer.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
