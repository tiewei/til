package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Slack struct{}

var (
	_ translation.Decodable    = (*Slack)(nil)
	_ translation.Translatable = (*Slack)(nil)
)

// Spec implements translation.Decodable.
func (*Slack) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"signing_secret": &hcldec.AttrSpec{
			Name:     "signing_secret",
			Type:     k8s.ObjectReferenceCty,
			Required: false,
		},
		"app_id": &hcldec.AttrSpec{
			Name:     "app_id",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Slack) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "SlackSource", id)

	if v := config.GetAttr("signing_secret"); !v.IsNull() {
		signingSecretSecretName := v.GetAttr("name").AsString()
		signingSecret := secrets.SecretKeyRefsSlackApp(signingSecretSecretName)
		s.SetNestedMap(signingSecret, "spec", "signingSecret", "valueFromSecret")
	}

	appID := config.GetAttr("app_id")
	if !appID.IsNull() {
		s.SetNestedField(appID.AsString(), "spec", "appID")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
