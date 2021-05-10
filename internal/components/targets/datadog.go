package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Datadog struct{}

var (
	_ translation.Decodable    = (*Datadog)(nil)
	_ translation.Translatable = (*Datadog)(nil)
	_ translation.Addressable  = (*Datadog)(nil)
)

// Spec implements translation.Decodable.
func (*Datadog) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"metric_prefix": &hcldec.AttrSpec{
			Name:     "metric_prefix",
			Type:     cty.String,
			Required: false,
		},
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Datadog) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "DatadogTarget", name)

	if v := config.GetAttr("metric_prefix"); !v.IsNull() {
		mp := config.GetAttr("metric_prefix").AsString()
		t.SetNestedField(mp, "spec", "metricPrefix")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	apiKeySecretRef := secrets.SecretKeyRefsDatadog(authSecretName)
	t.SetNestedMap(apiKeySecretRef, "spec", "apiKey", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "DatadogTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Datadog) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "DatadogTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
