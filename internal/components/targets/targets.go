package targets

// All includes all "target" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"aws_kinesis": (*AWSKinesis)(nil),
	"aws_lambda":  (*AWSLambda)(nil),
	"aws_s3":      (*AWSS3)(nil),
	"aws_sns":     (*AWSSNS)(nil),
	"aws_sqs":     (*AWSSQS)(nil),
	"container":   (*Container)(nil),
	"function":    (*Function)(nil),
	"kafka":       (*Kafka)(nil),
	"zendesk":     (*Zendesk)(nil),
}
