package k8s

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
)

// RFC1123Name sanitizes the given input string, ensuring it is a valid DNS
// subdomain name (as defined in RFC 1123) which can be used as a Kubernetes
// object name.
//
// The implementation is purposedly naive because it is assumed that this
// function will only ever sanitize HCL identifiers, which already contain a
// limited set of characters (see hclsyntax.ValidIdentifier).
func RFC1123Name(id string) string {
	return strings.ToLower(strings.ReplaceAll(id, "_", "-"))
}

// ObjectReferenceCty is a non-primitive cty.Type that represents a Kubernetes
// corev1.LocalObjectReference.
//
// It can be used wherever a simple name-reference to a known Kubernetes object
// is appropriate, such as in HCL attributes which are used to populate
// references to Kubernetes Secrets ("secretKeyRef") but the implementer is
// assumed to have some implicit knowledge about the keys which that Secret
// should contain.
var ObjectReferenceCty = cty.Object(map[string]cty.Type{
	"name": cty.String,
})

// NewObjectReference returns a new Kubernetes corev1.LocalObjectReference as a
// cty.Value which satisfies the ObjectReferenceCty type.
func NewObjectReference(name string) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"name": cty.StringVal(name),
	})
}

// IsObjectReference verifies that the given cty.Value conforms to the
// ObjectReferenceCty type.
func IsObjectReference(v cty.Value) bool {
	return v.Type().Equals(ObjectReferenceCty)
}
