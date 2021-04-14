package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"bridgedl/internal/sdk/secrets"
	"bridgedl/k8s"
	"bridgedl/translation"
)

type Confluent struct{}

var (
	_ translation.Decodable    = (*Confluent)(nil)
	_ translation.Translatable = (*Confluent)(nil)
	_ translation.Addressable  = (*Confluent)(nil)
)

// Spec implements translation.Decodable.
func (*Confluent) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"topic": &hcldec.AttrSpec{
			Name:     "topic",
			Type:     cty.String,
			Required: true,
		},
		"bootstrap_servers": &hcldec.AttrSpec{
			Name:     "bootstrap_servers",
			Type:     cty.List(cty.String),
			Required: true,
		},
		"security_protocol": &hcldec.AttrSpec{
			Name:     "security_protocol",
			Type:     cty.String,
			Required: false,
		},
		"sasl_mechanism": &hcldec.AttrSpec{
			Name:     "sasl_mechanism",
			Type:     cty.String,
			Required: false,
		},
		"username": &hcldec.AttrSpec{
			Name:     "username",
			Type:     cty.String,
			Required: true,
		},
		"sasl_auth": &hcldec.AttrSpec{
			Name:     "sasl_auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Confluent) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := &unstructured.Unstructured{}
	s.SetAPIVersion("targets.triggermesh.io/v1alpha1")
	s.SetKind("ConfluentTarget")
	s.SetName(k8s.RFC1123Name(id))

	arn := config.GetAttr("topic").AsString()
	_ = unstructured.SetNestedField(s.Object, arn, "spec", "topic")

	var bootstrapServers []interface{}

	bSrvsIter := config.GetAttr("bootstrap_servers").ElementIterator()
	for bSrvsIter.Next() {
		_, srv := bSrvsIter.Element()
		bootstrapServers = append(bootstrapServers, srv.AsString())
	}
	_ = unstructured.SetNestedSlice(s.Object, bootstrapServers, "spec", "bootstrapServers")

	if v := config.GetAttr("security_protocol"); !v.IsNull() {
		securityProtocol := v.AsString()
		_ = unstructured.SetNestedField(s.Object, securityProtocol, "spec", "securityProtocol")
	}

	if v := config.GetAttr("sasl_mechanism"); !v.IsNull() {
		saslMechanism := v.AsString()
		_ = unstructured.SetNestedField(s.Object, saslMechanism, "spec", "saslMechanism")
	}

	username := config.GetAttr("username").AsString()
	_ = unstructured.SetNestedField(s.Object, username, "spec", "username")

	saslAuthSecretName := config.GetAttr("sasl_auth").GetAttr("name").AsString()
	passwd := secrets.SecretKeyRefsConfluent(saslAuthSecretName)
	_ = unstructured.SetNestedMap(s.Object, passwd, "spec", "password", "secretKeyRef")

	return append(manifests, s)
}

// Address implements translation.Addressable.
func (*Confluent) Address(id string, _, _ cty.Value) cty.Value {
	return k8s.NewDestination("targets.triggermesh.io/v1alpha1", "ConfluentTarget", k8s.RFC1123Name(id))
}
