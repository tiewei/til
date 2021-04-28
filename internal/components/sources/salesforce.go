package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/internal/sdk/validation"
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
		"replay_id": &hcldec.ValidateSpec{
			Wrapped: &hcldec.AttrSpec{
				Name:     "replay_id",
				Type:     cty.Number,
				Required: false,
			},
			Func: validation.IsInt(),
		},
		"client_id": &hcldec.AttrSpec{
			Name:     "client_id",
			Type:     cty.String,
			Required: true,
		},
		"server": &hcldec.AttrSpec{
			Name:     "server",
			Type:     cty.String,
			Required: true,
		},
		"user": &hcldec.AttrSpec{
			Name:     "user",
			Type:     cty.String,
			Required: true,
		},
		"secret_key": &hcldec.AttrSpec{
			Name:     "secret_key",
			Type:     k8s.ObjectReferenceCty,
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

	if v := config.GetAttr("replay_id"); !v.IsNull() {
		replayID, _ := v.AsBigFloat().Int64()
		_ = unstructured.SetNestedField(s.Object, replayID, "spec", "subscription", "replayID")
	}

	clientID := config.GetAttr("client_id").AsString()
	_ = unstructured.SetNestedField(s.Object, clientID, "spec", "auth", "clientID")

	server := config.GetAttr("server").AsString()
	_ = unstructured.SetNestedField(s.Object, server, "spec", "auth", "server")

	user := config.GetAttr("user").AsString()
	_ = unstructured.SetNestedField(s.Object, user, "spec", "auth", "user")

	secrKeySecretName := config.GetAttr("secret_key").GetAttr("name").AsString()
	secrKeySecretRef := secrets.SecretKeyRefsSalesforceOAuthJWT(secrKeySecretName)
	_ = unstructured.SetNestedMap(s.Object, secrKeySecretRef, "spec", "auth", "certKey", "valueFromSecret")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
