package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/k8s"
	"bridgedl/translation"
)

// TODO(antoineco): this is a mock implementation. It exists only to allow
// testing the generator manually using the config.brg.hcl file at the root of
// the repo while this code is still at the stage of prototype.
type AWSSQS struct{}

var (
	_ translation.Decodable    = (*AWSSQS)(nil)
	_ translation.Translatable = (*AWSSQS)(nil)
	_ translation.Addressable  = (*AWSSQS)(nil)
)

// Spec implements translation.Decodable.
func (*AWSSQS) Spec() hcldec.Spec {
	return &hcldec.AttrSpec{
		Name:     "arn",
		Type:     cty.String,
		Required: true,
	}
}

// Manifests implements translation.Translatable.
func (*AWSSQS) Manifests(id string, config, eventDst cty.Value) []interface{} {
	return nil
}

// Address implements translation.Addressable.
func (*AWSSQS) Address(id string) cty.Value {
	return k8s.NewDestination("sources.triggermesh.io/v1alpha1", "AWSSQSSource", k8s.RFC1123Name(id))
}
