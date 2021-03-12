package functions

import (
	"github.com/zclconf/go-cty/cty"

	"bridgedl/k8s"
	"bridgedl/translation"
)

// "function" types are not supported yet, so we use this type as a stub,
// assuming there is only a single way to address and translate a function.
type Function struct{}

var (
	_ translation.Translatable = (*Function)(nil)
	_ translation.Addressable  = (*Function)(nil)
)

// Manifests implements translation.Translatable.
func (*Function) Manifests(id string, config, eventDst cty.Value) []interface{} {
	return nil
}

// Address implements translation.Addressable.
func (*Function) Address(id string) cty.Value {
	return k8s.NewDestination("serving.knative.dev/v1", "Service", k8s.RFC1123Name(id))
}
