package sources

import (
	"encoding/json"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
	"bridgedl/translation"
)

type Ping struct{}

var (
	_ translation.Decodable    = (*Ping)(nil)
	_ translation.Translatable = (*Ping)(nil)
)

// Spec implements translation.Decodable.
func (*Ping) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"schedule": &hcldec.AttrSpec{
			Name:     "schedule",
			Type:     cty.String,
			Required: false,
		},
		"data": &hcldec.AttrSpec{
			Name:     "data",
			Type:     cty.String,
			Required: true,
		},
		"content_type": &hcldec.AttrSpec{
			Name:     "content_type",
			Type:     cty.String,
			Required: false,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Ping) Manifests(id string, config, eventDst cty.Value) []interface{} {
	const defaultSchedule = "* * * * *" // every minute

	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.knative.dev/v1beta2")
	s.SetKind("PingSource")
	s.SetName(k8s.RFC1123Name(id))

	schedule := defaultSchedule
	if v := config.GetAttr("schedule"); !v.IsNull() {
		schedule = v.AsString()
	}
	_ = unstructured.SetNestedField(s.Object, schedule, "spec", "schedule")

	data := config.GetAttr("data").AsString()
	_ = unstructured.SetNestedField(s.Object, data, "spec", "data")

	if v := config.GetAttr("content_type"); !v.IsNull() {
		contentType := v.AsString()
		_ = unstructured.SetNestedField(s.Object, contentType, "spec", "contentType")
	} else if json.Valid([]byte(data)) {
		contentType := "application/json"
		_ = unstructured.SetNestedField(s.Object, contentType, "spec", "contentType")
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
