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

// Package graph provides primitives for building and traversing directed
// graphs.
//
// The implementation is purposely generic and could be used for representing
// any kind of data. In practice, the main application of those graphs is via
// the child "build" package for representing the different messaging
// components of a Bridge, as well as the data flow between those components.
package graph
