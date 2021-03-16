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
	nestedSpec := &hcldec.BlockListSpec{
		TypeName: "operation",
		Nested: &hcldec.ObjectSpec{
			"operation": &hcldec.BlockLabelSpec{
				Index: 0,
				Name:  "operation",
			},
			"path": &hcldec.BlockSpec{
				TypeName: "path",
				Nested: &hcldec.ObjectSpec{
					"key": &hcldec.AttrSpec{
						Name:     "key",
						Type:     cty.String,
						Required: true,
					},
					"value": &hcldec.AttrSpec{
						Name:     "value",
						Type:     cty.String,
						Required: true,
					},
				},
				Required: true,
			},
		},
		MinItems: 1,
	}

	return &hcldec.ObjectSpec{
		"context": &hcldec.BlockSpec{
			TypeName: "context",
			Nested:   nestedSpec,
		},
		"data": &hcldec.BlockSpec{
			TypeName: "data",
			Nested:   nestedSpec,
		},
	}

	/*
		Example of value decoded from the spec above, for a "context"
		and a "data" block, both containing two nested "operation" blocks:

		v: (map[string]interface {}) (len=2) {
		 "context": ([]interface {}) (len=2) {
		  (map[string]interface {}) (len=2) {
		   "operation": (string) "store",
		   "path": (map[string]interface {}) (len=2) {
		    "key": (string) "$id",
		    "value": (string) "id"
		   }
		  },
		  (map[string]interface {}) (len=2) {
		   "operation": (string) "add",
		   "path": (map[string]interface {}) (len=2) {
		    "key": (string) "id",
		    "value": (string) "${person}-${id}"
		   }
		  }
		 },
		 "data": ([]interface {}) (len=2) { ... }
		}
	*/
}

// Manifests implements translation.Translatable.
func (*Bumblebee) Manifests(id string, config, eventDst cty.Value) []interface{} {
	return nil
}

// Address implements translation.Addressable.
func (*Bumblebee) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("flow.triggermesh.io/v1alpha1", "Transformation", k8s.RFC1123Name(id))
}
