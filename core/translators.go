package core

import "bridgedl/translation"

// Translators encapsulates TranslatorProviders for all block types
// that support translation.
type Translators struct {
	Routers translation.TranslatorProvider
}
