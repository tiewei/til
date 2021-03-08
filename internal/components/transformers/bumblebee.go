package transformers

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/k8s"
	"bridgedl/translation"
)

// TODO(antoineco): this is a mock implementation. It exists only to allow
// testing the generator manually using the config.brg.hcl file at the root of
// the repo while this code is still at the stage of prototype.
type Bumblebee struct{}

var (
	_ translation.Decodable    = (*Bumblebee)(nil)
	_ translation.Translatable = (*Bumblebee)(nil)
	_ translation.Addressable  = (*Bumblebee)(nil)
)

// Spec implements translation.Decodable.
func (*Bumblebee) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{}
}

// Manifests implements translation.Translatable.
func (*Bumblebee) Manifests(id string, config cty.Value) []interface{} {
	panic("not implemented")
}

// Address implements translation.Addressable.
func (*Bumblebee) Address(id string) cty.Value {
	return k8s.NewDestination("flow.triggermesh.io/v1alpha1", "Transformation", k8s.RFC1123Name(id))
}
