package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
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
			Required: true,
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

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("WebhookSource")
	s.SetName(k8s.RFC1123Name(id))

	eventType := config.GetAttr("event_type").AsString()
	_ = unstructured.SetNestedField(s.Object, eventType, "spec", "eventType")

	eventSource := config.GetAttr("event_source").AsString()
	_ = unstructured.SetNestedField(s.Object, eventSource, "spec", "eventSource")

	basicAuthUsername := config.GetAttr("basic_auth_username")
	if !basicAuthUsername.IsNull() {
		_ = unstructured.SetNestedField(s.Object, basicAuthUsername.AsString(), "spec", "basicAuthUsername")
	}

	basicAuthPassword := config.GetAttr("basic_auth_password")
	if !basicAuthPassword.IsNull() {
		_ = unstructured.SetNestedField(s.Object, basicAuthPassword.AsString(), "spec", "basicAuthPassword", "value")
	}

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
