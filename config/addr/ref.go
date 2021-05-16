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

package addr

import "github.com/hashicorp/hcl/v2"

// Referenceable must be implemented by all address types that can be used as
// references in Bridge descriptions.
type Referenceable interface {
	// Addr produces a string representation of the component by which it
	// can be uniquely identified and referenced via HCL expressions.
	Addr() string
}

// Reference wraps a Referenceable together with the source location of the
// expression that references it.
type Reference struct {
	// Actual type being referenced.
	Subject Referenceable
	// Source location of the referencer.
	SourceRange hcl.Range
}
