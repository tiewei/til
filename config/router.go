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

// RouterBlockSchema is the shallow structure of a "router" block.
// Used for validation during decoding.
var RouterBlockSchema = &hcl.BodySchema{}

// Router represents a generic message router.
type Router struct {
	// Indicates which type of router is contained within the block.
	Type string
	// An identifier that is unique among all Routers within a Bridge.
	Identifier string

	// Configuration of the router.
	Config hcl.Body

	// Source location of the block.
	SourceRange hcl.Range
}
