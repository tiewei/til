package k8s

import "strings"

// RFC1123Name sanitizes the given input string, ensuring it is a valid DNS
// subdomain name (as defined in RFC 1123) which can be used as a Kubernetes
// object name.
//
// The implementation is purposedly naive because it is assumed that this
// function will only ever sanitize HCL identifiers, which already contain a
// limited set of characters (see hclsyntax.ValidIdentifier).
func RFC1123Name(id string) string {
	return strings.ToLower(strings.ReplaceAll(id, "_", "-"))
}
