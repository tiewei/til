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

package cli

import (
	"context"
	"errors"
	"flag"
)

// flagSetKey is the Context key for a flag.FlagSet.
type flagSetKey struct{}

// withFlagSet infuses ctx with the given flag.Flagset.
func withFlagSet(ctx context.Context, f *flag.FlagSet) context.Context {
	return context.WithValue(ctx, flagSetKey{}, f)
}

// FlagSetFromContext retrieves the flag.FlagSet injected into the given context.Context.
func FlagSetFromContext(ctx context.Context) *flag.FlagSet {
	f, ok := ctx.Value(flagSetKey{}).(*flag.FlagSet)
	if !ok {
		panic(errors.New("flag.FlagSet not found in context"))
	}

	return f
}

// uiKey is the Context key for a UI.
type uiKey struct{}

// withUI infuses ctx with the given UI.
func withUI(ctx context.Context, ui *UI) context.Context {
	return context.WithValue(ctx, uiKey{}, ui)
}

// UIFromContext retrieves the UI injected into the given context.Context.
func UIFromContext(ctx context.Context) *UI {
	ui, ok := ctx.Value(uiKey{}).(*UI)
	if !ok {
		panic(errors.New("UI not found in context"))
	}

	return ui
}
