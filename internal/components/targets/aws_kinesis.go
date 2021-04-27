package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
	"bridgedl/translation"
)

type AWSKinesis struct{}

var (
	_ translation.Decodable    = (*AWSKinesis)(nil)
	_ translation.Translatable = (*AWSKinesis)(nil)
	_ translation.Addressable  = (*AWSKinesis)(nil)
)

// Spec implements translation.Decodable.
func (*AWSKinesis) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"partition": &hcldec.AttrSpec{
			Name:     "partition",
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
func (*AWSKinesis) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("AWSKinesisTarget")
	s.SetName(k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	_ = unstructured.SetNestedField(s.Object, arn, "spec", "arn")

	partition := config.GetAttr("partition").AsString()
	_ = unstructured.SetNestedField(s.Object, partition, "spec", "partition")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	_ = unstructured.SetNestedMap(s.Object, accKeySecretRef, "spec", "awsApiKey", "secretKeyRef")
	_ = unstructured.SetNestedMap(s.Object, secrKeySecretRef, "spec", "awsApiSecret", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*AWSKinesis) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "AWSKinesisTarget", k8s.RFC1123Name(id))
}
