package k8s

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

const APIServing = "serving.knative.dev/v1"

// NewKnService returns a new Knative Service.
func NewKnService(name, image string, public bool) *unstructured.Unstructured {
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

	return s
}
