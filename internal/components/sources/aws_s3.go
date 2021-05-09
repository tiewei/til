package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk"
	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AWSS3 struct{}

var (
	_ translation.Decodable    = (*AWSS3)(nil)
	_ translation.Translatable = (*AWSS3)(nil)
)

// Spec implements translation.Decodable.
func (*AWSS3) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"event_types": &hcldec.AttrSpec{
			Name:     "event_types",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"queue_arn": &hcldec.AttrSpec{
			Name:     "queue_arn",
			Type:     cty.String,
			Required: false,
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

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "AWSS3Source", k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	s.SetNestedField(arn, "spec", "arn")

	eventTypes := sdk.DecodeStringSlice(config.GetAttr("event_types"))
	s.SetNestedSlice(eventTypes, "spec", "eventTypes")

	if !config.GetAttr("queue_arn").IsNull() {
		queueARN := config.GetAttr("queue_arn").AsString()
		s.SetNestedField(queueARN, "spec", "queueARN")
	}

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
