package lang

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"

	"bridgedl/config"
	"bridgedl/config/addr"
)

// BlockReferencesInBody returns all the block references contained in the
// given hcl.Body. The provided Spec is used to infer a schema that allows
// discovering variables in the body.
//
// It is assumed that every hcl.Traversal attribute is a block reference in the
// Bridge Description Language, therefore error diagnostics are returned
// whenever a hcl.Traversal which doesn't match this predicate is encountered.
func BlockReferencesInBody(b hcl.Body, s hcldec.Spec) ([]*addr.Reference, hcl.Diagnostics) {
	return blockReferences(hcldec.Variables(b, s))
}

// blockReferences returns the list of block references that can be parsed from
// the given hcl.Traversals.
func blockReferences(ts []hcl.Traversal) ([]*addr.Reference, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	var refs []*addr.Reference

	for _, t := range ts {
		ref, parseDiags := ParseBlockReference(t)
		diags = diags.Extend(parseDiags)

		if ref != nil {
			refs = append(refs, ref)
		}
	}

	return refs, diags
}

// ParseBlockReference attempts to extract a block reference from a
// hcl.Traversal.
//
// The caller is responsible for checking that a corresponding block exists
// within the Bridge.
func ParseBlockReference(attr hcl.Traversal) (*addr.Reference, hcl.Diagnostics) {
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
		return nil, diags
	}

	ref := &addr.Reference{
		SourceRange: attr.SourceRange(),
	}

	identifier := ts.Rel[0].(hcl.TraverseAttr).Name

	switch blkType {
	case config.BlkChannel:
		ref.Subject = addr.Channel{
			Identifier: identifier,
		}

	case config.BlkRouter:
		ref.Subject = addr.Router{
			Identifier: identifier,
		}

	case config.BlkTransf:
		ref.Subject = addr.Transformer{
			Identifier: identifier,
		}

	case config.BlkTarget:
		ref.Subject = addr.Target{
			Identifier: identifier,
		}

	case config.BlkFunc:
		ref.Subject = addr.Function{
			Identifier: identifier,
		}

	default:
		// should never occur, the list returned by
		// referenceableTypes() is exhaustive
		diags = diags.Append(badRefTypeDiagnostic(blkType, ts.Abs.SourceRange()))
	}

	return ref, diags
}

// referenceableTypes returns a set containing the block types that can be
// referenced inside expressions.
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
