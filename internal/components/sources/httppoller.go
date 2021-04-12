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
			Required: false,
		},
		"interval": &hcldec.AttrSpec{
			Name:     "interval",
			Type:     cty.String,
			Required: false,
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

	eventSource := config.GetAttr("event_source")
	if !eventSource.IsNull() {
		_ = unstructured.SetNestedField(s.Object, eventSource.AsString(), "spec", "eventSource")
	}

	endpoint := config.GetAttr("endpoint").AsString()
	_ = unstructured.SetNestedField(s.Object, endpoint, "spec", "endpoint")

	method := config.GetAttr("method")
	if !method.IsNull() {
		_ = unstructured.SetNestedField(s.Object, method.AsString(), "spec", "method")
	}

	interval := config.GetAttr("interval")
	if !interval.IsNull() {
		_ = unstructured.SetNestedField(s.Object, interval.AsString(), "spec", "interval")
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
