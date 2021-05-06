package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type GitHub struct{}

var (
	_ translation.Decodable    = (*GitHub)(nil)
	_ translation.Translatable = (*GitHub)(nil)
)

// Spec implements translation.Decodable.
func (*GitHub) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"event_types": &hcldec.AttrSpec{
			Name:     "event_types",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"owner_and_repository": &hcldec.AttrSpec{
			Name:     "owner_and_repository",
			Type:     cty.String,
			Required: true,
		},
		"tokens": &hcldec.AttrSpec{
			Name:     "tokens",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*GitHub) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject("sources.knative.dev/v1alpha1", "GitHubSource", id)

	eventTypesVals := config.GetAttr("event_types").AsValueSlice()
	eventTypes := make([]interface{}, 0, len(eventTypesVals))
	for _, v := range eventTypesVals {
		eventTypes = append(eventTypes, v.AsString())
	}
	s.SetNestedSlice(eventTypes, "spec", "eventTypes")

	ownerAndRepository := config.GetAttr("owner_and_repository").AsString()
	s.SetNestedField(ownerAndRepository, "spec", "ownerAndRepository")

	tokens := config.GetAttr("tokens").GetAttr("name").AsString()
	accTokenSecretRef, webhookSecretRef := secrets.SecretKeyRefsGitHub(tokens)
	s.SetNestedMap(accTokenSecretRef, "spec", "accessToken", "secretKeyRef")
	s.SetNestedMap(webhookSecretRef, "spec", "secretToken", "secretKeyRef")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
