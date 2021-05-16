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

package funcs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"bridgedl/lang/k8s"
)

// SecretRefFunc returns the implementation of the "secret_ref" function.
//
// secret_ref() creates a reference to a data key inside a Kubernetes Secret
// object.
//
// Parameters:
//  * name: name of the Secret object. Must be a valid Kubernetes object name (RFC 1123 subdomain).
//  * key:  data key to reference. Must consist of alphanumeric characters, '-', '_' or '.'.
//
func SecretRefFunc() function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "name",
				Type: cty.String,
			},
			{
				Name: "key",
				Type: cty.String,
			},
		},

		Type: function.StaticReturnType(k8s.SecretKeySelectorCty),
		Impl: secretRefFuncImpl,
	})
}

func secretRefFuncImpl(args []cty.Value, _ cty.Type) (cty.Value, error) {
	name := args[0].AsString()
	key := args[1].AsString()

	errs := append(
		validation.IsDNS1123Subdomain(name),
		validation.IsConfigMapKey(key)...,
	)
	if len(errs) > 0 {
		return cty.UnknownVal(cty.String), fmt.Errorf("invalid Kubernetes Secret reference: %v", errs)
	}

	return k8s.NewSecretKeySelector(name, key), nil
}
