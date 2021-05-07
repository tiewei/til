package funcs_test

import (
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/util/validation"

	. "bridgedl/lang/funcs"
)

func TestSecretRefFunc(t *testing.T) {
	secrRefFn := SecretRefFunc()

	testCases := map[string]struct {
		params    []cty.Value
		expectErr bool
	}{
		"valid name and key": {
			params: []cty.Value{
				cty.StringVal("123-abc"),
				cty.StringVal("key.name"),
			},
			expectErr: false,
		},
		"name is invalid": {
			params: []cty.Value{
				cty.StringVal("123_abc"),
				cty.StringVal("key.name"),
			},
			expectErr: true,
		},
		"key is too long": {
			params: []cty.Value{
				cty.StringVal("123_abc"),
				cty.StringVal(strings.Repeat("a", validation.DNS1123SubdomainMaxLength+1)),
			},
			expectErr: true,
		},
		"no argument": {
			params:    []cty.Value{},
			expectErr: true,
		},
		"no key": {
			params: []cty.Value{
				cty.StringVal("123-abc"),
			},
			expectErr: true,
		},
		"too many arguments": {
			params: []cty.Value{
				cty.StringVal("123-abc"),
				cty.StringVal("key.name"),
				cty.StringVal("extra-arg"),
			},
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			_, err := secrRefFn.Call(tc.params)

			if tc.expectErr && err == nil {
				t.Error("Expected function call to return an error")
			}
			if !tc.expectErr && err != nil {
				t.Error("Function call returned an error:", err)
			}
		})
	}
}
