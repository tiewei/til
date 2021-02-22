package addr

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
)

// Referenceable must be implemented by all address types that can be used as
// references in Bridge descriptions.
type Referenceable interface {
	// Addr returns a string representation of the Referenceable as an address
	Addr() string
}

// Reference wraps a Referenceable together with the source location of the
// expression that references it.
type Reference struct {
	// Actual type being referenced
	Subject Referenceable
	// Source location of the referencer
	SourceRange hcl.Range
}

// ParseBlockRef attempt to extracts a Reference from an expression attribute
// parsed as a hcl.Traversal.
func ParseBlockRef(attr hcl.Traversal) (*Reference, hcl.Diagnostics) {
	if attr == nil {
		return nil, nil
	}

	var diags hcl.Diagnostics

	ts := attr.SimpleSplit()
	blkType := ts.RootName()

	if !referenceableTypes().Has(blkType) {
		diags = diags.Append(badRefTypeDiagnostic(blkType, ts.Abs.SourceRange()))
		return nil, diags
	}

	if len(ts.Rel) != 1 {
		diags = diags.Append(badRefFormatDiagnostic(attr.SourceRange()))
		diags = diags.Append(badRefFormatDiagnostic(ts.Abs.SourceRange()))
		return nil, diags
	}

	ref := &Reference{
		SourceRange: attr.SourceRange(),
	}

	identifier := ts.Rel[0].(hcl.TraverseAttr).Name

	switch blkType {
	case config.BlkChannel:
		ref.Subject = Channel{
			Identifier: identifier,
		}

	case config.BlkRouter:
		ref.Subject = Router{
			Identifier: identifier,
		}

	case config.BlkTransf:
		ref.Subject = Transformer{
			Identifier: identifier,
		}

	case config.BlkTarget:
		ref.Subject = Target{
			Identifier: identifier,
		}

	case config.BlkFunc:
		ref.Subject = Function{
			Identifier: identifier,
		}

	default:
		// should never occur
		diags = diags.Append(badRefTypeDiagnostic(blkType, ts.Abs.SourceRange()))
	}

	return ref, diags
}

// referenceableTypes returns a set containing the block types that can be
// referenced from expressions.
func referenceableTypes() stringSet {
	var refTypes stringSet

	refTypes.Add(
		config.BlkChannel,
		config.BlkRouter,
		config.BlkTransf,
		config.BlkTarget,
		config.BlkFunc,
	)

	return refTypes
}

type stringSet map[string]struct{}

var _ fmt.Stringer = (stringSet)(nil)

// Add adds or replaces elements in the set.
func (s *stringSet) Add(elems ...string) {
	if len(elems) == 0 {
		return
	}

	if *s == nil {
		*s = make(stringSet, len(elems))
	}

	for _, e := range elems {
		(*s)[e] = struct{}{}
	}
}

// Has returns whether the set contains the given element.
func (s stringSet) Has(elem string) bool {
	if s == nil {
		return false
	}

	_, contains := s[elem]
	return contains
}

// String implements fmt.Stringer.
func (s stringSet) String() string {
	elems := make([]string, 0, len(s))

	for e := range s {
		elems = append(elems, e)
	}

	sort.Strings(elems)

	return fmt.Sprintf("%q", elems)
}
