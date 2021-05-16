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

package config

import "strconv"

// ComponentCategory represents a category of Bridge component.
//
// The purpose of this type is to expose an exhaustive and strongly typed list
// of block types that can appear in Bridge Description Files, which is mostly
// useful outside of this package.
type ComponentCategory int8

const (
	CategoryUnknown ComponentCategory = iota - 1 // invalid
	CategoryChannels
	CategoryRouters
	CategoryTransformers
	CategorySources
	CategoryTargets
)

// String implements fmt.Stringer.
func (c ComponentCategory) String() string {
	switch c {
	case CategoryChannels:
		return BlkChannel
	case CategoryRouters:
		return BlkRouter
	case CategoryTransformers:
		return BlkTransf
	case CategorySources:
		return BlkSource
	case CategoryTargets:
		return BlkTarget
	default:
		return "config.ComponentCategory(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

// AsComponentCategory returns the category of Bridge component corresponding
// to the given string.
func AsComponentCategory(s string) ComponentCategory {
	switch s {
	case BlkChannel:
		return CategoryChannels
	case BlkRouter:
		return CategoryRouters
	case BlkTransf:
		return CategoryTransformers
	case BlkSource:
		return CategorySources
	case BlkTarget:
		return CategoryTargets
	default:
		return CategoryUnknown
	}
}
