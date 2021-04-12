package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
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
			Type:     cty.String,
			Required: false,
		},
		"app_ID": &hcldec.AttrSpec{
			Name:     "app_ID",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Slack) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("SlackSource")
	s.SetName(k8s.RFC1123Name(id))

	signingSecret := config.GetAttr("signing_secret")
	if !signingSecret.IsNull() {
		_ = unstructured.SetNestedField(s.Object, signingSecret.AsString(), "spec", "signingSecret", "value")
	}

	appID := config.GetAttr("app_ID")
	if !appID.IsNull() {
		_ = unstructured.SetNestedField(s.Object, appID.AsString(), "spec", "appID")
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
