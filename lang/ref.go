package lang

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/hcl/v2"

	"bridgedl/config"
	"bridgedl/config/addr"
)

// BlockReferences takes an instance of a Go type and returns all the block
// references it contains.
//
// This function is typically called with structs decoded from hcl.Body.
func BlockReferences(i interface{}) []*addr.Reference {
	traversals := traversals(reflect.ValueOf(i))

	refs := make([]*addr.Reference, 0, len(traversals))

	for _, tr := range traversals {
		ref, diags := ParseBlockReference(tr)
		if diags.HasErrors() {
			continue
		}
		refs = append(refs, ref)
	}

	return refs
}

// traversals returns a list of all the hcl.Traversals nested inside the given
// reflect.Value.
func traversals(v reflect.Value) []hcl.Traversal {
	switch v.Kind() {
	case reflect.Interface:
		expr, ok := v.Interface().(hcl.Expression)
		if !ok {
			return nil
		}

		tr, diags := hcl.AbsTraversalForExpr(expr)
		if diags.HasErrors() {
			return nil
		}

		return []hcl.Traversal{tr}

	case reflect.Ptr:
		return traversals(v.Elem())

	case reflect.Struct:
		var refs []hcl.Traversal
		for i := 0; i < v.NumField(); i++ {
			refs = append(refs, traversals(v.Field(i))...)
		}
		return refs

	case reflect.Slice:
		var refs []hcl.Traversal
		for i := 0; i < v.Len(); i++ {
			refs = append(refs, traversals(v.Index(i))...)
		}
		return refs

	default:
		return nil
	}
}

// ParseBlockReference attempts to extract a block reference from an expression
// attribute parsed as a hcl.Traversal.
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
		// should never occur
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
