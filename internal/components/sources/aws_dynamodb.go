package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
	"bridgedl/translation"
)

type AWSDynamoDB struct{}

var (
	_ translation.Decodable    = (*AWSDynamoDB)(nil)
	_ translation.Translatable = (*AWSDynamoDB)(nil)
)

// Spec implements translation.Decodable.
func (*AWSDynamoDB) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"access_key": &hcldec.AttrSpec{
			Name:     "access_key",
			Type:     cty.String,
			Required: true,
		},
		"secret_key": &hcldec.AttrSpec{
			Name:     "secret_key",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*AWSDynamoDB) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("AWSDynamoDB")
	s.SetName(k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	_ = unstructured.SetNestedField(s.Object, arn, "spec", "arn")

	accessKey := config.GetAttr("access_key").AsString()
	secretKey := config.GetAttr("secret_key").AsString()
	_ = unstructured.SetNestedField(s.Object, accessKey, "spec", "credentials", "accessKeyID", "value")
	_ = unstructured.SetNestedField(s.Object, secretKey, "spec", "credentials", "secretAccessKey", "value")

	sinkRef := eventDst.GetAttr("ref")

	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}

	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
