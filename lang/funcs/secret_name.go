package funcs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"bridgedl/lang/k8s"
)

// SecretNameFunc returns the implementation of the "secret_name" function.
//
// secret_name() converts the name of a Kubernetes Secret object into an object
// reference to this Secret.
//
// Parameters:
//  * name: name of the Secret object. Must be a valid Kubernetes object name (RFC 1123 subdomain).
//
func SecretNameFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "name",
				Type: cty.String,
			},
		},

		Type: function.StaticReturnType(k8s.ObjectReferenceCty),
		Impl: secretNameFuncImpl,
	})
}

func secretNameFuncImpl(args []cty.Value, _ cty.Type) (cty.Value, error) {
	name := args[0].AsString()

	errs := validation.IsDNS1123Subdomain(name)
	if len(errs) > 0 {
		return cty.UnknownVal(cty.String), fmt.Errorf("invalid Kubernetes object name: %v", errs)
	}

	return k8s.NewObjectReference(name), nil
}
