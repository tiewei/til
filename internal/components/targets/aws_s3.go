package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
	"bridgedl/translation"
)

type AWSS3 struct{}

var (
	_ translation.Decodable    = (*AWSS3)(nil)
	_ translation.Translatable = (*AWSS3)(nil)
	_ translation.Addressable  = (*AWSS3)(nil)
)

// Spec implements translation.Decodable.
func (*AWSS3) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"credentials": &hcldec.AttrSpec{
			Name:     "credentials",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*AWSS3) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("AWSS3Target")
	s.SetName(k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	_ = unstructured.SetNestedField(s.Object, arn, "spec", "arn")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	_ = unstructured.SetNestedMap(s.Object, accKeySecretRef, "spec", "awsApiKey", "secretKeyRef")
	_ = unstructured.SetNestedMap(s.Object, secrKeySecretRef, "spec", "awsApiSecret", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*AWSS3) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "AWSS3Target", k8s.RFC1123Name(id))
}
