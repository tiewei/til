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
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"til/config/globals"
	"til/internal/sdk/k8s"
	"til/internal/sdk/secrets"
	"til/internal/sdk/validation"
	"til/translation"
)

type AWSCloudWatch struct{}

var (
	_ translation.Decodable    = (*AWSCloudWatch)(nil)
	_ translation.Translatable = (*AWSCloudWatch)(nil)
)

// Spec implements translation.Decodable.
func (*AWSCloudWatch) Spec() hcldec.Spec {
	metricQuerySpec := &hcldec.ValidateSpec{
		Wrapped: &hcldec.ObjectSpec{
			"name": &hcldec.BlockLabelSpec{
				Index: 0,
				Name:  "name",
			},
			"expression": &hcldec.AttrSpec{
				Name:     "expression",
				Type:     cty.String,
				Required: false,
			},
			"metric": &hcldec.BlockSpec{
				TypeName: "metric",
				Nested: &hcldec.ObjectSpec{
					"period": &hcldec.ValidateSpec{
						Wrapped: &hcldec.AttrSpec{
							Name:     "period",
							Type:     cty.Number,
							Required: false,
						},
						Func: validation.IsInt,
					},
					"stat": &hcldec.AttrSpec{
						Name:     "stat",
						Type:     cty.String,
						Required: true,
					},
					"unit": &hcldec.AttrSpec{
						Name:     "unit",
						Type:     cty.String,
						Required: true,
					},
					"name": &hcldec.AttrSpec{
						Name:     "name",
						Type:     cty.String,
						Required: true,
					},
					"namespace": &hcldec.AttrSpec{
						Name:     "namespace",
						Type:     cty.String,
						Required: true,
					},
					"dimension": &hcldec.BlockSetSpec{
						TypeName: "dimension",
						Nested: &hcldec.ObjectSpec{
							"name": &hcldec.BlockLabelSpec{
								Index: 0,
								Name:  "name",
							},
							"value": &hcldec.AttrSpec{
								Name:     "value",
								Type:     cty.String,
								Required: true,
							},
						},
						MinItems: 1,
					},
				},
				Required: false,
			},
		},
		Func: validateCloudWatchAttrMetricQuery,
	}

	return &hcldec.ObjectSpec{
		"region": &hcldec.AttrSpec{
			Name:     "region",
			Type:     cty.String,
			Required: true,
		},
		"polling_interval": &hcldec.AttrSpec{
			Name:     "polling_interval",
			Type:     cty.String,
			Required: false,
		},
		"metric_query": &hcldec.BlockSetSpec{
			TypeName: "metric_query",
			Nested:   metricQuerySpec,
			MinItems: 1,
		},
		"credentials": &hcldec.AttrSpec{
			Name:     "credentials",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}

	/*
		Example of value decoded from the spec above, for a config with
		two "metric_query" blocks, one with an "expression" attribute,
		one with a "metric" block:

		v: (map[string]interface {}) (len=4) {
		 region": (string) "us-east-2",
		 polling_interval": (interface {}) <nil>,
		 credentials": (map[string]interface {}) (len=1) {
		  name": (string) "my-aws-access-keys"
		 },
		 "metric_query": (set.Set) {
		  vals: (map[int][]interface {}) (len=2) {
		   (int) 2800467745: ([]interface {}) (len=1 cap=1) {
		    (map[string]interface {}) (len=3) {
		     "name": (string) "query1",
		     "expression": (string) "SEARCH(' {AWS/EC2} MetricName=\"CPUUtilization\" ', 'Average', 300)",
		     "metric": (interface {}) <nil>
		    }
		   },
		   (int) 719348152: ([]interface {}) (len=1 cap=1) {
		    (map[string]interface {}) (len=3) {
		     name": (string) "query2",
		     "expression": (interface {}) <nil>,
		     "metric": (map[string]interface {}) (len=6) {
		      "period": (*big.Float)(0xc0001dacf0)(60),
		      "stat": (string) "p90",
		      "unit": (string) "Milliseconds",
		      "name": (string) "Duration",
		      "namespace": (string) "AWS/Lambda",
		      "dimension": (set.Set) {
		       vals: (map[int][]interface {}) (len=1) {
		        (int) 828842044: ([]interface {}) (len=1 cap=1) {
		         (map[string]interface {}) (len=2) {
		          "value": (string) "lambdadumper",
		          "name": (string) "FunctionName"
		         }
		        }
		       }
		      }
		     }
		    }
		   }
		  }
		 }
		}
	*/
}

// Manifests implements translation.Translatable.
func (*AWSCloudWatch) Manifests(id string, config, eventDst cty.Value, glb globals.Accessor) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	manifests, eventDst = k8s.MaybeAppendChannel(name, manifests, eventDst, glb)

	s := k8s.NewObject(k8s.APISources, "AWSCloudWatchSource", name)

	region := config.GetAttr("region").AsString()
	s.SetNestedField(region, "spec", "region")

	if v := config.GetAttr("polling_interval"); !v.IsNull() {
		pollingInterval := v.AsString()
		s.SetNestedField(pollingInterval, "spec", "pollingInterval")
	}

	metricQueriesVal := config.GetAttr("metric_query")
	metricQueries := make([]interface{}, 0, metricQueriesVal.LengthInt())
	for mqIter := metricQueriesVal.ElementIterator(); mqIter.Next(); {
		_, metricQueryVal := mqIter.Element()

		metricQuery := make(map[string]interface{}, 2)

		metricQuery["name"] = metricQueryVal.GetAttr("name").AsString()

		if v := metricQueryVal.GetAttr("expression"); !v.IsNull() {
			expression := v.AsString()
			metricQuery["expression"] = expression
		}

		if metricVal := metricQueryVal.GetAttr("metric"); !metricVal.IsNull() {
			metric := make(map[string]interface{}, 4)

			period, _ := metricVal.GetAttr("period").AsBigFloat().Int64()
			metric["period"] = period
			metric["stat"] = metricVal.GetAttr("stat").AsString()
			metric["unit"] = metricVal.GetAttr("unit").AsString()

			metricmetric := make(map[string]interface{}, 3)
			metricmetric["metricName"] = metricVal.GetAttr("name").AsString()
			metricmetric["namespace"] = metricVal.GetAttr("namespace").AsString()

			dimensionsVal := metricVal.GetAttr("dimension")
			dimensions := make([]interface{}, 0, dimensionsVal.LengthInt())
			for dimIter := dimensionsVal.ElementIterator(); dimIter.Next(); {
				_, dimensionVal := dimIter.Element()

				dimension := map[string]interface{}{
					"name":  dimensionVal.GetAttr("name").AsString(),
					"value": dimensionVal.GetAttr("value").AsString(),
				}

				dimensions = append(dimensions, dimension)
			}
			metricmetric["dimensions"] = dimensions

			metric["metric"] = metricmetric

			metricQuery["metric"] = metric
		}

		metricQueries = append(metricQueries, metricQuery)
	}
	s.SetNestedSlice(metricQueries, "spec", "metricQueries")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}

func validateCloudWatchAttrMetricQuery(val cty.Value) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if val.GetAttr("expression").IsNull() == val.GetAttr("metric").IsNull() {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Mutually exclusive attributes",
			Detail:   `The "expression" and "metric" attributes are mutually exclusive.`,
		})
	}

	return diags
}
