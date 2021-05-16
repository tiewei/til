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

package k8s

import (
	"sort"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const APIServing = "serving.knative.dev/v1"

// NewKnService returns a new Knative Service.
func NewKnService(name, image string, public bool, opts ...KnServiceOption) *unstructured.Unstructured {
	validateDNS1123Subdomain(name)

	s := &unstructured.Unstructured{}

	s.SetAPIVersion(APIServing)
	s.SetKind("Service")
	s.SetName(name)

	if !public {
		s.SetLabels(map[string]string{
			"networking.knative.dev/visibility": "cluster-local",
		})
	}

	container := make(map[string]interface{}, 1)
	_ = unstructured.SetNestedField(container, image, "image")
	_ = unstructured.SetNestedSlice(s.Object, []interface{}{container}, "spec", "template", "spec", "containers")

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// KnServiceOption is a functional option of a Knative Service.
type KnServiceOption func(*unstructured.Unstructured)

// EnvVars sets environment variables on the container of a Knative Service.
func EnvVars(evs map[string]cty.Value) KnServiceOption {
	return func(o *unstructured.Unstructured) {
		envVars := make([]interface{}, 0, len(evs))

		for k, v := range evs {
			envVar := map[string]interface{}{
				"name": k,
			}

			switch {
			case v.Type() == cty.String:
				envVar["value"] = v.AsString()
			case IsSecretKeySelector(v):
				envVar["valueFrom"] = map[string]interface{}{
					"secretKeyRef": map[string]interface{}{
						"name": v.GetAttr("name").AsString(),
						"key":  v.GetAttr("key").AsString(),
					},
				}
			default:
				continue
			}

			envVars = append(envVars, envVar)
		}

		if len(envVars) == 0 {
			return
		}

		sort.Slice(envVars, func(i, j int) bool {
			return envVars[i].(map[string]interface{})["name"].(string) <
				envVars[j].(map[string]interface{})["name"].(string)
		})

		containers, _, _ := unstructured.NestedSlice(o.Object, "spec", "template", "spec", "containers")
		_ = unstructured.SetNestedSlice(containers[0].(map[string]interface{}), envVars, "env")

		_ = unstructured.SetNestedSlice(o.Object, containers, "spec", "template", "spec", "containers")
	}
}
