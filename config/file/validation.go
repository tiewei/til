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

package file

import (
	"math/big"

	"github.com/zclconf/go-cty/cty"
)

// isInt64 returns whether the given cty.Value is a cty.Number that can be
// represented as an int64.
func isInt64(v cty.Value) bool {
	if v.Type() != cty.Number {
		return false
	}

	bigInt, accuracy := v.AsBigFloat().Int(nil)
	if accuracy != big.Exact {
		return false
	}

	return bigInt.IsInt64()
}
