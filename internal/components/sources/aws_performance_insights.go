/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AWSPerformanceInsights struct{}

var (
	_ translation.Decodable    = (*AWSPerformanceInsights)(nil)
	_ translation.Translatable = (*AWSPerformanceInsights)(nil)
)

// Spec implements translation.Decodable.
func (*AWSPerformanceInsights) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"polling_interval": &hcldec.AttrSpec{
			Name:     "polling_interval",
			Type:     cty.String,
			Required: true,
		},
		"credentials": &hcldec.AttrSpec{
			Name:     "credentials",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
		"metric_query": &hcldec.AttrSpec{
			Name:     "metric_query",
			Type:     cty.String,
			Required: true,
		},
		"identifier": &hcldec.AttrSpec{
			Name:     "identifier",
			Type:     cty.String,
			Required: true,
		},
		"service_type": &hcldec.AttrSpec{
			Name:     "service_type",
			Type:     cty.String,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*AWSPerformanceInsights) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject(k8s.APISources, "AWSPerformanceInsightsSource", k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	s.SetNestedField(arn, "spec", "arn")

	pi := config.GetAttr("polling_interval").AsString()
	s.SetNestedField(pi, "spec", "pollingInterval")

	mq := config.GetAttr("metric_query").AsString()
	s.SetNestedField(mq, "spec", "metricQuery")

	i := config.GetAttr("identifier").AsString()
	s.SetNestedField(i, "spec", "identifier")

	st := config.GetAttr("service_type").AsString()
	s.SetNestedField(st, "spec", "serviceType")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
