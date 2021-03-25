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
	testFS.CreateFile(filepath.Join(baseDir, fakeFileRelPath),
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
