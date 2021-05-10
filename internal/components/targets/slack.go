package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Slack struct{}

var (
	_ translation.Decodable    = (*Slack)(nil)
	_ translation.Translatable = (*Slack)(nil)
	_ translation.Addressable  = (*Slack)(nil)
)

// Spec implements translation.Decodable.
func (*Slack) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Slack) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "SlackTarget", name)

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	apiTokenSecretRef := secrets.SecretKeyRefsSlack(authSecretName)
	t.SetNestedMap(apiTokenSecretRef, "spec", "token", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "SlackTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Slack) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "SlackTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
