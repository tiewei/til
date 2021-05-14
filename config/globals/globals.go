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

package globals

import "github.com/zclconf/go-cty/cty"

// Accessor exposes global Bridge settings.
type Accessor interface {
	Delivery() *Delivery
}

// Delivery is a clone of config.Delivery where the DeadLetterSink expression
// field has been evaluated to a concrete value using a hcl.EvalContext.
type Delivery struct {
	Retries        *int64
	DeadLetterSink cty.Value
}
