package lang

import (
	"github.com/zclconf/go-cty/cty/function"

	"bridgedl/fs"
	"bridgedl/lang/funcs"
)

// Functions returns all functions supported by the Bridge Description
// Language, scoped at basedir for functions that access the file system via
// the fs interface.
func Functions(basedir string, fs fs.FS) map[string]function.Function {
	return map[string]function.Function{
		"file":        funcs.FileFunc(basedir, fs),
		"secret_name": funcs.SecretNameFunc(),
		"secret_ref":  funcs.SecretRefFunc(),
	}
}
