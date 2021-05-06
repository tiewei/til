package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("sources.knative.dev/v1alpha1")
	s.SetKind("GitHubSource")
	s.SetName(k8s.RFC1123Name(id))

	var eventTypes []interface{}
	eTIter := config.GetAttr("event_types").ElementIterator()
	for eTIter.Next() {
		_, srv := eTIter.Element()
		eventTypes = append(eventTypes, srv.AsString())
	}
	_ = unstructured.SetNestedSlice(s.Object, eventTypes, "spec", "eventTypes")

	ownerAndRepository := config.GetAttr("owner_and_repository").AsString()
	_ = unstructured.SetNestedField(s.Object, ownerAndRepository, "spec", "ownerAndRepository")

	tokens := config.GetAttr("tokens").GetAttr("name").AsString()
	accTokenSecretRef, webhookSecretRef := secrets.SecretKeyRefsGitHub(tokens)
	_ = unstructured.SetNestedMap(s.Object, accTokenSecretRef, "spec", "accessToken", "secretKeyRef")
	_ = unstructured.SetNestedMap(s.Object, webhookSecretRef, "spec", "secretToken", "secretKeyRef")

	sinkRef := eventDst.GetAttr("ref")
	sink := map[string]interface{}{
		"apiVersion": sinkRef.GetAttr("apiVersion").AsString(),
		"kind":       sinkRef.GetAttr("kind").AsString(),
		"name":       sinkRef.GetAttr("name").AsString(),
	}
	_ = unstructured.SetNestedMap(s.Object, sink, "spec", "sink", "ref")

	return append(manifests, s)
}
