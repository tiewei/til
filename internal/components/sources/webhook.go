package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type Webhook struct{}

var (
	_ translation.Decodable    = (*Webhook)(nil)
	_ translation.Translatable = (*Webhook)(nil)
)

// Spec implements translation.Decodable.
func (*Webhook) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"event_type": &hcldec.AttrSpec{
			Name:     "event_type",
			Type:     cty.String,
			Required: true,
		},
		"event_source": &hcldec.AttrSpec{
			Name:     "event_source",
			Type:     cty.String,
			Required: false,
		},
		"basic_auth_username": &hcldec.AttrSpec{
			Name:     "basic_auth_username",
			Type:     cty.String,
			Required: false,
		},
		"basic_auth_password": &hcldec.AttrSpec{
			Name:     "basic_auth_password",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Webhook) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "WebhookSource", id)

	eventType := config.GetAttr("event_type").AsString()
	s.SetNestedField(eventType, "spec", "eventType")

	if v := config.GetAttr("event_source"); !v.IsNull() {
		eventSource := v.AsString()
		s.SetNestedField(eventSource, "spec", "eventSource")
	}

	basicAuthUsername := config.GetAttr("basic_auth_username")
	if !basicAuthUsername.IsNull() {
		s.SetNestedField(basicAuthUsername.AsString(), "spec", "basicAuthUsername")
	}

	basicAuthPassword := config.GetAttr("basic_auth_password")
	if !basicAuthPassword.IsNull() {
		s.SetNestedField(basicAuthPassword.AsString(), "spec", "basicAuthPassword", "value")
	}

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
