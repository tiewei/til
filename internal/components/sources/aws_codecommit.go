package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
	"bridgedl/translation"
)

type AWSCodeCommit struct{}

var (
	_ translation.Decodable    = (*AWSCodeCommit)(nil)
	_ translation.Translatable = (*AWSCodeCommit)(nil)
)

// Spec implements translation.Decodable.
func (*AWSCodeCommit) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"branch": &hcldec.AttrSpec{
			Name:     "branch",
			Type:     cty.String,
			Required: true,
		},
		"event_types": &hcldec.AttrSpec{
			Name:     "event_types",
			Type:     cty.List(cty.String),
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
func (*AWSCodeCommit) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("AWSCodeCommitSource")
	s.SetName(k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	_ = unstructured.SetNestedField(s.Object, arn, "spec", "arn")

	branch := config.GetAttr("branch").AsString()
	_ = unstructured.SetNestedField(s.Object, branch, "spec", "branch")

	eventTypes := config.GetAttr("event_types").AsValueSlice()
	var stringSlice []string
	for _, v := range eventTypes {
		stringSlice = append(stringSlice, v.AsString())
	}
	_ = unstructured.SetNestedStringSlice(s.Object, stringSlice, "spec", "eventTypes")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	_ = unstructured.SetNestedMap(s.Object, accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	_ = unstructured.SetNestedMap(s.Object, secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
