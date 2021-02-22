package config

import "github.com/hashicorp/hcl/v2"

// decodeBlockRef decodes an expression attribute representing a reference to
// another block.
//
// References are expected to be in the format "block_type.identifier".
// Those attributes are ultimately used as static references to construct a
// graph of connections between messaging components within the Bridge.
// In the parsing/decoding phase, we are not interested in evaluating or
// validating their value (the actual hcl.Block that is referenced).
//
// Under the hood, those expressions are interpreted as a hcl.Traversal which,
// in terms of HCL language definition, is a variable name followed by zero or
// more attributes or index operators with constant operands (e.g. "foo",
// "foo.bar", "foo[0]").
func decodeBlockRef(attr *hcl.Attribute) (hcl.Traversal, hcl.Diagnostics) {
	if attr == nil {
		return nil, nil
	}

	return hcl.AbsTraversalForExpr(attr.Expr)
}
