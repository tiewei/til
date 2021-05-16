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

package funcs_test

import (
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/util/validation"

	. "bridgedl/lang/funcs"
)

func TestSecretNameFunc(t *testing.T) {
	secrNameFn := SecretNameFunc()

	testCases := map[string]struct {
		params    []cty.Value
		expectErr bool
	}{
		"valid RFC 1123 subdomain": {
			params: []cty.Value{
				cty.StringVal("123-abc"),
			},
			expectErr: false,
		},
		"contains invalid characters": {
			params: []cty.Value{
				cty.StringVal("123_abc"),
			},
			expectErr: true,
		},
		"starts with a non-alphanumeric character": {
			params: []cty.Value{
				cty.StringVal("-123-abc"),
			},
			expectErr: true,
		},
		"too long": {
			params: []cty.Value{
				cty.StringVal(strings.Repeat("a", validation.DNS1123SubdomainMaxLength+1)),
			},
			expectErr: true,
		},
		"no argument": {
			params:    []cty.Value{},
			expectErr: true,
		},
		"too many arguments": {
			params: []cty.Value{
				cty.StringVal("123-abc"),
				cty.StringVal("my-name"),
			},
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			_, err := secrNameFn.Call(tc.params)

			if tc.expectErr && err == nil {
				t.Error("Expected function call to return an error")
			}
			if !tc.expectErr && err != nil {
				t.Error("Function call returned an error:", err)
			}
		})
	}
}
