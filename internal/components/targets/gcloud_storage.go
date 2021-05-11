package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
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

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "GoogleCloudStorageTarget", name)

	bucketName := config.GetAttr("bucket_name").AsString()
	t.SetNestedField(bucketName, "spec", "bucketName")

	svcAccountSecretName := config.GetAttr("service_account").GetAttr("name").AsString()
	keySecretRef := secrets.SecretKeyRefsGCloudServiceAccount(svcAccountSecretName)
	t.SetNestedMap(keySecretRef, "spec", "credentialsJson", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "GoogleCloudStorageTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*GCloudStorage) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "GoogleCloudStorageTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
