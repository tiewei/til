package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
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

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("HTTPPollerSource")
	s.SetName(k8s.RFC1123Name(id))

	eventType := config.GetAttr("event_type").AsString()
	_ = unstructured.SetNestedField(s.Object, eventType, "spec", "eventType")

	if v := config.GetAttr("event_source"); !v.IsNull() {
		eventSource := v.AsString()
		_ = unstructured.SetNestedField(s.Object, eventSource, "spec", "eventSource")
	}

	endpoint := config.GetAttr("endpoint").AsString()
	_ = unstructured.SetNestedField(s.Object, endpoint, "spec", "endpoint")

	method := config.GetAttr("method").AsString()
	_ = unstructured.SetNestedField(s.Object, method, "spec", "method")

	interval := config.GetAttr("interval").AsString()
	_ = unstructured.SetNestedField(s.Object, interval, "spec", "interval")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
