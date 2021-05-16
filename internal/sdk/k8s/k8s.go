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

package k8s

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation"
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

// validateDNS1123Subdomain panics if the given value is not a valid Kubernetes
// object name.
// It is intended to be used inside constructors for throwing loud errors
// whenever component implementations forget to sanitize their input.
func validateDNS1123Subdomain(v string) {
	errs := validation.IsDNS1123Subdomain(v)
	if len(errs) > 0 {
		panic(fmt.Errorf("%q is not a valid Kubernetes object name: %v", v, errs))
	}
}

// Object wraps an instance of unstructured.Unstructured to expose convenience
// field setters.
type Object struct {
	u *unstructured.Unstructured
}

// NewObject returns an Object initialized with the given group/version, kind
// and name.
// Panics if the name parameter is not a valid Kubernetes object name.
func NewObject(apiVersion, kind, name string) *Object {
	validateDNS1123Subdomain(name)

	o := &Object{
		u: &unstructured.Unstructured{},
	}

	o.u.SetAPIVersion(apiVersion)
	o.u.SetKind(kind)
	o.u.SetName(name)

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
