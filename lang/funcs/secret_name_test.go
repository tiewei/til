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
