package k8s

import "github.com/zclconf/go-cty/cty"

// DestinationCty is a non-primitive cty.Type that represents a duck Destination.
var DestinationCty = cty.ObjectWithOptionalAttrs(
	map[string]cty.Type{
		"ref": cty.Object(map[string]cty.Type{
			"apiVersion": cty.String,
			"kind":       cty.String,
			"name":       cty.String,
		}),
		"uri": cty.String,
	},
	[]string{"ref", "uri"},
)
