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
	"path/filepath"
	"testing"

	"github.com/zclconf/go-cty/cty"

	"bridgedl/fs"
	. "bridgedl/lang/funcs"
)

func TestFileFunc(t *testing.T) {
	const baseDir = "/fake"

	const fakeFileRelPath = "somedir/somefile.json"
	const fakeFileContents = `{"msg": "Hello, World!"}`

	testFS := fs.NewMemFS()
	_ = testFS.CreateFile(filepath.Join(baseDir, fakeFileRelPath),
		[]byte(fakeFileContents),
	)

	fileFn := FileFunc(baseDir, testFS)

	testCases := map[string]struct {
		params    []cty.Value
		expectErr bool
	}{
		"valid relative path": {
			params: []cty.Value{
				cty.StringVal(fakeFileRelPath),
			},
			expectErr: false,
		},
		"valid absolute path": {
			params: []cty.Value{
				cty.StringVal(filepath.Join(baseDir, fakeFileRelPath)),
			},
			expectErr: false,
		},
		"file not found": {
			params: []cty.Value{
				cty.StringVal("no/file/here"),
			},
			expectErr: true,
		},
		"wrong param type": {
			params: []cty.Value{
				cty.NumberIntVal(42),
			},
			expectErr: true,
		},
		"no argument": {
			params:    []cty.Value{},
			expectErr: true,
		},
		"too many arguments": {
			params: []cty.Value{
				cty.StringVal(fakeFileRelPath),
				cty.StringVal(fakeFileRelPath),
			},
			expectErr: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			_, err := fileFn.Call(tc.params)

			if tc.expectErr && err == nil {
				t.Error("Expected function call to return an error")
			}
			if !tc.expectErr && err != nil {
				t.Error("Function call returned an error:", err)
			}
		})
	}
}
