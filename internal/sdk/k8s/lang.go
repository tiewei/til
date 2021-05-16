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

package k8s

import "bridgedl/lang/k8s"

// This file contains aliases to variables and functions from the lang/k8s
// package that are commonly used in component implementations, so that those
// components only have to depend on a single import.

var NewDestination = k8s.NewDestination
var DestinationCty = k8s.DestinationCty

var ObjectReferenceCty = k8s.ObjectReferenceCty
var IsObjectReference = k8s.IsObjectReference

var IsSecretKeySelector = k8s.IsSecretKeySelector
