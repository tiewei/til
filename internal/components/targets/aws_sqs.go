package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

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

	t := k8s.NewObject("targets.triggermesh.io/v1alpha1", "AWSSQSTarget", id)

	arn := config.GetAttr("arn").AsString()
	t.SetNestedField(arn, "spec", "arn")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	t.SetNestedMap(accKeySecretRef, "spec", "awsApiKey", "secretKeyRef")
	t.SetNestedMap(secrKeySecretRef, "spec", "awsApiSecret", "secretKeyRef")

	return append(manifests, t.Unstructured())
}

// Address implements translation.Addressable.
func (*AWSSQS) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "AWSSQSTarget", k8s.RFC1123Name(id))
}
