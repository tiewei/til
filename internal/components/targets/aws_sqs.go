package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AWSSQS struct{}

var (
	_ translation.Decodable    = (*AWSSQS)(nil)
	_ translation.Translatable = (*AWSSQS)(nil)
	_ translation.Addressable  = (*AWSSQS)(nil)
)

// Spec implements translation.Decodable.
func (*AWSSQS) Spec() hcldec.Spec {
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
func (*AWSSQS) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("AWSSQSTarget")
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
func (*AWSSQS) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "AWSSQSTarget", k8s.RFC1123Name(id))
}
