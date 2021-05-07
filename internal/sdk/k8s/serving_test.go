package k8s_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "bridgedl/internal/sdk/k8s"
)

func TestNewKnService(t *testing.T) {
	const (
		name  = "test"
		image = "example.com/myapp:1.0"
	)

	ksvc := NewKnService(name, image, false,
		EnvVars(map[string]cty.Value{
			"TEST_ENV_FOO": cty.StringVal("foo"),
			"TEST_ENV_BAR": cty.StringVal("bar"),
		}),
	)

	expectKsvc := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": APIServing,
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name": name,
				"labels": map[string]interface{}{
					"networking.knative.dev/visibility": "cluster-local",
				},
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"image": image,
								"env": []interface{}{
									map[string]interface{}{
										"name":  "TEST_ENV_BAR",
										"value": "bar",
									},
									map[string]interface{}{
										"name":  "TEST_ENV_FOO",
										"value": "foo",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if d := cmp.Diff(expectKsvc, ksvc); d != "" {
		t.Errorf("Unexpected diff: (-:expect, +:got) %s", d)
	}
}
