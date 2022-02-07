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

	"til/config/globals"
	"til/internal/sdk"
	"til/internal/sdk/k8s"
	"til/internal/sdk/secrets"
	"til/translation"
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
		"metric_queries": &hcldec.AttrSpec{
			Name:     "metric_queries",
			Type:     cty.List(cty.String),
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*AWSPerformanceInsights) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "AWSPerformanceInsightsSource", name)

	arn := config.GetAttr("arn").AsString()
	s.SetNestedField(arn, "spec", "arn")

	pollingInterval := config.GetAttr("polling_interval").AsString()
	s.SetNestedField(pollingInterval, "spec", "pollingInterval")

	metricQueries := sdk.DecodeStringSlice(config.GetAttr("metric_queries"))
	s.SetNestedSlice(metricQueries, "spec", "metricQueries")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "auth", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "auth", "credentials", "secretAccessKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
