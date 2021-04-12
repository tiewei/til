package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/k8s"
	"bridgedl/translation"
)

type Salesforce struct{}

var (
	_ translation.Decodable    = (*Salesforce)(nil)
	_ translation.Translatable = (*Salesforce)(nil)
)

// Spec implements translation.Decodable.
func (*Salesforce) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"channel": &hcldec.AttrSpec{
			Name:     "channel",
			Type:     cty.String,
			Required: true,
		},
		"replay_ID": &hcldec.AttrSpec{
			Name:     "replay_ID",
			Type:     cty.String,
			Required: false,
		},
		"client_ID": &hcldec.AttrSpec{
			Name:     "client_ID",
			Type:     cty.String,
			Required: true,
		},
		"server": &hcldec.AttrSpec{
			Name:     "server",
			Type:     cty.String,
			Required: false,
		},
		"user": &hcldec.AttrSpec{
			Name:     "user",
			Type:     cty.String,
			Required: true,
		},
		"cert_key": &hcldec.AttrSpec{
			Name:     "cert_key",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Salesforce) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.triggermesh.io/v1alpha1")
	s.SetKind("SalesforceSource")
	s.SetName(k8s.RFC1123Name(id))

	channel := config.GetAttr("channel").AsString()
	_ = unstructured.SetNestedField(s.Object, channel, "spec", "subscription", "channel")

	replayID := config.GetAttr("replay_ID")
	if !replayID.IsNull() {
		_ = unstructured.SetNestedField(s.Object, replayID.AsString(), "spec", "subscription", "replayID")
	}

	clientID := config.GetAttr("client_ID").AsString()
	_ = unstructured.SetNestedField(s.Object, clientID, "spec", "auth", "clientID")

	server := config.GetAttr("server")
	if !server.IsNull() {
		_ = unstructured.SetNestedField(s.Object, server.AsString(), "spec", "auth", "server")
	}

	user := config.GetAttr("user").AsString()
	_ = unstructured.SetNestedField(s.Object, user, "spec", "auth", "user")

	certKey := config.GetAttr("cert_key").AsString()
	_ = unstructured.SetNestedField(s.Object, certKey, "spec", "auth", "certKey", "value")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
