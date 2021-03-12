package lookup

import (
	"bridgedl/config"
	"bridgedl/config/addr"
)

// ChannelKey returns a lookup key which can be used to index the given Channel
// inside a config.Bridge.
func ChannelKey(ch *config.Channel) interface{} {
	return addr.Channel{Identifier: ch.Identifier}
}

// RouterKey returns a lookup key which can be used to index the given Router
// inside a config.Bridge.
func RouterKey(rtr *config.Router) interface{} {
	return addr.Router{Identifier: rtr.Identifier}
}

// TransformerKey returns a lookup key which can be used to index the given Transformer
// inside a config.Bridge.
func TransformerKey(trsf *config.Transformer) interface{} {
	return addr.Transformer{Identifier: trsf.Identifier}
}

// SourceKey returns a lookup key which can be used to index the given Source
// inside a config.Bridge.
func SourceKey(src *config.Source) interface{} {
	return addr.Source{Identifier: src.Identifier}
}

// TargetKey returns a lookup key which can be used to index the given Target
// inside a config.Bridge.
func TargetKey(rtr *config.Target) interface{} {
	return addr.Target{Identifier: rtr.Identifier}
}

// FunctionKey returns a lookup key which can be used to index the given Function
// inside a config.Bridge.
func FunctionKey(fn *config.Function) interface{} {
	return addr.Function{Identifier: fn.Identifier}
}
