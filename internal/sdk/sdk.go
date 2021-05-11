package sdk

import "github.com/zclconf/go-cty/cty"

// DecodeStringSlice decodes a cty.Value into a slice of strings, in a format
// that is compatible with the "unstructured" package from k8s.io/apimachinery.
// Panics if the given value is not a non-null collection type containing
// exclusively string values.
func DecodeStringSlice(val cty.Value) []interface{} {
	out := make([]interface{}, 0, val.LengthInt())
	for iter := val.ElementIterator(); iter.Next(); {
		_, v := iter.Element()
		out = append(out, v.AsString())
	}

	return out
}

// DecodeStringMap decodes a cty.Value into a map of string elements, in a
// format that is compatible with the "unstructured" package from k8s.io/apimachinery.
// Panics if the given value is not a non-null collection type containing
// exclusively string values.
func DecodeStringMap(val cty.Value) map[string]interface{} {
	out := make(map[string]interface{}, val.LengthInt())
	for iter := val.ElementIterator(); iter.Next(); {
		if k, v := iter.Element(); !v.IsNull() {
			out[k.AsString()] = v.AsString()
		}
	}

	return out
}
