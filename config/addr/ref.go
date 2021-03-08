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
