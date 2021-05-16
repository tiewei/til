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
