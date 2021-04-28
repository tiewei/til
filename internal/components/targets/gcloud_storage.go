package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
	"bridgedl/translation"
)

type GCloudStorage struct{}

var (
	_ translation.Decodable    = (*GCloudStorage)(nil)
	_ translation.Translatable = (*GCloudStorage)(nil)
	_ translation.Addressable  = (*GCloudStorage)(nil)
)

// Spec implements translation.Decodable.
func (*GCloudStorage) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"bucket_name": &hcldec.AttrSpec{
			Name:     "bucket_name",
			Type:     cty.String,
			Required: true,
		},
		"service_account": &hcldec.AttrSpec{
			Name:     "service_account",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*GCloudStorage) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("GoogleCloudStorageTarget")
	s.SetName(k8s.RFC1123Name(id))

	bucketName := config.GetAttr("bucket_name").AsString()
	_ = unstructured.SetNestedField(s.Object, bucketName, "spec", "bucketName")

	svcAccountSecretName := config.GetAttr("service_account").GetAttr("name").AsString()
	keySecretRef := secrets.SecretKeyRefsGCloudServiceAccount(svcAccountSecretName)
	_ = unstructured.SetNestedMap(s.Object, keySecretRef, "spec", "credentialsJson", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*GCloudStorage) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "GoogleCloudStorageTarget", k8s.RFC1123Name(id))
}
