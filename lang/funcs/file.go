package funcs

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"bridgedl/fs"
)

// FileFunc returns the implementation of the "file" function.
//
// file() reads the contents of a file and returns it as a string.
//
// Parameters:
//  * path: path of the file, either absolute or relative to the directory of
//    the file that calls the function.
//
func FileFunc(basedir string, filesyst fs.FS) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{
			{
				Name: "path",
				Type: cty.String,
			},
		},

		Type: function.StaticReturnType(cty.String),
		Impl: fileFuncImpl(basedir, filesyst),
	})
}

func fileFuncImpl(basedir string, fs fs.FS) function.ImplFunc {
	return func(args []cty.Value, _ cty.Type) (cty.Value, error) {
		path := args[0].AsString()
		if !filepath.IsAbs(path) {
			path = filepath.Join(basedir, path)
		}

		fd, err := fs.Open(path)
		if err != nil {
			return cty.UnknownVal(cty.String), fmt.Errorf("opening file: %w", err)
		}
		defer fd.Close()

		data, err := io.ReadAll(fd)
		if err != nil {
			return cty.UnknownVal(cty.String), fmt.Errorf("reading file contents: %w", err)
		}

		return cty.StringVal(string(data)), nil
	}
}
