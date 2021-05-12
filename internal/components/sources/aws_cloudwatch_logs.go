package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AWSCloudWatchLogs struct{}

var (
	_ translation.Decodable    = (*AWSCloudWatchLogs)(nil)
	_ translation.Translatable = (*AWSCloudWatchLogs)(nil)
)

// Spec implements translation.Decodable.
func (*AWSCloudWatchLogs) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"polling_interval": &hcldec.AttrSpec{
			Name:     "polling_interval",
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
func (*AWSCloudWatchLogs) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "AWSCloudWatchLogsSource", k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	s.SetNestedField(arn, "spec", "arn")

	if v := config.GetAttr("polling_interval"); !v.IsNull() {
		pi := v.AsString()
		s.SetNestedField(pi, "spec", "pollingInterval")
	}

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
