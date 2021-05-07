package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
	"bridgedl/translation"
)

type Infra struct{}

var (
	_ translation.Decodable    = (*Infra)(nil)
	_ translation.Translatable = (*Infra)(nil)
	_ translation.Addressable  = (*Infra)(nil)
)

// Spec implements translation.Decodable.
func (*Infra) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"code": &hcldec.AttrSpec{
			Name:     "code",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Infra) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("InfraTarget")
	s.SetName(k8s.RFC1123Name(id))

	code := config.GetAttr("code").AsString()
	_ = unstructured.SetNestedField(s.Object, code, "spec", "script", "code")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Infra) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "InfraTarget", k8s.RFC1123Name(id))
}
