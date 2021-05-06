package k8s

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

// Object wraps an instance of unstructured.Unstructured to expose convenience
// field setters.
type Object struct {
	u *unstructured.Unstructured
}

// NewObject returns an Object initialized with the given group/version, kind
// and identifier.
func NewObject(apiVersion, kind, id string) *Object {
	o := &Object{
		u: &unstructured.Unstructured{},
	}

	o.u.SetAPIVersion(apiVersion)
	o.u.SetKind(kind)
	o.u.SetName(RFC1123Name(id))

	return o
}

// Unstructured returns the underlying unstructured.Unstructured.
func (o *Object) Unstructured() *unstructured.Unstructured {
	return o.u
}

// SetNestedField invokes unstructured.SetNestedField on the underlying
// object's data.
func (o *Object) SetNestedField(value interface{}, fields ...string) {
	err := unstructured.SetNestedField(o.u.Object, value, fields...)
	if err != nil {
		panic(err)
	}
}

// SetNestedSlice invokes unstructured.SetNestedSlice on the underlying
// object's data.
func (o *Object) SetNestedSlice(value []interface{}, fields ...string) {
	err := unstructured.SetNestedSlice(o.u.Object, value, fields...)
	if err != nil {
		panic(err)
	}
}

// SetNestedMap invokes unstructured.SetNestedMap on the underlying
// object's data.
func (o *Object) SetNestedMap(value map[string]interface{}, fields ...string) {
	err := unstructured.SetNestedMap(o.u.Object, value, fields...)
	if err != nil {
		panic(err)
	}
}
