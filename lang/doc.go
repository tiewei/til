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

// Package lang deals with the runtime aspects of the Bridge Description
// Language, such as decoding the remainder of partial contents, or evaluating
// functions and expressions.
//
// It can be considered as a sibling package of "config", which deals with
// static parsing and validation only.
package lang
