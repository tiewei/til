package config

import "github.com/hashicorp/hcl/v2"

// BlockRef represents a reference to another block.
//
// References are expected to be expression attributes in the format "block_type.identifier".
// Under the hood, those expressions are interpreted as a hcl.Traversal, which
// is a variable name followed by zero or more attribute access or index
// operators with constant operands (e.g. "foo", "foo.bar", "foo[0]").
type BlockRef struct {
	BlockType  string
	Identifier string
}

// decodeBlockRef decodes an expression attribute into a BlockRef.
//
// BlockRefs attributes are ultimately used as static references to construct a
// graph of connections between messaging components within the Bridge.
// Therefore, we are not interested in evaluating their value (the actual
// hcl.Block that is referenced).
func decodeBlockRef(attr *hcl.Attribute) (BlockRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	var blkRef BlockRef

	if attr == nil {
		return blkRef, diags
	}

	traversal, travDiags := hcl.AbsTraversalForExpr(attr.Expr)
	diags = diags.Extend(travDiags)

	if traversal == nil {
		return blkRef, diags
	}

	switch blkType := traversal.RootName(); blkType {
	case BlkChannel, BlkRouter, BlkTransf, BlkTarget, BlkFunc:
		rel := traversal.SimpleSplit().Rel
		if len(rel) != 1 {
			diags = diags.Append(badReferenceDiagnostic(attr.Expr))
			break
		}

		blkRef.BlockType = blkType
		blkRef.Identifier = rel[0].(hcl.TraverseAttr).Name

	default:
		diags = diags.Append(badReferenceDiagnostic(attr.Expr))
	}

	return blkRef, diags
}
