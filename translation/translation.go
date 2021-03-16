package translation

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// Decodable is implemented by component types that embed a HCL configuration
// body which can be decoded into a concrete value.
type Decodable interface {
	// Spec that can be used by facilities of the hcldec package to decode
	// a hcl.Body.
	Spec() hcldec.Spec
}

// Translatable is implemented by component types that can be translated into
// Kubernetes manifests.
//
// An id parameter is required wherever we need to generate deterministic
// object names, because the state of previously generated manifests is not
// persisted.
type Translatable interface {
	// Kubernetes manifests satisfying the given configuration.
	//
	// NOTE(antoineco): we deliberately keep the return type open for the
	// time being to be able to experiment with different formats until
	// this code matures. The initial implementation uses
	// unstructured.Unstructured from k8s/apimachinery.
	Manifests(id string, config, eventDst cty.Value) []interface{}
}

// Addressable is implemented by component types that can receive events from
// other components.
//
// Currently, only KReference destinations are supported, meaning that all
// types that implement Addressable are also Translatable. This may change in
// the future.
// The provided id must be used to generate a deterministic value for the name
// field of this KReference, and that name must match the one of the
// corresponding object generated by the Translatable interface.
type Addressable interface {
	Translatable

	// Address of the component expressed as a Knative "duck" destination
	// in the cty type system.
	Address(id string, config, eventDst cty.Value) cty.Value
}
