package targets

// All includes all "target" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"container": (*Container)(nil),
	"function":  (*Function)(nil),
	"kafka":     (*Kafka)(nil),
}
