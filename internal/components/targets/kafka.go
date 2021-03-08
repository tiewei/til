package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/k8s"
	"bridgedl/translation"
)

// TODO(antoineco): this is a mock implementation. It exists only to allow
// testing the generator manually using the config.brg.hcl file at the root of
// the repo while this code is still at the stage of prototype.
type Kafka struct{}

var (
	_ translation.Decodable    = (*Kafka)(nil)
	_ translation.Translatable = (*Kafka)(nil)
	_ translation.Addressable  = (*Kafka)(nil)
)

// Spec implements translation.Decodable.
func (*Kafka) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{}
}

// Manifests implements translation.Translatable.
func (*Kafka) Manifests(id string, config cty.Value) []interface{} {
	panic("not implemented")
}

// Address implements translation.Addressable.
func (*Kafka) Address(id string) cty.Value {
	return k8s.NewDestination("eventing.knative.dev", "KafkaSink", k8s.RFC1123Name(id))
}
