package k8s

import "github.com/zclconf/go-cty/cty"

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

// SecretKeySelectorCty is a non-primitive cty.Type that represents a Kubernetes
// corev1.SecretKeySelector.
//
// It can be used to indicate that the value of an environment variable should
// be read from a Kubernetes Secret.
var SecretKeySelectorCty = cty.Object(map[string]cty.Type{
	"name": cty.String,
	"key":  cty.String,
})

// NewSecretKeySelector returns a new Kubernetes corev1.SecretKeySelector as a
// cty.Value which satisfies the SecretKeySelectorCty type.
func NewSecretKeySelector(name, key string) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"name": cty.StringVal(name),
		"key":  cty.StringVal(key),
	})
}

// IsSecretKeySelector verifies that the given cty.Value conforms to the
// SecretKeySelectorCty type.
func IsSecretKeySelector(v cty.Value) bool {
	return v.Type().Equals(SecretKeySelectorCty)
}
