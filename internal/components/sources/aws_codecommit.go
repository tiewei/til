package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
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

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "AWSCodeCommitSource", id)

	arn := config.GetAttr("arn").AsString()
	s.SetNestedField(arn, "spec", "arn")

	branch := config.GetAttr("branch").AsString()
	s.SetNestedField(branch, "spec", "branch")

	eventTypesVals := config.GetAttr("event_types").AsValueSlice()
	eventTypes := make([]interface{}, 0, len(eventTypesVals))
	for _, v := range eventTypesVals {
		eventTypes = append(eventTypes, v.AsString())
	}
	s.SetNestedSlice(eventTypes, "spec", "eventTypes")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
