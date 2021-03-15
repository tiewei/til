module bridgedl

go 1.16

require (
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/zclconf/go-cty v1.8.0
)

// TODO(antoineco): this dependency increases the executable's size from 5MiB
// to 15MiB. Considering we only import it for the unstructured.Unstructured
// type, it might be worth implementing our own marshalable type instead.
require k8s.io/apimachinery v0.20.4
