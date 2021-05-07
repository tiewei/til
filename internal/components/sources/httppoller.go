package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/translation"
)

type HTTPPoller struct{}

var (
	_ translation.Decodable    = (*HTTPPoller)(nil)
	_ translation.Translatable = (*HTTPPoller)(nil)
)

// Spec implements translation.Decodable.
func (*HTTPPoller) Spec() hcldec.Spec {
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
		"endpoint": &hcldec.AttrSpec{
			Name:     "endpoint",
			Type:     cty.String,
			Required: true,
		},
		"method": &hcldec.AttrSpec{
			Name:     "method",
			Type:     cty.String,
			Required: true,
		},
		"interval": &hcldec.AttrSpec{
			Name:     "interval",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*HTTPPoller) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.triggermesh.io/v1alpha1", "HTTPPollerSource", id)

	eventType := config.GetAttr("event_type").AsString()
	s.SetNestedField(eventType, "spec", "eventType")

	if v := config.GetAttr("event_source"); !v.IsNull() {
		eventSource := v.AsString()
		s.SetNestedField(eventSource, "spec", "eventSource")
	}

	endpoint := config.GetAttr("endpoint").AsString()
	s.SetNestedField(endpoint, "spec", "endpoint")

	method := config.GetAttr("method").AsString()
	s.SetNestedField(method, "spec", "method")

	interval := config.GetAttr("interval").AsString()
	s.SetNestedField(interval, "spec", "interval")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
